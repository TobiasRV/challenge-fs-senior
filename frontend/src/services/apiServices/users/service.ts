import server from "../..";
import handleAxiosErrors from "../../axios.helper";
import { getEmailExistsRoute } from "./routes";

export const emailExistsService = async (email: string) => {
    try {
        const response = await server.get(getEmailExistsRoute(), {
            params: {
                email
            }
        });
        return {...response.data, statusCode: response.status}
    } catch (error) {
        return handleAxiosErrors(error);
    }
}