import { ProjectStatusEnum, UserRolesEnum } from "./enums";

export const UserRolesTranslation = {
    [UserRolesEnum.ADMIN]: "Admin",
    [UserRolesEnum.MANAGER]: "Manager",
    [UserRolesEnum.MEMBER]: "Miembro"
}

export const ProjectStatusTranslation = {
    [ProjectStatusEnum.ON_HOLD]: "En pausa",
    [ProjectStatusEnum.IN_PROGRESS]: "En curso",
    [ProjectStatusEnum.COMPLETED]: "Completo",
}