import { HttpStatusCode } from "axios";
import { create } from "zustand";
import { ICreateTaskBody, IGetTasksParams, ITask, IUpdateTaskParams } from "./tasks.interface";
import { createTaskService, deleteTaskService, getTasksService, updateTaskService } from "@/src/services/apiServices/tasks/service";

interface TasksStore {
  loading: boolean;
  error: boolean;
  statusCode: number | null;
  tasks: {
    data: Array<ITask>;
    prev: string;
    next: string;
  };
  getTasks: (filters: IGetTasksParams) => Promise<void>;
  createTask: (body: ICreateTaskBody) => Promise<number>;
  updateTask: (body: IUpdateTaskParams) => Promise<number>;
  deleteTask: (id: string) => Promise<number>;
  clearRequestState: () => void;
}

const initialState = {
  isLoggedIn: false,
  user: null,
  loading: false,
  error: false,
  statusCode: null,
  tasks: {
    data: [],
    prev: "",
    next: "",
  },
};

export const useTaskStore = create<TasksStore>((set, get) => ({
  ...initialState,

  getTasks: async (filters: IGetTasksParams): Promise<void> => {
    try {
      set({ loading: true });
      const response = await getTasksService(filters);
      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
          tasks: {
            data: [],
            prev: "",
            next: "",
          },
        });
        return;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.statusCode,
        tasks: {
          data: response.data || [],
          prev: response.pagination.prev_cursor,
          next: response.pagination.next_cursor,
        },
      });

      return;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
        tasks: {
          data: [],
          prev: "",
          next: "",
        },
      });
      return;
    }
  },
  createTask: async (body: ICreateTaskBody): Promise<number> => {
    try {
      set({ loading: true });
      const response = await createTaskService(body);
      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
        });
        return response.statusCode;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.statusCode,
      });

      return response.statusCode;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
      });
      return HttpStatusCode.InternalServerError;
    }
  },
  updateTask: async (params: IUpdateTaskParams): Promise<number> => {
    try {
      set({ loading: true });
      const response = await updateTaskService(params);
      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
        });
        return response.statusCode;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.statusCode,
      });

      return response.statusCode;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
      });
      return HttpStatusCode.InternalServerError;
    }
  },
  deleteTask: async (id: string): Promise<number> => {
    try {
      set({ loading: true });
      const response = await deleteTaskService(id);
      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
        });
        return response.statusCode;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.statusCode,
      });

      return response.statusCode;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
      });
      return HttpStatusCode.InternalServerError;
    }
  },
  clearRequestState: (): void => {
    set({
      error: false,
      statusCode: null,
      loading: false,
    });
  },
}));
