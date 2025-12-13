package interfaces

import (
	"context"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
)

type IProjectRepository interface {
	CreateProject(context.Context, database.CreateProjectParams) (models.Project, error)
	GetProjects(context.Context, GetProjectsFilters) ([]GetProjectsResponse, error)
	UpdateProject(context.Context, UpdateProjectData) (models.Project, error)
	GetProjectById(context.Context, uuid.UUID) (models.Project, error)
	DeleteProject(context.Context, uuid.UUID) error
}

type GetProjectsParams struct {
	Name      string    `query:"name"`
	Limit     uint64    `query:"limit,required"`
	Cursor    string    `query:"cursor"`
	WithStats bool      `query:"withStats"`
	TeamId    uuid.UUID `query:"teamId,uuid"`
	ManagerId uuid.UUID `query:"managerId,uuid"`
}

type GetProjectsFilters struct {
	Name            string
	TeamId          uuid.UUID
	ManagerId       uuid.UUID
	Limit           uint64
	IsFirstPage     bool
	PointsNext      bool
	CursorCreatedAt time.Time
	CursorId        uuid.UUID
	WithStats       bool
}

type CreateProjectPayload struct {
	Name string `json:"name" validate:"required"`
}

type GetProjectsResponse struct {
	ID              uuid.UUID            `json:"id"`
	CreatedAt       time.Time            `json:"createdAt"`
	UpdatedAt       time.Time            `json:"updatedAt"`
	Name            string               `json:"name"`
	TeamID          uuid.UUID            `json:"teamId"`
	ManagerID       uuid.UUID            `json:"managerId"`
	Status          models.Projectstatus `json:"status"`
	ToDoTasks       int                  `json:"toDoTasks"`
	InProgressTasks int                  `json:"inProgressTasks"`
	DoneTasks       int                  `json:"doneTasks"`
}

type UpdateProjectData struct {
	Name      string
	Status    models.Projectstatus
	UpdatedAt time.Time
	ID        uuid.UUID
}

type UpdateProjectPayload struct {
	Name   string               `json:"name" validate:"required"`
	Status models.Projectstatus `json:"status" validate:"required,oneof=OnHold InProgress Completed"`
}
