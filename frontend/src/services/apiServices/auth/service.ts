import handleAxiosErrors from "../../axios.helper";
import { getLoginRoute, getRefreshTokenRoute, getRegisterAdminRoute } from "./routes";
import server from "../..";
import { ILogInForm, ISignInForm } from "@/src/stores/auth/auth.interfaces";

export const loginService = async ( payload: ILogInForm ) => {
    const {email, password} = payload;

    try {
        const response = await server.post(getLoginRoute(), {
            email,
            password
        });
        console.log({ response })
        return { ...response.data, statusCode: response.status}

    } catch (error) {
        console.log({ error })
        return handleAxiosErrors(error);
    }
}

export const refreshTokenService = async (payload: { refreshToken: string })  => {
    const { refreshToken } = payload;

    try {
        const response = await server.post(getRefreshTokenRoute(), {
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
        const response = await server.post(getRegisterAdminRoute(), {
            username,
            password,
            email
        });
        return {...response.data, statusCode: response.status}
    } catch (error) {
        return handleAxiosErrors(error);
    }
}
