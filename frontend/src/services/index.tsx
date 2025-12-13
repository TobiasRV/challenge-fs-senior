import axios, { AxiosResponse, HttpStatusCode, InternalAxiosRequestConfig } from "axios";

import { formatUser } from "../utils/helpers";
import { localStorageKeys } from "../utils/consts";
import { clearLs, getLsItem, setLsItem } from "../utils/localStorage";
import { refreshTokenService } from "./apiServices/auth/service";
import { loginRoute, registerAdminRoute } from "./apiServices/auth/routes";


const BASE_URL_API = process.env.NEXT_PUBLIC_BASE_URL_API;
let isRefreshing = false;
let refreshTokenPromise: Promise<any> | null = null;

const server = axios.create({
  baseURL: BASE_URL_API,
});


const redirectToLogin = () => {
  clearLs();
  if (typeof window !== "undefined") {
    window.location.replace("/");
  }
};

const refreshToken = async () => {
  const refreshToken = getLsItem(localStorageKeys.REFRESH_TOKEN);
  if (refreshToken) {
    const response = (await refreshTokenService({ refreshToken })) as any;
    if (
      response.statusCode === HttpStatusCode.BadRequest ||
      response.statusCode === HttpStatusCode.InternalServerError ||
      response.statusCode === HttpStatusCode.Unauthorized
    ) {
      redirectToLogin();
      return;
    }
    return response.accessToken;
  } else {
    redirectToLogin();
  }
};

server.interceptors.request.use(
  (config: InternalAxiosRequestConfig<any>) => {
    const token = getLsItem(localStorageKeys.TOKEN);
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error: any) => Promise.reject(error),
);

server.interceptors.response.use(
  (response: AxiosResponse<any, any, {}>) => {
    if (
      response.config.url && [loginRoute(), registerAdminRoute()].includes(response.config.url)
    ) {
      const { accessToken, refreshToken, user } = response.data;

      if (accessToken) {
        setLsItem(localStorageKeys.TOKEN, accessToken);
      }

      if (refreshToken) {
        setLsItem(localStorageKeys.REFRESH_TOKEN, refreshToken);
      }

      if (user) {
        setLsItem(localStorageKeys.USER, formatUser(user))
      }
    }

    return response;
  },
  async (error: any) => {
    const originalRequest = error.config;

    const ingoredRoutes: Array<string> = [];

    if (ingoredRoutes.includes(originalRequest.url)) {
      return Promise.reject(error);
    }

    // Error Control to refresh tokens
    if (
      error.response.status === HttpStatusCode.Unauthorized &&
      !isRefreshing
    ) {
      isRefreshing = true;
      refreshTokenPromise = refreshTokenPromise || refreshToken();

      try {
        const token = await refreshTokenPromise;
        originalRequest.headers.Authorization = `Bearer ${token}`;

        return axios(originalRequest);
      } finally {
        isRefreshing = false;
        refreshTokenPromise = null;
      }
    }
    return Promise.reject(error);
  },
);

export default server;
