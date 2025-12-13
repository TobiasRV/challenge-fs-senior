export interface ITeam {
    id: string;
    createdAt: string;
    updatedAt: string;
    name: string;
    ownerId: string;
}

export interface ICreateTeam {
    name: string
}