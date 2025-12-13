import {
  ICreateProjectBody,
  IGetProjectParams,
  IUpdateProjectParams,
} from "@/src/stores/projects/projects.interface";
import server from "../..";
import handleAxiosErrors from "../../axios.helper";
import { ICreateTaskBody, IGetTasksParams, IUpdateTaskParams } from "@/src/stores/tasks/tasks.interface";
import { createTasksRoute, deleteTasksRoute, getTasksRoute, updateTasksRoute } from "./routes";


export const getTasks = async (params: IGetTasksParams) => {
  try {
    const response = await server.get(getTasksRoute(), {
      params,
    });
    return { ...response.data, statusCode: response.status };
  } catch (error) {
    return handleAxiosErrors(error);
  }
};

export const createTask = async (body: ICreateTaskBody) => {
  try {
    const response = await server.post(createTasksRoute(), body);
    return { ...response.data, statusCode: response.status };
  } catch (error) {
    return handleAxiosErrors(error);
  }
};

export const updateTask = async (params: IUpdateTaskParams) => {
    try {
    const response = await server.put(updateTasksRoute(params.id), {
        title: params.title,
        description: params.description,
        status: params.status,
        userId: params.userId,
    });
    return { ...response.data, statusCode: response.status };
  } catch (error) {
    return handleAxiosErrors(error);
  }
};

export const deleteTask = async (id: string) => {
    try {
    const response = await server.delete(deleteTasksRoute(id));
    return { ...response.data, statusCode: response.status };
  } catch (error) {
    return handleAxiosErrors(error);
  }
};

