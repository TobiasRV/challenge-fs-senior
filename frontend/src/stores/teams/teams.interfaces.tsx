export interface ITeam {
    id: string;
    createdAt: Date;
    updatedAt: Date;
    name: string;
    ownerId: string;
}

export interface ICreateTeam {
    name: string
}