import { ProjectStatusEnum, TaskStatusEnum, UserRolesEnum } from "./enums";

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

export const TasksStatusTranslation = {
    [TaskStatusEnum.TODO]: "Por hacer",
    [TaskStatusEnum.IN_PROGRESS]: "En curso",
    [TaskStatusEnum.DONE]: "Completa",
}