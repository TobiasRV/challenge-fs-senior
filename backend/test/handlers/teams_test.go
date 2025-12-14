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
	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/test/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateTeam(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		requestBody    map[string]interface{}
		setupMocks     func(*mocks.MockUserRepository, *mocks.MockTeamRepository)
		expectedStatus int
	}{
		{
			name:     "Unauthorized - non-admin user",
			userRole: "Member",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"name": "Test Team",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockTeamRepo *mocks.MockTeamRepository) {
			},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:        "Missing name field",
			userRole:    "Admin",
			userId:      userId.String(),
			requestBody: map[string]interface{}{},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockTeamRepo *mocks.MockTeamRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:     "User not found",
			userRole: "Admin",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"name": "Test Team",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockTeamRepo *mocks.MockTeamRepository) {
				mockUserRepo.On("GetUserById", mock.Anything, userId).Return(models.User{}, sql.ErrNoRows)
			},
			expectedStatus: fiber.StatusNotFound,
		},
		{
			name:     "User already has a team",
			userRole: "Admin",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"name": "Test Team",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockTeamRepo *mocks.MockTeamRepository) {
				mockUserRepo.On("GetUserById", mock.Anything, userId).Return(models.User{
					ID:       userId,
					Username: "admin",
					Email:    "admin@example.com",
					Role:     models.UserrolesAdmin,
				}, nil)
				mockTeamRepo.On("GetTeamByOwner", mock.Anything, userId).Return(true, models.Team{
					ID:      teamId,
					Name:    "Existing Team",
					OwnerID: userId,
				}, nil)
			},
			expectedStatus: fiber.StatusConflict,
		},
		{
			name:     "Successfully create team",
			userRole: "Admin",
			userId:   userId.String(),
			requestBody: map[string]interface{}{
				"name": "New Team",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockTeamRepo *mocks.MockTeamRepository) {
				mockUserRepo.On("GetUserById", mock.Anything, userId).Return(models.User{
					ID:       userId,
					Username: "admin",
					Email:    "admin@example.com",
					Role:     models.UserrolesAdmin,
				}, nil)
				mockTeamRepo.On("GetTeamByOwner", mock.Anything, userId).Return(false, models.Team{}, nil)
				mockTeamRepo.On("CreateTeam", mock.Anything, mock.AnythingOfType("models.Team")).Return(models.Team{
					ID:        teamId,
					Name:      "New Team",
					OwnerID:   userId,
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

			tt.setupMocks(mockUserRepo, mockTeamRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Post("/teams", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.CreateTeam(c)
			})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockUserRepo.AssertExpectations(t)
			mockTeamRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetTeamByOwner(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		userRole       string
		userId         string
		setupMocks     func(*mocks.MockTeamRepository)
		expectedStatus int
	}{
		{
			name:     "Unauthorized - non-admin user",
			userRole: "Member",
			userId:   userId.String(),
			setupMocks: func(mockTeamRepo *mocks.MockTeamRepository) {
			},
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:     "Team not found - returns exists: false",
			userRole: "Admin",
			userId:   userId.String(),
			setupMocks: func(mockTeamRepo *mocks.MockTeamRepository) {
				mockTeamRepo.On("GetTeamByOwner", mock.Anything, userId).Return(false, models.Team{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:     "Successfully get team",
			userRole: "Admin",
			userId:   userId.String(),
			setupMocks: func(mockTeamRepo *mocks.MockTeamRepository) {
				mockTeamRepo.On("GetTeamByOwner", mock.Anything, userId).Return(true, models.Team{
					ID:        teamId,
					Name:      "Test Team",
					OwnerID:   userId,
					CreatedAt: now,
					UpdatedAt: now,
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

			tt.setupMocks(mockTeamRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Get("/teams/by-owner", func(c *fiber.Ctx) error {
				c.Locals("userRole", tt.userRole)
				c.Locals("userId", tt.userId)
				return handler.GetTeamByOwner(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/teams/by-owner", nil)

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockTeamRepo.AssertExpectations(t)
		})
	}
}
