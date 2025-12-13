import { HttpStatusCode } from "axios";
import { create } from "zustand";
import {
  ICreateProjectBody,
  IGetProjectParams,
  IProject,
  IUpdateProjectParams,
} from "./projects.interface";
import {
  createProjectsService,
  deleteProjectService,
  getProjectsService,
  updateProjectService,
} from "@/src/services/apiServices/projects/service";

interface ProjectStore {
  loading: boolean;
  error: boolean;
  statusCode: number | null;
  projects: {
    data: Array<IProject>;
    prev: string;
    next: string;
  };
  getProjects: (filters: IGetProjectParams) => Promise<void>;
  createProject: (body: ICreateProjectBody) => Promise<number>;
  updateProject: (body: IUpdateProjectParams) => Promise<number>;
  deleteProject: (id: string) => Promise<number>;
  clearRequestState: () => void;
}

const initialState = {
  isLoggedIn: false,
  user: null,
  loading: false,
  error: false,
  statusCode: null,
  projects: {
    data: [],
    prev: "",
    next: "",
  },
};

export const useProjectStore = create<ProjectStore>((set, get) => ({
  ...initialState,

  getProjects: async (filters: IGetProjectParams): Promise<void> => {
    try {
      set({ loading: true });
      const response = await getProjectsService(filters);
      if (response.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.statusCode,
          projects: {
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
        projects: {
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
        projects: {
          data: [],
          prev: "",
          next: "",
        },
      });
      return;
    }
  },
  createProject: async (body: ICreateProjectBody): Promise<number> => {
    try {
      set({ loading: true });
      const response = await createProjectsService(body);
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
  updateProject: async (params: IUpdateProjectParams): Promise<number> => {
    try {
      set({ loading: true });
      const response = await updateProjectService(params);
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
  deleteProject: async (id: string): Promise<number> => {
    try {
      set({ loading: true });
      const response = await deleteProjectService(id);
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
