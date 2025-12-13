import { UserRolesEnum } from "@/src/utils/enums";

export interface IUser {
    id: string;
    username: string;
    email: string;
    role: UserRolesEnum
    teamId?: string;
}

export interface IGetUsersParams {
    email?: string;
    teamId: string;
    limit: number;
    cursor: string;
}

export interface IUpdateUsersParams {
    username: string;
    email: string;
    id: string
}