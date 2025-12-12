import { emailExistsService } from "@/src/services/apiServices/users/service";
import { HttpStatusCode } from "axios";
import { create } from "zustand";

interface UserStore {
  loading: boolean;
  error: boolean;
  statusCode: number | null;
  emailAlreadyExists: (email: string) => Promise<boolean>;
  clearRequestState: () => void;
}

const initialState = {
  isLoggedIn: false,
  user: null,
  loading: false,
  error: false,
  statusCode: null,
};

export const useUserStore = create<UserStore>((set, get) => ({
  ...initialState,
  emailAlreadyExists: async (email: string): Promise<boolean> => {
    try {
      set({ loading: true });
      const response = await emailExistsService(email);

      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.status,
        });
        return true;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.status
      });

      return response.exists;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
      });
      return true;
    }
  },

  clearRequestState: (): void => {
    set({
      error: false,
      statusCode: null,
      loading: false,
    });
  },
}));
