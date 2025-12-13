package interfaces

import (
	"context"
	"database/sql"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/TobiasRV/challenge-fs-senior/internals/utils"
	"github.com/google/uuid"
)

type ITaskRepository interface {
	CreateTask(context.Context, database.CreateTasksParams) (models.Task, error)
	GetTasks(context.Context, GetTasksFilters) ([]GetTasksResponse, error)
	GetTaskById(context.Context, uuid.UUID) (models.Task, error)
	UpdateTask(context.Context, UpdateTaskData) (models.Task, error)
	DeleteTask(context.Context, uuid.UUID) error
}

type TasksListResponse struct {
	Data       []GetTasksResponse `json:"data"`
	Pagination utils.Pagination   `json:"pagination"`
}

type DeleteTaskResponse struct {
	Deleted bool `json:"deleted" example:"true"`
}

type CreateTaksPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	ProjectId   string `json:"projectId" validate:"required,uuid"`
	UserId      string `json:"userId"`
}

type GetTasksParams struct {
	Limit     uint64 `query:"limit,required"`
	Cursor    string `query:"cursor"`
	Title     string `query:"title"`
	ProjectId string `query:"projectId"`
}

type GetTasksFilters struct {
	Limit           uint64
	IsFirstPage     bool
	PointsNext      bool
	CursorCreatedAt time.Time
	CursorId        uuid.UUID
	Title           string
	ProjectId       uuid.UUID
	UserId          uuid.UUID
}

type GetTasksResponse struct {
	ID          uuid.UUID         `json:"id"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	ProjectID   uuid.UUID         `json:"projectId"`
	UserID      uuid.NullUUID     `json:"userId"`
	Status      models.Taskstatus `json:"status"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	ProjectName string            `json:"projectName"`
	UserName    string            `json:"userName"`
}

type UpdateTaskData struct {
	Title       string
	Description sql.NullString
	Status      models.Taskstatus
	UserId      uuid.NullUUID
	UpdatedAt   time.Time
	ID          uuid.UUID
}

type UpdateTaskPayload struct {
	Title       string            `json:"title" validate:"required"`
	Description string            `json:"description"`
	Status      models.Taskstatus `json:"status" validate:"required,oneof=ToDo InProgress Done"`
	UserId      uuid.UUID         `json:"userId"`
}
