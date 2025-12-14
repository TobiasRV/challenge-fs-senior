package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/handlers"
	"github.com/TobiasRV/challenge-fs-senior/internals/interfaces"
	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/test/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateTask(t *testing.T) {
	userId := uuid.New()
	projectId := uuid.New()
	taskId := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		requestBody    map[string]interface{}
		setupMocks     func(*mocks.MockTaskRepository)
		expectedStatus int
	}{
		{
			name:     "Unauthorized - non-manager user",
			userRole: "Member",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"title":     "Test Task",
				"projectId": projectId.String(),
			},
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:     "Unauthorized - admin user",
			userRole: "Admin",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"title":     "Test Task",
				"projectId": projectId.String(),
			},
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:     "Missing title field",
			userRole: "Manager",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"projectId": projectId.String(),
			},
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:     "Missing projectId field",
			userRole: "Manager",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"title": "Test Task",
			},
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:     "Successfully create task without userId",
			userRole: "Manager",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"title":       "New Task",
				"description": "Task description",
				"projectId":   projectId.String(),
			},
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("CreateTask", mock.Anything, mock.AnythingOfType("database.CreateTasksParams")).Return(models.Task{
					ID:        taskId,
					Title:     "New Task",
					ProjectID: projectId,
					Status:    models.TaskstatusToDo,
					Description: sql.NullString{
						String: "Task description",
						Valid:  true,
					},
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
			},
			expectedStatus: fiber.StatusCreated,
		},
		{
			name:     "Successfully create task with userId",
			userRole: "Manager",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"title":       "Assigned Task",
				"description": "Task description",
				"projectId":   projectId.String(),
				"userId":      userId.String(),
			},
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("CreateTask", mock.Anything, mock.AnythingOfType("database.CreateTasksParams")).Return(models.Task{
					ID:        taskId,
					Title:     "Assigned Task",
					ProjectID: projectId,
					UserID: uuid.NullUUID{
						UUID:  userId,
						Valid: true,
					},
					Status: models.TaskstatusToDo,
					Description: sql.NullString{
						String: "Task description",
						Valid:  true,
					},
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
			},
			expectedStatus: fiber.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()

			mockUserRepo := mocks.NewMockUserRepository()
			mockRefreshTokenRepo := mocks.NewMockRefreshTokenRepository()
			mockTeamRepo := mocks.NewMockTeamRepository()
			mockProjectRepo := mocks.NewMockProjectRepository()
			mockTaskRepo := mocks.NewMockTaskRepository()

			tt.setupMocks(mockTaskRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Post("/tasks", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.CreateTask(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockTaskRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetTasks(t *testing.T) {
	userId := uuid.New()
	projectId := uuid.New()
	taskId := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		queryParams    string
		setupMocks     func(*mocks.MockTaskRepository)
		expectedStatus int
	}{
		{
			name:           "Admin without projectId - returns error",
			userRole:       "Admin",
			userId:         userId.String(),
			queryParams:    "?limit=10",
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "Manager without projectId - returns error",
			userRole:       "Manager",
			userId:         userId.String(),
			queryParams:    "?limit=10",
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:        "Member without projectId - returns own tasks",
			userRole:    "Member",
			userId:      userId.String(),
			queryParams: "?limit=10",
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("GetTasks", mock.Anything, mock.MatchedBy(func(filters interfaces.GetTasksFilters) bool {
					return filters.UserId == userId && filters.Limit == 10
				})).Return([]interfaces.GetTasksResponse{
					{
						ID:        taskId,
						Title:     "My Task",
						ProjectID: projectId,
						UserID: uuid.NullUUID{
							UUID:  userId,
							Valid: true,
						},
						Status:    models.TaskstatusToDo,
						CreatedAt: now,
						UpdatedAt: now,
					},
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:        "Successfully get tasks by projectId as Admin",
			userRole:    "Admin",
			userId:      userId.String(),
			queryParams: "?limit=10&projectId=" + projectId.String(),
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("GetTasks", mock.Anything, mock.MatchedBy(func(filters interfaces.GetTasksFilters) bool {
					return filters.ProjectId == projectId && filters.Limit == 10
				})).Return([]interfaces.GetTasksResponse{
					{
						ID:        taskId,
						Title:     "Test Task",
						ProjectID: projectId,
						Status:    models.TaskstatusToDo,
						CreatedAt: now,
						UpdatedAt: now,
					},
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:        "Successfully get tasks by projectId as Manager",
			userRole:    "Manager",
			userId:      userId.String(),
			queryParams: "?limit=10&projectId=" + projectId.String(),
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("GetTasks", mock.Anything, mock.MatchedBy(func(filters interfaces.GetTasksFilters) bool {
					return filters.ProjectId == projectId && filters.Limit == 10
				})).Return([]interfaces.GetTasksResponse{
					{
						ID:        taskId,
						Title:     "Test Task",
						ProjectID: projectId,
						Status:    models.TaskstatusInProgress,
						CreatedAt: now,
						UpdatedAt: now,
					},
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:        "Empty tasks list",
			userRole:    "Admin",
			userId:      userId.String(),
			queryParams: "?limit=10&projectId=" + projectId.String(),
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("GetTasks", mock.Anything, mock.Anything).Return([]interfaces.GetTasksResponse{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()

			mockUserRepo := mocks.NewMockUserRepository()
			mockRefreshTokenRepo := mocks.NewMockRefreshTokenRepository()
			mockTeamRepo := mocks.NewMockTeamRepository()
			mockProjectRepo := mocks.NewMockProjectRepository()
			mockTaskRepo := mocks.NewMockTaskRepository()

			tt.setupMocks(mockTaskRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Get("/tasks", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.GetTasks(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/tasks"+tt.queryParams, nil)

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockTaskRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_UpdateTask(t *testing.T) {
	userId := uuid.New()
	projectId := uuid.New()
	taskId := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		taskId         string
		requestBody    map[string]interface{}
		setupMocks     func(*mocks.MockTaskRepository)
		expectedStatus int
	}{
		{
			name:     "Unauthorized - non-manager user",
			userRole: "Member",
			userId:   userId.String(),
			taskId:   taskId.String(),
			requestBody: map[string]interface{}{
				"title":  "Updated Task",
				"status": "InProgress",
			},
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:     "Invalid task ID - empty",
			userRole: "Manager",
			userId:   userId.String(),
			taskId:   "",
			requestBody: map[string]interface{}{
				"title":  "Updated Task",
				"status": "InProgress",
			},
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusNotFound,
		},
		{
			name:     "Task not found",
			userRole: "Manager",
			userId:   userId.String(),
			taskId:   taskId.String(),
			requestBody: map[string]interface{}{
				"title":  "Updated Task",
				"status": "InProgress",
			},
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("GetTaskById", mock.Anything, taskId).Return(models.Task{}, sql.ErrNoRows)
			},
			expectedStatus: fiber.StatusNotFound,
		},
		{
			name:        "Missing required fields",
			userRole:    "Manager",
			userId:      userId.String(),
			taskId:      taskId.String(),
			requestBody: map[string]interface{}{},
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("GetTaskById", mock.Anything, taskId).Return(models.Task{
					ID:        taskId,
					Title:     "Test Task",
					ProjectID: projectId,
					Status:    models.TaskstatusToDo,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:     "Successfully update task",
			userRole: "Manager",
			userId:   userId.String(),
			taskId:   taskId.String(),
			requestBody: map[string]interface{}{
				"title":       "Updated Task",
				"description": "New description",
				"status":      "InProgress",
			},
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("GetTaskById", mock.Anything, taskId).Return(models.Task{
					ID:        taskId,
					Title:     "Test Task",
					ProjectID: projectId,
					Status:    models.TaskstatusToDo,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
				mockTaskRepo.On("UpdateTask", mock.Anything, mock.MatchedBy(func(data interfaces.UpdateTaskData) bool {
					return data.ID == taskId && data.Title == "Updated Task" && data.Status == "InProgress"
				})).Return(models.Task{
					ID:        taskId,
					Title:     "Updated Task",
					ProjectID: projectId,
					Status:    models.TaskstatusInProgress,
					Description: sql.NullString{
						String: "New description",
						Valid:  true,
					},
					CreatedAt: now,
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:     "Successfully update task with userId",
			userRole: "Manager",
			userId:   userId.String(),
			taskId:   taskId.String(),
			requestBody: map[string]interface{}{
				"title":  "Updated Task",
				"status": "Done",
				"userId": userId.String(),
			},
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("GetTaskById", mock.Anything, taskId).Return(models.Task{
					ID:        taskId,
					Title:     "Test Task",
					ProjectID: projectId,
					Status:    models.TaskstatusInProgress,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
				mockTaskRepo.On("UpdateTask", mock.Anything, mock.AnythingOfType("interfaces.UpdateTaskData")).Return(models.Task{
					ID:        taskId,
					Title:     "Updated Task",
					ProjectID: projectId,
					UserID: uuid.NullUUID{
						UUID:  userId,
						Valid: true,
					},
					Status:    models.TaskstatusDone,
					CreatedAt: now,
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()

			mockUserRepo := mocks.NewMockUserRepository()
			mockRefreshTokenRepo := mocks.NewMockRefreshTokenRepository()
			mockTeamRepo := mocks.NewMockTeamRepository()
			mockProjectRepo := mocks.NewMockProjectRepository()
			mockTaskRepo := mocks.NewMockTaskRepository()

			tt.setupMocks(mockTaskRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Put("/tasks/:id", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.UpdateTask(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			path := "/tasks/" + tt.taskId
			req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockTaskRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteTask(t *testing.T) {
	userId := uuid.New()
	taskId := uuid.New()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		taskId         string
		setupMocks     func(*mocks.MockTaskRepository)
		expectedStatus int
	}{
		{
			name:           "Unauthorized - non-manager user",
			userRole:       "Member",
			userId:         userId.String(),
			taskId:         taskId.String(),
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:           "Unauthorized - admin user",
			userRole:       "Admin",
			userId:         userId.String(),
			taskId:         taskId.String(),
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:           "Invalid task ID - empty",
			userRole:       "Manager",
			userId:         userId.String(),
			taskId:         "",
			setupMocks:     func(mockTaskRepo *mocks.MockTaskRepository) {},
			expectedStatus: fiber.StatusNotFound,
		},
		{
			name:     "Successfully delete task",
			userRole: "Manager",
			userId:   userId.String(),
			taskId:   taskId.String(),
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("DeleteTask", mock.Anything, taskId).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:     "Delete task - database error",
			userRole: "Manager",
			userId:   userId.String(),
			taskId:   taskId.String(),
			setupMocks: func(mockTaskRepo *mocks.MockTaskRepository) {
				mockTaskRepo.On("DeleteTask", mock.Anything, taskId).Return(sql.ErrConnDone)
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()

			mockUserRepo := mocks.NewMockUserRepository()
			mockRefreshTokenRepo := mocks.NewMockRefreshTokenRepository()
			mockTeamRepo := mocks.NewMockTeamRepository()
			mockProjectRepo := mocks.NewMockProjectRepository()
			mockTaskRepo := mocks.NewMockTaskRepository()

			tt.setupMocks(mockTaskRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.DeleteTask(c)
			})

			path := "/tasks/" + tt.taskId
			req := httptest.NewRequest(http.MethodDelete, path, nil)

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockTaskRepo.AssertExpectations(t)
		})
	}
}
