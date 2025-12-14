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

func TestProjectRepository_CreateProject(t *testing.T) {
	projectId := uuid.New()
	teamId := uuid.New()
	managerId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name         string
		params       database.CreateProjectParams
		mockSetup    func(sqlmock.Sqlmock)
		expectError  bool
		expectedId   uuid.UUID
		expectedName string
	}{
		{
			name: "Successfully create project",
			params: database.CreateProjectParams{
				Name:      "Test Project",
				TeamID:    teamId,
				ManagerID: managerId,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"}).
					AddRow(projectId, now, now, "Test Project", teamId, managerId, "OnHold")
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO projects")).
					WithArgs(now, now, "Test Project", teamId, managerId).
					WillReturnRows(rows)
			},
			expectError:  false,
			expectedId:   projectId,
			expectedName: "Test Project",
		},
		{
			name: "Database error on create",
			params: database.CreateProjectParams{
				Name:      "Test Project",
				TeamID:    teamId,
				ManagerID: managerId,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO projects")).
					WithArgs(now, now, "Test Project", teamId, managerId).
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
			repo := repository.NewProjectRepository(queries, db)

			result, err := repo.CreateProject(context.Background(), tt.params)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedId, result.ID)
				assert.Equal(t, tt.expectedName, result.Name)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProjectRepository_GetProjectById(t *testing.T) {
	projectId := uuid.New()
	teamId := uuid.New()
	managerId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name        string
		projectId   uuid.UUID
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectedId  uuid.UUID
	}{
		{
			name:      "Successfully get project by ID",
			projectId: projectId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"}).
					AddRow(projectId, now, now, "Test Project", teamId, managerId, "OnHold")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, name, team_id, manager_id, status FROM projects")).
					WithArgs(projectId).
					WillReturnRows(rows)
			},
			expectError: false,
			expectedId:  projectId,
		},
		{
			name:      "Project not found",
			projectId: projectId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, name, team_id, manager_id, status FROM projects")).
					WithArgs(projectId).
					WillReturnError(sql.ErrNoRows)
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
			repo := repository.NewProjectRepository(queries, db)

			result, err := repo.GetProjectById(context.Background(), tt.projectId)

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

func TestProjectRepository_GetProjectByManager(t *testing.T) {
	projectId := uuid.New()
	teamId := uuid.New()
	managerId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name        string
		managerId   uuid.UUID
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectedId  uuid.UUID
	}{
		{
			name:      "Successfully get project by manager",
			managerId: managerId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"}).
					AddRow(projectId, now, now, "Test Project", teamId, managerId, "InProgress")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, name, team_id, manager_id, status FROM projects")).
					WithArgs(managerId).
					WillReturnRows(rows)
			},
			expectError: false,
			expectedId:  projectId,
		},
		{
			name:      "Project not found for manager",
			managerId: managerId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, name, team_id, manager_id, status FROM projects")).
					WithArgs(managerId).
					WillReturnError(sql.ErrNoRows)
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
			repo := repository.NewProjectRepository(queries, db)

			result, err := repo.GetProjectByManager(context.Background(), tt.managerId)

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

func TestProjectRepository_GetProjects(t *testing.T) {
	projectId := uuid.New()
	teamId := uuid.New()
	managerId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name          string
		filters       interfaces.GetProjectsFilters
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		expectedCount int
	}{
		{
			name: "Get projects with team filter",
			filters: interfaces.GetProjectsFilters{
				TeamId:      teamId,
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"}).
					AddRow(projectId, now, now, "Test Project", teamId, managerId, "OnHold")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get projects with manager filter",
			filters: interfaces.GetProjectsFilters{
				ManagerId:   managerId,
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"}).
					AddRow(projectId, now, now, "Test Project", teamId, managerId, "InProgress")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get projects with stats",
			filters: interfaces.GetProjectsFilters{
				TeamId:      teamId,
				Limit:       10,
				IsFirstPage: true,
				WithStats:   true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status", "ToDoTasks", "InProgressTasks", "DoneTasks"}).
					AddRow(projectId, now, now, "Test Project", teamId, managerId, "OnHold", 5, 3, 2)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get projects with name filter",
			filters: interfaces.GetProjectsFilters{
				TeamId:      teamId,
				Name:        "test",
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"}).
					AddRow(projectId, now, now, "Test Project", teamId, managerId, "OnHold")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get projects - empty result",
			filters: interfaces.GetProjectsFilters{
				TeamId:      uuid.New(),
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 0,
		},
		{
			name: "Get projects with pagination - next page",
			filters: interfaces.GetProjectsFilters{
				TeamId:          teamId,
				Limit:           10,
				IsFirstPage:     false,
				PointsNext:      true,
				CursorCreatedAt: now.Add(-time.Hour),
				CursorId:        uuid.New(),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"}).
					AddRow(projectId, now, now, "Test Project", teamId, managerId, "OnHold")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Database error",
			filters: interfaces.GetProjectsFilters{
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
			repo := repository.NewProjectRepository(queries, db)

			result, err := repo.GetProjects(context.Background(), tt.filters)

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

func TestProjectRepository_UpdateProject(t *testing.T) {
	projectId := uuid.New()
	teamId := uuid.New()
	managerId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name           string
		updateData     interfaces.UpdateProjectData
		mockSetup      func(sqlmock.Sqlmock)
		expectError    bool
		expectedName   string
		expectedStatus models.Projectstatus
	}{
		{
			name: "Successfully update project",
			updateData: interfaces.UpdateProjectData{
				ID:        projectId,
				Name:      "Updated Project",
				Status:    models.ProjectstatusInProgress,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"}).
					AddRow(projectId, now, now, "Updated Project", teamId, managerId, "InProgress")
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE projects")).
					WithArgs("Updated Project", database.ProjectstatusInProgress, now, projectId).
					WillReturnRows(rows)
			},
			expectError:    false,
			expectedName:   "Updated Project",
			expectedStatus: models.ProjectstatusInProgress,
		},
		{
			name: "Update project to completed",
			updateData: interfaces.UpdateProjectData{
				ID:        projectId,
				Name:      "Completed Project",
				Status:    models.ProjectstatusCompleted,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "team_id", "manager_id", "status"}).
					AddRow(projectId, now, now, "Completed Project", teamId, managerId, "Completed")
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE projects")).
					WithArgs("Completed Project", database.ProjectstatusCompleted, now, projectId).
					WillReturnRows(rows)
			},
			expectError:    false,
			expectedName:   "Completed Project",
			expectedStatus: models.ProjectstatusCompleted,
		},
		{
			name: "Update project - not found",
			updateData: interfaces.UpdateProjectData{
				ID:        projectId,
				Name:      "Updated Project",
				Status:    models.ProjectstatusInProgress,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE projects")).
					WithArgs("Updated Project", database.ProjectstatusInProgress, now, projectId).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
		},
		{
			name: "Update project - database error",
			updateData: interfaces.UpdateProjectData{
				ID:        projectId,
				Name:      "Updated Project",
				Status:    models.ProjectstatusInProgress,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE projects")).
					WithArgs("Updated Project", database.ProjectstatusInProgress, now, projectId).
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
			repo := repository.NewProjectRepository(queries, db)

			result, err := repo.UpdateProject(context.Background(), tt.updateData)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedName, result.Name)
				assert.Equal(t, tt.expectedStatus, result.Status)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProjectRepository_DeleteProject(t *testing.T) {
	projectId := uuid.New()

	tests := []struct {
		name        string
		projectId   uuid.UUID
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:      "Successfully delete project",
			projectId: projectId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM projects")).
					WithArgs(projectId).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name:      "Delete project - database error",
			projectId: projectId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM projects")).
					WithArgs(projectId).
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
			repo := repository.NewProjectRepository(queries, db)

			err = repo.DeleteProject(context.Background(), tt.projectId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
