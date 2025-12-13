import {
  ICreateProjectBody,
  IGetProjectParams,
  IUpdateProjectParams,
} from "@/src/stores/projects/projects.interface";
import server from "../..";
import handleAxiosErrors from "../../axios.helper";
import { createProjectsRoute, deleteProjectsRoute, getProjectsRoute, updateProjectsRoute } from "./routes";

export const getProjects = async (params: IGetProjectParams) => {
  try {
    const response = await server.get(getProjectsRoute(), {
      params,
    });
    return { ...response.data, statusCode: response.status };
  } catch (error) {
    return handleAxiosErrors(error);
  }
};

export const createProjects = async (body: ICreateProjectBody) => {
  try {
    const response = await server.post(createProjectsRoute(), body);
    return { ...response.data, statusCode: response.status };
  } catch (error) {
    return handleAxiosErrors(error);
  }
};

export const updateProject = async (params: IUpdateProjectParams) => {
    try {
    const response = await server.put(updateProjectsRoute(params.id), {
        name: params.name,
        status: params.status
    });
    return { ...response.data, statusCode: response.status };
  } catch (error) {
    return handleAxiosErrors(error);
  }
};

export const deleteProject = async (id: string) => {
    try {
    const response = await server.delete(deleteProjectsRoute(id));
    return { ...response.data, statusCode: response.status };
  } catch (error) {
    return handleAxiosErrors(error);
  }
};

