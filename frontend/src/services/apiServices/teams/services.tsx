import { ICreateTeam } from "@/src/stores/teams/teams.interfaces";
import server from "../..";
import handleAxiosErrors from "../../axios.helper";
import { getTeamByOwnerRoute, getCreateTeamRoute } from "./routes";

export const getTeamByOwner = async () => {
    try {
        const response = await server.get(getTeamByOwnerRoute());
        return {...response.data, statusCode: response.status}
    } catch (error) {
        return handleAxiosErrors(error);
    }
}

export const createTeam = async (payload: ICreateTeam) => {
    const { name } = payload;
    try {
        const response = await server.post(getCreateTeamRoute(), {
            name
        });
        return {...response.data, statusCode: response.status};
    } catch (error) {
        return handleAxiosErrors(error);
    }
}