import {
  createUsersService,
  deleteUserService,
  emailExistsService,
  getUsersService,
  updateUsersService,
} from "@/src/services/apiServices/users/service";
import { HttpStatusCode } from "axios";
import { create } from "zustand";
import { IGetUsersParams, IUpdateUsersParams, IUser } from "./users.interface";

interface UserStore {
  loading: boolean;
  error: boolean;
  statusCode: number | null;
  users: {
    data: Array<IUser>;
    prev: string;
    next: string;
  };
  emailAlreadyExists: (email: string) => Promise<boolean>;
  getUsers: (filters: IGetUsersParams) => Promise<void>;
  createUser: (userData: Partial<IUser>) => Promise<number>;
  updateUser: (userData: IUpdateUsersParams) => Promise<number>;
  deleteUser: (id: string) => Promise<number>;
  clearRequestState: () => void;
  clearState: () => void;
}

const initialState = {
  isLoggedIn: false,
  user: null,
  loading: false,
  error: false,
  statusCode: null,
  users: {
    data: [],
    prev: "",
    next: "",
  },
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
        statusCode: response.status,
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
  getUsers: async (filters: IGetUsersParams): Promise<void> => {
    try {
      set({ loading: true });
      const response = await getUsersService(filters);
      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
          users: {
            data: [],
            prev: "",
            next: "",
          },
        });
        return;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.statusCode,
        users: {
          data: response.data || [],
          prev: response.pagination.prev_cursor,
          next: response.pagination.next_cursor,
        },
      });

      return;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
        users: {
          data: [],
          prev: "",
          next: "",
        },
      });
      return;
    }
  },
  createUser: async (userData: Partial<IUser>): Promise<number> => {
    try {
      set({ loading: true });
      const response = await createUsersService(userData);

      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
        });
        return response.statusCode;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.statusCode,
      });

      return response.statusCode;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
      });
      return HttpStatusCode.InternalServerError;
    }
  },
  updateUser: async (userData: IUpdateUsersParams): Promise<number> => {
    try {
      set({ loading: true });
      const response = await updateUsersService(userData);

      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
        });
        return response.statusCode;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.statusCode,
      });

      return response.statusCode;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
      });
      return HttpStatusCode.InternalServerError;
    }
  },
  deleteUser: async (id: string): Promise<number> => {
    try {
      set({ loading: true });
      const response = await deleteUserService(id);

      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
        });
        return response.statusCode;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.statusCode,
      });

      return response.statusCode;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
      });
      return HttpStatusCode.InternalServerError;
    }
  },
  clearRequestState: (): void => {
    set({
      error: false,
      statusCode: null,
      loading: false,
    });
  },
  clearState: (): void => {
    set(initialState);
  },
}));
