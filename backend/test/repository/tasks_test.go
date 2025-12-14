package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TobiasRV/challenge-fs-senior/internals/interfaces"
	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/repository"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_CreateTask(t *testing.T) {
	taskId := uuid.New()
	projectId := uuid.New()
	userId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name          string
		params        database.CreateTasksParams
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		expectedId    uuid.UUID
		expectedTitle string
	}{
		{
			name: "Successfully create task without user assignment",
			params: database.CreateTasksParams{
				Title:       "Test Task",
				ProjectID:   projectId,
				Description: sql.NullString{String: "Task description", Valid: true},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description"}).
					AddRow(taskId, now, now, projectId, nil, "ToDo", "Test Task", "Task description")
				// SQL: INSERT INTO tasks (created_at, updated_at, project_id, title, description, user_id)
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO tasks")).
					WithArgs(now, now, projectId, "Test Task", sql.NullString{String: "Task description", Valid: true}, uuid.NullUUID{}).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedId:    taskId,
			expectedTitle: "Test Task",
		},
		{
			name: "Successfully create task with user assignment",
			params: database.CreateTasksParams{
				Title:       "Assigned Task",
				ProjectID:   projectId,
				UserID:      uuid.NullUUID{UUID: userId, Valid: true},
				Description: sql.NullString{String: "Assigned task", Valid: true},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description"}).
					AddRow(taskId, now, now, projectId, userId, "ToDo", "Assigned Task", "Assigned task")
				// SQL: INSERT INTO tasks (created_at, updated_at, project_id, title, description, user_id)
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO tasks")).
					WithArgs(now, now, projectId, "Assigned Task", sql.NullString{String: "Assigned task", Valid: true}, uuid.NullUUID{UUID: userId, Valid: true}).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedId:    taskId,
			expectedTitle: "Assigned Task",
		},
		{
			name: "Database error on create",
			params: database.CreateTasksParams{
				Title:     "Test Task",
				ProjectID: projectId,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO tasks")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			queries := database.New(db)
			repo := repository.NewTaskRepository(queries, db)

			result, err := repo.CreateTask(context.Background(), tt.params)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedId, result.ID)
				assert.Equal(t, tt.expectedTitle, result.Title)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepository_GetTaskById(t *testing.T) {
	taskId := uuid.New()
	projectId := uuid.New()
	userId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name        string
		taskId      uuid.UUID
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectedId  uuid.UUID
	}{
		{
			name:   "Successfully get task by ID",
			taskId: taskId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description"}).
					AddRow(taskId, now, now, projectId, userId, "ToDo", "Test Task", "Task description")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, project_id, user_id, status, title, description FROM tasks")).
					WithArgs(taskId).
					WillReturnRows(rows)
			},
			expectError: false,
			expectedId:  taskId,
		},
		{
			name:   "Task not found",
			taskId: taskId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, project_id, user_id, status, title, description FROM tasks")).
					WithArgs(taskId).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
		},
		{
			name:   "Database error",
			taskId: taskId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, project_id, user_id, status, title, description FROM tasks")).
					WithArgs(taskId).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			queries := database.New(db)
			repo := repository.NewTaskRepository(queries, db)

			result, err := repo.GetTaskById(context.Background(), tt.taskId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedId, result.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepository_GetTasks(t *testing.T) {
	taskId := uuid.New()
	projectId := uuid.New()
	userId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name          string
		filters       interfaces.GetTasksFilters
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		expectedCount int
	}{
		{
			name: "Get tasks with project filter",
			filters: interfaces.GetTasksFilters{
				ProjectId:   projectId,
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description", "name", "username"}).
					AddRow(taskId, now, now, projectId, userId, "ToDo", "Test Task", "Description", "Project Name", "testuser")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get tasks with user filter",
			filters: interfaces.GetTasksFilters{
				UserId:      userId,
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description", "name", "username"}).
					AddRow(taskId, now, now, projectId, userId, "InProgress", "My Task", "Description", "Project Name", "testuser")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get tasks with title filter",
			filters: interfaces.GetTasksFilters{
				ProjectId:   projectId,
				Title:       "test",
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description", "name", "username"}).
					AddRow(taskId, now, now, projectId, userId, "ToDo", "Test Task", "Description", "Project Name", "testuser")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get tasks - empty result",
			filters: interfaces.GetTasksFilters{
				ProjectId:   uuid.New(),
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description", "name", "username"})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 0,
		},
		{
			name: "Get tasks with pagination - next page",
			filters: interfaces.GetTasksFilters{
				ProjectId:       projectId,
				Limit:           10,
				IsFirstPage:     false,
				PointsNext:      true,
				CursorCreatedAt: now.Add(-time.Hour),
				CursorId:        uuid.New(),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description", "name", "username"}).
					AddRow(taskId, now, now, projectId, userId, "ToDo", "Test Task", "Description", "Project Name", "testuser")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get tasks with pagination - previous page",
			filters: interfaces.GetTasksFilters{
				ProjectId:       projectId,
				Limit:           10,
				IsFirstPage:     false,
				PointsNext:      false,
				CursorCreatedAt: now.Add(time.Hour),
				CursorId:        uuid.New(),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description", "name", "username"}).
					AddRow(taskId, now, now, projectId, userId, "Done", "Test Task", "Description", "Project Name", "testuser")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Database error",
			filters: interfaces.GetTasksFilters{
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			queries := database.New(db)
			repo := repository.NewTaskRepository(queries, db)

			result, err := repo.GetTasks(context.Background(), tt.filters)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepository_UpdateTask(t *testing.T) {
	taskId := uuid.New()
	projectId := uuid.New()
	userId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name           string
		updateData     interfaces.UpdateTaskData
		mockSetup      func(sqlmock.Sqlmock)
		expectError    bool
		expectedTitle  string
		expectedStatus models.Taskstatus
	}{
		{
			name: "Successfully update task",
			updateData: interfaces.UpdateTaskData{
				ID:          taskId,
				Title:       "Updated Task",
				Description: sql.NullString{String: "Updated description", Valid: true},
				Status:      models.TaskstatusInProgress,
				UpdatedAt:   now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description"}).
					AddRow(taskId, now, now, projectId, userId, "InProgress", "Updated Task", "Updated description")
				// SQL: SET title = $1, description=$2 ,user_id=$3, status = $4, updated_at = $5 WHERE id = $6
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE tasks")).
					WithArgs("Updated Task", sql.NullString{String: "Updated description", Valid: true}, uuid.NullUUID{}, database.TaskstatusInProgress, now, taskId).
					WillReturnRows(rows)
			},
			expectError:    false,
			expectedTitle:  "Updated Task",
			expectedStatus: models.TaskstatusInProgress,
		},
		{
			name: "Update task to done",
			updateData: interfaces.UpdateTaskData{
				ID:        taskId,
				Title:     "Completed Task",
				Status:    models.TaskstatusDone,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description"}).
					AddRow(taskId, now, now, projectId, userId, "Done", "Completed Task", nil)
				// SQL: SET title = $1, description=$2 ,user_id=$3, status = $4, updated_at = $5 WHERE id = $6
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE tasks")).
					WithArgs("Completed Task", sql.NullString{}, uuid.NullUUID{}, database.TaskstatusDone, now, taskId).
					WillReturnRows(rows)
			},
			expectError:    false,
			expectedTitle:  "Completed Task",
			expectedStatus: models.TaskstatusDone,
		},
		{
			name: "Update task with user assignment",
			updateData: interfaces.UpdateTaskData{
				ID:        taskId,
				Title:     "Assigned Task",
				Status:    models.TaskstatusToDo,
				UserId:    uuid.NullUUID{UUID: userId, Valid: true},
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "project_id", "user_id", "status", "title", "description"}).
					AddRow(taskId, now, now, projectId, userId, "ToDo", "Assigned Task", nil)
				// SQL: SET title = $1, description=$2 ,user_id=$3, status = $4, updated_at = $5 WHERE id = $6
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE tasks")).
					WithArgs("Assigned Task", sql.NullString{}, uuid.NullUUID{UUID: userId, Valid: true}, database.TaskstatusToDo, now, taskId).
					WillReturnRows(rows)
			},
			expectError:    false,
			expectedTitle:  "Assigned Task",
			expectedStatus: models.TaskstatusToDo,
		},
		{
			name: "Update task - not found",
			updateData: interfaces.UpdateTaskData{
				ID:        taskId,
				Title:     "Updated Task",
				Status:    models.TaskstatusInProgress,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// SQL: SET title = $1, description=$2 ,user_id=$3, status = $4, updated_at = $5 WHERE id = $6
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE tasks")).
					WithArgs("Updated Task", sql.NullString{}, uuid.NullUUID{}, database.TaskstatusInProgress, now, taskId).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
		},
		{
			name: "Update task - database error",
			updateData: interfaces.UpdateTaskData{
				ID:        taskId,
				Title:     "Updated Task",
				Status:    models.TaskstatusInProgress,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// SQL: SET title = $1, description=$2 ,user_id=$3, status = $4, updated_at = $5 WHERE id = $6
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE tasks")).
					WithArgs("Updated Task", sql.NullString{}, uuid.NullUUID{}, database.TaskstatusInProgress, now, taskId).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			queries := database.New(db)
			repo := repository.NewTaskRepository(queries, db)

			result, err := repo.UpdateTask(context.Background(), tt.updateData)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTitle, result.Title)
				assert.Equal(t, tt.expectedStatus, result.Status)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepository_DeleteTask(t *testing.T) {
	taskId := uuid.New()

	tests := []struct {
		name        string
		taskId      uuid.UUID
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:   "Successfully delete task",
			taskId: taskId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM tasks")).
					WithArgs(taskId).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name:   "Delete task - database error",
			taskId: taskId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM tasks")).
					WithArgs(taskId).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			queries := database.New(db)
			repo := repository.NewTaskRepository(queries, db)

			err = repo.DeleteTask(context.Background(), tt.taskId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
