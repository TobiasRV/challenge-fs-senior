package models

import (
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
)

type Projectstatus string

const (
	ProjectstatusOnHold     Projectstatus = "OnHold"
	ProjectstatusInProgress Projectstatus = "InProgress"
	ProjectstatusCompleted  Projectstatus = "Completed"
)

type Project struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Name      string        `json:"name"`
	TeamID    uuid.UUID     `json:"teamId"`
	ManagerID uuid.UUID     `json:"managerId"`
	Status    Projectstatus `json:"status"`
}

func DatabaseProjectToProject(dbProject database.Project) Project {
	return Project{
		ID:        dbProject.ID,
		CreatedAt: dbProject.CreatedAt,
		UpdatedAt: dbProject.UpdatedAt,
		Name:      dbProject.Name,
		TeamID:    dbProject.TeamID,
		ManagerID: dbProject.ManagerID,
		Status:    Projectstatus(dbProject.Status),
	}
}

func DatabaseProjectsToProjects(dbProjects []database.Project) []Project {
	res := []Project{}
	for _, p := range dbProjects {
		res = append(res, DatabaseProjectToProject(p))
	}

	return res
}
