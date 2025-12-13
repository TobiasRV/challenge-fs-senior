import { IGetUsersParams, IUpdateUsersParams, IUser } from "@/src/stores/users/users.interface";
import server from "../..";
import handleAxiosErrors from "../../axios.helper";
import { createUsersRoute, deleteUsersRoute, getEmailExistsRoute, getUsersRoute, updateUsersRoute } from "./routes";

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

export const getUsers = async (params: IGetUsersParams) => {
    try {
        const response = await server.get(getUsersRoute(), {
            params
        });
        return {...response.data, statusCode: response.status}
    } catch (error) {
        return handleAxiosErrors(error);
    }
}

export const createUsers = async (body: Partial<IUser>) => {
    try {
        const response = await server.post(createUsersRoute(), body);
        return {...response.data, statusCode: response.status}
    } catch (error) {
        return handleAxiosErrors(error);
    }
}

export const updateUsers = async (body: IUpdateUsersParams) => {
    try {
        const response = await server.put(updateUsersRoute(body.id), {
            username: body.username,
            email: body.email
        });
        return {...response.data, statusCode: response.status}
    } catch (error) {
        return handleAxiosErrors(error);
    }
}

export const deleteUser = async (id: string) => {
    try {
        const response = await server.delete(deleteUsersRoute(id));
        return {...response.data, statusCode: response.status}
    } catch (error) {
        return handleAxiosErrors(error);
    }
}