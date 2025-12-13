import {
  loginService,
  logoutService,
  signInService,
} from "@/src/services/apiServices/auth/service";
import { localStorageKeys } from "@/src/utils/consts";
import { getLsItem, setLsItem, removeLsItem } from "@/src/utils/localStorage";
import { HttpStatusCode } from "axios";
import { create } from "zustand";
import { ILogInForm, ISignInForm } from "./auth.interfaces";
import { IUser } from "../users/users.interface";

interface AuthStore {
  isLoggedIn: boolean;
  user: IUser | null;
  loading: boolean;
  error: boolean;
  statusCode: number | null;
  logIn: (payload: ILogInForm) => Promise<number>;
  logOut: () => Promise<void>;
  registerAdmin: (payload: ISignInForm) => Promise<number>;
  clearRequestState: () => void;
  clearState: () => void;
}

const initialState = {
  isLoggedIn: false,
  user: null,
  loading: false,
  error: false,
  statusCode: null,
};

export const useAuthStore = create<AuthStore>((set, get) => ({
  ...initialState,
  user:
    (typeof window !== "undefined" && getLsItem(localStorageKeys.USER)) ||
    initialState.user,
  isLoggedIn:
    getLsItem(localStorageKeys.IS_LOGGED_IN) || initialState.isLoggedIn,
  logIn: async (payload: ILogInForm): Promise<number> => {
    try {
      set({ loading: true });

      const response = await loginService(payload);
      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
        });
        return response.statusCode;
      } else {
        set({
          error: false,
          loading: false,
          statusCode: response.statusCode,
          isLoggedIn: true,
          user: response.user,
        });
        setLsItem(localStorageKeys.IS_LOGGED_IN, true);
        setLsItem(localStorageKeys.USER, response.user);
      }
      return response.statusCode;
    } catch (error) {
      set({
        error: true,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
      });

      return HttpStatusCode.InternalServerError;
    }
  },
  logOut: async (): Promise<void> => {
    try {
      await logoutService();
    } catch (error) {
    } finally {
      set({
        ...initialState,
      });
      removeLsItem(localStorageKeys.USER);
      removeLsItem(localStorageKeys.IS_LOGGED_IN);
    }
  },
  registerAdmin: async (payload: ISignInForm): Promise<number> => {
    try {
      set({ loading: true });
      const response = await signInService(payload);
      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
        });
      } else {
        set({
          error: false,
          loading: false,
          statusCode: response.statusCode,
          isLoggedIn: true,
          user: response.user,
        });
        setLsItem(localStorageKeys.IS_LOGGED_IN, true);
        setLsItem(localStorageKeys.USER, response.user);
      }
      return response.statusCode;
    } catch (error) {
      set({
        error: true,
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
    set({
      ...initialState,
    });
    removeLsItem(localStorageKeys.USER);
    removeLsItem(localStorageKeys.IS_LOGGED_IN);
  },
}));
