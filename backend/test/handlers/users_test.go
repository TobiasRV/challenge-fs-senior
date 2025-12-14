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

func TestHandler_CreateUserAdmin(t *testing.T) {
	userId := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		setupMocks     func(*mocks.MockUserRepository, *mocks.MockRefreshTokenRepository)
		expectedStatus int
	}{
		{
			name: "Missing username field",
			requestBody: map[string]interface{}{
				"password": "password123",
				"email":    "test@example.com",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Missing password field",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Missing email field",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Invalid email format",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
				"email":    "invalid-email",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "User already exists",
			requestBody: map[string]interface{}{
				"username": "existinguser",
				"password": "password123",
				"email":    "existing@example.com",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				mockUserRepo.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(models.User{
					ID:       userId,
					Email:    "existing@example.com",
					Username: "existinguser",
				}, nil)
			},
			expectedStatus: fiber.StatusConflict,
		},
		{
			name: "Successfully create admin user",
			requestBody: map[string]interface{}{
				"username": "newadmin",
				"password": "password123",
				"email":    "newadmin@example.com",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				mockUserRepo.On("GetUserByEmail", mock.Anything, "newadmin@example.com").Return(models.User{}, sql.ErrNoRows)
				mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("models.User")).Return(models.User{
					ID:        userId,
					Username:  "newadmin",
					Email:     "newadmin@example.com",
					Role:      models.UserrolesAdmin,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
				mockRefreshTokenRepo.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("models.RefreshToken")).Return(nil)
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

			tt.setupMocks(mockUserRepo, mockRefreshTokenRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Post("/auth/register-admin", handler.CreateUserAdmin)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/register-admin", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockUserRepo.AssertExpectations(t)
			mockRefreshTokenRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_UserExistsByEmail(t *testing.T) {
	userId := uuid.New()

	tests := []struct {
		name           string
		email          string
		setupMocks     func(*mocks.MockUserRepository)
		expectedStatus int
	}{
		{
			name:  "Missing email parameter",
			email: "",
			setupMocks: func(mockUserRepo *mocks.MockUserRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:  "User exists",
			email: "existing@example.com",
			setupMocks: func(mockUserRepo *mocks.MockUserRepository) {
				mockUserRepo.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(models.User{
					ID:    userId,
					Email: "existing@example.com",
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:  "User does not exist",
			email: "notfound@example.com",
			setupMocks: func(mockUserRepo *mocks.MockUserRepository) {
				mockUserRepo.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return(models.User{}, sql.ErrNoRows)
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

			tt.setupMocks(mockUserRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Get("/users/exists-by-email", handler.UserExistsByEmail)

			url := "/users/exists-by-email"
			if tt.email != "" {
				url += "?email=" + tt.email
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockUserRepo.AssertExpectations(t)
		})
	}
}
