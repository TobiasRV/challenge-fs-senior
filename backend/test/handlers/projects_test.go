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

func TestHandler_CreateProject(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	projectId := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		requestBody    map[string]interface{}
		setupMocks     func(*mocks.MockUserRepository, *mocks.MockProjectRepository)
		expectedStatus int
	}{
		{
			name:     "Unauthorized - non-manager user",
			userRole: "Member",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"name": "Test Project",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockProjectRepo *mocks.MockProjectRepository) {
			},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:     "Unauthorized - admin user",
			userRole: "Admin",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"name": "Test Project",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockProjectRepo *mocks.MockProjectRepository) {
			},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:        "Missing name field",
			userRole:    "Manager",
			userId:      userId.String(),
			requestBody: map[string]interface{}{},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockProjectRepo *mocks.MockProjectRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:     "Manager not found",
			userRole: "Manager",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"name": "Test Project",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockProjectRepo *mocks.MockProjectRepository) {
				mockUserRepo.On("GetUserById", mock.Anything, userId).Return(models.User{}, sql.ErrNoRows)
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:     "Successfully create project",
			userRole: "Manager",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"name": "New Project",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockProjectRepo *mocks.MockProjectRepository) {
				mockUserRepo.On("GetUserById", mock.Anything, userId).Return(models.User{
					ID:       userId,
					Username: "manager",
					Email:    "manager@example.com",
					Role:     models.UserrolesManager,
					TeamId:   teamId,
				}, nil)
				mockProjectRepo.On("CreateProject", mock.Anything, mock.AnythingOfType("database.CreateProjectParams")).Return(models.Project{
					ID:        projectId,
					Name:      "New Project",
					TeamID:    teamId,
					ManagerID: userId,
					Status:    models.ProjectstatusOnHold,
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

			tt.setupMocks(mockUserRepo, mockProjectRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Post("/projects", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.CreateProject(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockUserRepo.AssertExpectations(t)
			mockProjectRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetProjects(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	projectId := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		queryParams    string
		setupMocks     func(*mocks.MockProjectRepository)
		expectedStatus int
	}{
		{
			name:           "Unauthorized - member user",
			userRole:       "Member",
			userId:         userId.String(),
			queryParams:    "?limit=10&teamId=" + teamId.String(),
			setupMocks:     func(mockProjectRepo *mocks.MockProjectRepository) {},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:           "Missing teamId and managerId",
			userRole:       "Admin",
			userId:         userId.String(),
			queryParams:    "?limit=10",
			setupMocks:     func(mockProjectRepo *mocks.MockProjectRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:        "Successfully get projects by teamId as Admin",
			userRole:    "Admin",
			userId:      userId.String(),
			queryParams: "?limit=10&teamId=" + teamId.String(),
			setupMocks: func(mockProjectRepo *mocks.MockProjectRepository) {
				mockProjectRepo.On("GetProjects", mock.Anything, mock.MatchedBy(func(filters interfaces.GetProjectsFilters) bool {
					return filters.TeamId == teamId && filters.Limit == 10
				})).Return([]interfaces.GetProjectsResponse{
					{
						ID:        projectId,
						Name:      "Test Project",
						TeamID:    teamId,
						ManagerID: userId,
						Status:    models.ProjectstatusOnHold,
						CreatedAt: now,
						UpdatedAt: now,
					},
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:        "Successfully get projects by managerId as Manager",
			userRole:    "Manager",
			userId:      userId.String(),
			queryParams: "?limit=10&managerId=" + userId.String(),
			setupMocks: func(mockProjectRepo *mocks.MockProjectRepository) {
				mockProjectRepo.On("GetProjects", mock.Anything, mock.MatchedBy(func(filters interfaces.GetProjectsFilters) bool {
					return filters.ManagerId == userId && filters.Limit == 10
				})).Return([]interfaces.GetProjectsResponse{
					{
						ID:        projectId,
						Name:      "Test Project",
						TeamID:    teamId,
						ManagerID: userId,
						Status:    models.ProjectstatusInProgress,
						CreatedAt: now,
						UpdatedAt: now,
					},
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:        "Empty projects list",
			userRole:    "Admin",
			userId:      userId.String(),
			queryParams: "?limit=10&teamId=" + teamId.String(),
			setupMocks: func(mockProjectRepo *mocks.MockProjectRepository) {
				mockProjectRepo.On("GetProjects", mock.Anything, mock.Anything).Return([]interfaces.GetProjectsResponse{}, nil)
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

			tt.setupMocks(mockProjectRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Get("/projects", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.GetProjects(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/projects"+tt.queryParams, nil)

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockProjectRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_UpdateProject(t *testing.T) {
	userId := uuid.New()
	otherUserId := uuid.New()
	teamId := uuid.New()
	projectId := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		projectId      string
		requestBody    map[string]interface{}
		setupMocks     func(*mocks.MockProjectRepository)
		expectedStatus int
	}{
		{
			name:      "Unauthorized - non-manager user",
			userRole:  "Member",
			userId:    userId.String(),
			projectId: projectId.String(),
			requestBody: map[string]interface{}{
				"name":   "Updated Project",
				"status": "InProgress",
			},
			setupMocks:     func(mockProjectRepo *mocks.MockProjectRepository) {},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:      "Invalid project ID",
			userRole:  "Manager",
			userId:    userId.String(),
			projectId: "invalid-uuid",
			requestBody: map[string]interface{}{
				"name":   "Updated Project",
				"status": "InProgress",
			},
			setupMocks:     func(mockProjectRepo *mocks.MockProjectRepository) {},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name:      "Project not found",
			userRole:  "Manager",
			userId:    userId.String(),
			projectId: projectId.String(),
			requestBody: map[string]interface{}{
				"name":   "Updated Project",
				"status": "InProgress",
			},
			setupMocks: func(mockProjectRepo *mocks.MockProjectRepository) {
				mockProjectRepo.On("GetProjectById", mock.Anything, projectId).Return(models.Project{}, sql.ErrNoRows)
			},
			expectedStatus: fiber.StatusNotFound,
		},
		{
			name:      "Unauthorized - not project manager",
			userRole:  "Manager",
			userId:    userId.String(),
			projectId: projectId.String(),
			requestBody: map[string]interface{}{
				"name":   "Updated Project",
				"status": "InProgress",
			},
			setupMocks: func(mockProjectRepo *mocks.MockProjectRepository) {
				mockProjectRepo.On("GetProjectById", mock.Anything, projectId).Return(models.Project{
					ID:        projectId,
					Name:      "Test Project",
					TeamID:    teamId,
					ManagerID: otherUserId, // Different manager
					Status:    models.ProjectstatusOnHold,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
			},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:        "Missing required fields",
			userRole:    "Manager",
			userId:      userId.String(),
			projectId:   projectId.String(),
			requestBody: map[string]interface{}{
				// Missing name and status
			},
			setupMocks: func(mockProjectRepo *mocks.MockProjectRepository) {
				mockProjectRepo.On("GetProjectById", mock.Anything, projectId).Return(models.Project{
					ID:        projectId,
					Name:      "Test Project",
					TeamID:    teamId,
					ManagerID: userId,
					Status:    models.ProjectstatusOnHold,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:      "Successfully update project",
			userRole:  "Manager",
			userId:    userId.String(),
			projectId: projectId.String(),
			requestBody: map[string]interface{}{
				"name":   "Updated Project",
				"status": "InProgress",
			},
			setupMocks: func(mockProjectRepo *mocks.MockProjectRepository) {
				mockProjectRepo.On("GetProjectById", mock.Anything, projectId).Return(models.Project{
					ID:        projectId,
					Name:      "Test Project",
					TeamID:    teamId,
					ManagerID: userId,
					Status:    models.ProjectstatusOnHold,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
				mockProjectRepo.On("UpdateProject", mock.Anything, mock.MatchedBy(func(data interfaces.UpdateProjectData) bool {
					return data.ID == projectId && data.Name == "Updated Project" && data.Status == "InProgress"
				})).Return(models.Project{
					ID:        projectId,
					Name:      "Updated Project",
					TeamID:    teamId,
					ManagerID: userId,
					Status:    models.ProjectstatusInProgress,
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

			tt.setupMocks(mockProjectRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Put("/projects/:id", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.UpdateProject(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			path := "/projects/" + tt.projectId
			req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockProjectRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteProject(t *testing.T) {
	userId := uuid.New()
	projectId := uuid.New()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		projectId      string
		setupMocks     func(*mocks.MockProjectRepository)
		expectedStatus int
	}{
		{
			name:           "Unauthorized - non-manager user",
			userRole:       "Member",
			userId:         userId.String(),
			projectId:      projectId.String(),
			setupMocks:     func(mockProjectRepo *mocks.MockProjectRepository) {},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:           "Invalid project ID - invalid uuid",
			userRole:       "Manager",
			userId:         userId.String(),
			projectId:      "invalid-uuid",
			setupMocks:     func(mockProjectRepo *mocks.MockProjectRepository) {},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name:      "Successfully delete project",
			userRole:  "Manager",
			userId:    userId.String(),
			projectId: projectId.String(),
			setupMocks: func(mockProjectRepo *mocks.MockProjectRepository) {
				mockProjectRepo.On("DeleteProject", mock.Anything, projectId).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:      "Delete project - database error",
			userRole:  "Manager",
			userId:    userId.String(),
			projectId: projectId.String(),
			setupMocks: func(mockProjectRepo *mocks.MockProjectRepository) {
				mockProjectRepo.On("DeleteProject", mock.Anything, projectId).Return(sql.ErrConnDone)
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

			tt.setupMocks(mockProjectRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Delete("/projects/:id", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.DeleteProject(c)
			})

			path := "/projects/" + tt.projectId
			req := httptest.NewRequest(http.MethodDelete, path, nil)

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockProjectRepo.AssertExpectations(t)
		})
	}
}
