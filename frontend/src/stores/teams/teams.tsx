import { HttpStatusCode } from "axios";
import { create } from "zustand";
import { ICreateTeam, ITeam } from "./teams.interfaces";
import { createTeam, getTeamByOwner } from "@/src/services/apiServices/teams/services";

interface TeamsStore {
  loading: boolean;
  error: boolean;
  statusCode: number | null;
  team: {
    exists: boolean,
    team: ITeam | null
  }
  getTeamByOwner: () => Promise<void>;
  createTeam: (payload: { name: string }) => Promise<number>;
  clearRequestState: () => void;
}

const initialState = {
  isLoggedIn: false,
  user: null,
  loading: false,
  error: false,
  statusCode: null,
  team: {
    exists: false,
    team: null
  }
};

export const useTeamStore = create<TeamsStore>((set, get) => ({
  ...initialState,
  getTeamByOwner: async (): Promise<void> => {
    try {
      set({ loading: true });
      const response = await getTeamByOwner();
      if (response?.error) {
        set({
          error: true,
          loading: false,
          statusCode: response.status,
          team: {
            exists: false,
            team: null
          }
        });
        return;
      }

      set({
        error: false,
        loading: false,
        statusCode: response.status,
        team: {
            exists: response.exists,
            team: response.team
        }
      });

      return;
    } catch (error) {
      set({
        error: false,
        loading: false,
        statusCode: HttpStatusCode.InternalServerError,
        team: {
            exists: false,
            team: null
        }
      });
      return;
    }
  },
  createTeam: async (payload: ICreateTeam): Promise<number> => {
    try {
      set({ loading: true });
      const response = await createTeam(payload);
      if (response?.error) {
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
