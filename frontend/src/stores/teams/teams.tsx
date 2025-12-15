import { HttpStatusCode } from "axios";
import { create } from "zustand";
import { ICreateTeam, ITeam } from "./teams.interfaces";
import { createTeamService, getTeamByOwnerService } from "@/src/services/apiServices/teams/services";
import { setLsItem } from "@/src/utils/localStorage";
import { localStorageKeys } from "@/src/utils/consts";

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
  clearState: () => void;
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
      const response = await getTeamByOwnerService();
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

      setLsItem(localStorageKeys.TEAM_ID, response.team.id)

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
      console.log("ENTRE")
      const response = await createTeamService(payload);
      console.log({ response });
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

      setLsItem(localStorageKeys.TEAM_ID, response.id)

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
  clearState: (): void => {
    set(initialState);
  },
}));
