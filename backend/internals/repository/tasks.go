package repository

import (
	"context"
	stdsql "database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/TobiasRV/challenge-fs-senior/internals/interfaces"
	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
)

type TaskRepository struct {
	queries *database.Queries
	db      *stdsql.DB
}

func NewTaskRepository(queries *database.Queries, db *stdsql.DB) *TaskRepository {
	return &TaskRepository{
		queries: queries,
		db:      db,
	}
}

func (tsr *TaskRepository) CreateTask(c context.Context, data database.CreateTasksParams) (models.Task, error) {
	newTask, err := tsr.queries.CreateTasks(c, data)

	if err != nil {
		return models.Task{}, err
	}

	return models.DatabaseTaskToTask(newTask), nil
}

func (tsr *TaskRepository) GetTasks(c context.Context, filters interfaces.GetTasksFilters) ([]interfaces.GetTasksResponse, error) {

	sql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("t.id", "t.created_at", "t.updated_at", "t.project_id", "t.user_id", "t.status", "t.title", "t.description", "p.name", "u.username").From("tasks as t").LeftJoin("projects p ON p.id = t.project_id").LeftJoin("users u ON u.id = t.user_id")

	if filters.ProjectId != uuid.Nil {
		sql = sql.Where(sq.Eq{"t.project_id": filters.ProjectId})
	}

	if filters.UserId != uuid.Nil {
		sql = sql.Where(sq.Eq{"t.user_id": filters.UserId})
	}

	orderAsc := true

	// Handle cursor pagination
	if !filters.IsFirstPage {
		if filters.PointsNext {
			sql = sql.Where(sq.Or{
				sq.Gt{"t.created_at": filters.CursorCreatedAt},
				sq.And{
					sq.Eq{"t.created_at": filters.CursorCreatedAt},
					sq.Gt{"t.id": filters.CursorId},
				},
			})
			orderAsc = true
		} else {
			sql = sql.Where(sq.Or{
				sq.Lt{"t.created_at": filters.CursorCreatedAt},
				sq.And{
					sq.Eq{"t.created_at": filters.CursorCreatedAt},
					sq.Lt{"t.id": filters.CursorId},
				},
			})
			orderAsc = false
		}
	}

	if filters.Title != "" {
		titleLower := strings.ToLower(filters.Title)
		sql = sql.Where(sq.Like{"LOWER(t.title)": fmt.Sprintf("%%%v%%", titleLower)})
	}

	if orderAsc {
		sql = sql.OrderBy("t.created_at ASC, t.id ASC").Limit(filters.Limit + 1)
	} else {
		sql = sql.OrderBy("t.created_at DESC, t.id DESC").Limit(filters.Limit + 1)
	}

	queryString, arg, err := sql.ToSql()

	if err != nil {
		return []interfaces.GetTasksResponse{}, err
	}

	rows, err := tsr.db.Query(queryString, arg...)

	if err != nil {
		return []interfaces.GetTasksResponse{}, err
	}

	defer rows.Close()

	var tasks []interfaces.GetTasksResponse

	for rows.Next() {
		var task interfaces.GetTasksResponse
		var description, projectName, userName stdsql.NullString
		if err := rows.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt, &task.ProjectID, &task.UserID, &task.Status, &task.Title, &description, &projectName, &userName); err != nil {
			return tasks, err
		}

		if description.Valid {
			task.Description = description.String
		} else {
			task.Description = ""
		}

		if projectName.Valid {
			task.ProjectName = projectName.String
		} else {
			task.ProjectName = ""
		}

		if userName.Valid {
			task.UserName = userName.String
		} else {
			task.UserName = ""
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return []interfaces.GetTasksResponse{}, err
	}

	if !orderAsc && len(tasks) > 0 {
		for i, j := 0, len(tasks)-1; i < j; i, j = i+1, j-1 {
			tasks[i], tasks[j] = tasks[j], tasks[i]
		}
	}

	return tasks, nil

}

func (tsr *TaskRepository) GetTaskById(c context.Context, id uuid.UUID) (models.Task, error) {

	task, err := tsr.queries.GetTaskById(c, id)

	if err != nil {
		return models.Task{}, err
	}

	return models.DatabaseTaskToTask(task), nil

}

func (tsr *TaskRepository) UpdateTask(c context.Context, data interfaces.UpdateTaskData) (models.Task, error) {
	user, err := tsr.queries.UpdateTask(c, database.UpdateTaskParams{
		Title:       data.Title,
		UserID:      data.UserId,
		Status:      database.Taskstatus(data.Status),
		Description: data.Description,
		UpdatedAt:   data.UpdatedAt,
		ID:          data.ID,
	})

	if err != nil {
		return models.Task{}, err
	}

	return models.DatabaseTaskToTask(user), nil
}

func (tsr *TaskRepository) DeleteTask(c context.Context, id uuid.UUID) error {
	err := tsr.queries.DeleteTask(c, id)

	return err
}
