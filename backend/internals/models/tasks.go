package models

import (
	"database/sql"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
)

type Taskstatus string

const (
	TaskstatusToDo       Taskstatus = "ToDo"
	TaskstatusInProgress Taskstatus = "InProgress"
	TaskstatusDone       Taskstatus = "Done"
)

type Task struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	ProjectID   uuid.UUID      `json:"projectId"`
	UserID      uuid.NullUUID  `json:"userId"`
	Status      Taskstatus     `json:"status"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
}

func DatabaseTaskToTask(dbTask database.Task) Task {
	return Task{
		ID:        dbTask.ID,
		CreatedAt: dbTask.CreatedAt,
		UpdatedAt: dbTask.UpdatedAt,
		ProjectID: dbTask.ProjectID,
		UserID: uuid.NullUUID{
			UUID:  dbTask.UserID.UUID,
			Valid: dbTask.UserID.Valid,
		},
		Status:      Taskstatus(dbTask.Status),
		Title:       dbTask.Title,
		Description: dbTask.Description,
	}
}

func DatabaseTasksToTasks(dbTasks []database.Task) []Task {
	res := []Task{}
	for _, u := range dbTasks {
		res = append(res, DatabaseTaskToTask(u))
	}

	return res
}
