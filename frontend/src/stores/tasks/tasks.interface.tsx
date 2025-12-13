import { TaskStatusEnum } from "@/src/utils/enums";

export interface ITask {
  id: string;
  createdAt: string;
  updatedAt: string;
  projectId: string;
  userId: string;
  status: TaskStatusEnum;
  title: string;
  description: string;
  projectName: string;
  userName: string;
}


export interface IGetTasksParams {
    projectId?: string;
    title?: string;
    limit: number;
    cursor: string;
}


export interface ICreateTaskBody {
    title: string;
    description?: string;
    projectId: string;
}


export interface IUpdateTaskParams {
    title: string;
    description?: string;
    status: TaskStatusEnum;
    userId: string;
    id: string;
}