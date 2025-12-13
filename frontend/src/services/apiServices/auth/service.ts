import handleAxiosErrors from "../../axios.helper";
import { loginRoute, logoutRoute, refreshTokenRoute, registerAdminRoute } from "./routes";
import server from "../..";
import { ILogInForm, ISignInForm } from "@/src/stores/auth/auth.interfaces";

export const loginService = async ( payload: ILogInForm ) => {
    const {email, password} = payload;

    try {
        const response = await server.post(loginRoute(), {
            email,
            password
        });
        return { ...response.data, statusCode: response.status}

    } catch (error) {
        return handleAxiosErrors(error);
    }
}

export const refreshTokenService = async (payload: { refreshToken: string })  => {
    const { refreshToken } = payload;

    try {
        const response = await server.post(refreshTokenRoute(), {
            refreshToken
        });

        return { accessToken: response.data, statusCode: response.status };

    } catch (error) {
        return handleAxiosErrors(error);
    }
}

export const signInService = async (payload: ISignInForm) => {
    const { username, password, email} = payload;
    try {
        const response = await server.post(registerAdminRoute(), {
            username,
            password,
            email
        });
        return {...response.data, statusCode: response.status}
    } catch (error) {
        return handleAxiosErrors(error);
    }
}

export const logoutService = async ( ) => {
    try {
        const response = await server.delete(logoutRoute());
        return { ...response.data, statusCode: response.status}

    } catch (error) {
        return handleAxiosErrors(error);
    }
}