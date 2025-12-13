import { ProjectStatusEnum } from "@/src/utils/enums";

export interface IProject {
    id: string,
    createdAt: string,
    updatedAt: string,
    name: string,
    teamId: string,
    managerId: string,
    status: ProjectStatusEnum
    toDoTasks: number,
    inProgressTasks: number,
    doneTasks: number
}

export interface IGetProjectParams {
    teamId?: string;
    managerId?: string;
    name?: string;
    limit: number;
    cursor: string;
    withStats: boolean;
}


export interface ICreateProjectBody {
    name: string;
}


export interface IUpdateProjectParams {
    name: string;
    status: ProjectStatusEnum;
    id: string;
}