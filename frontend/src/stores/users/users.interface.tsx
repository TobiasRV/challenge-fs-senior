import { UserRolesEnum } from "@/src/utils/enums";

export interface IUser {
    id: string;
    username: string;
    email: string;
    role: UserRolesEnum
}