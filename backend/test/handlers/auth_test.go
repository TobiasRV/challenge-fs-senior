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
	"golang.org/x/crypto/bcrypt"
)

// setupTestApp creates a new Fiber app for testing
func setupTestApp() *fiber.App {
	return fiber.New()
}

func TestHandler_LogIn(t *testing.T) {
	userId := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		setupMocks     func(*mocks.MockUserRepository, *mocks.MockRefreshTokenRepository)
		expectedStatus int
	}{
		{
			name: "Missing email field",
			requestBody: map[string]interface{}{
				"password": "password123",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				// No mocks needed for validation failure
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Missing password field",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				// No mocks needed for validation failure
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "User not found",
			requestBody: map[string]interface{}{
				"email":    "notfound@example.com",
				"password": "password123",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				mockUserRepo.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return(models.User{}, sql.ErrNoRows)
			},
			expectedStatus: fiber.StatusConflict,
		},
		{
			name: "Invalid password",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				mockUserRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(models.User{
					ID:       userId,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedStatus: fiber.StatusConflict,
		},
		{
			name: "Successful login",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			setupMocks: func(mockUserRepo *mocks.MockUserRepository, mockRefreshTokenRepo *mocks.MockRefreshTokenRepository) {
				mockUserRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(models.User{
					ID:        userId,
					Email:     "test@example.com",
					Password:  string(hashedPassword),
					Username:  "testuser",
					Role:      models.UserrolesMember,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
				mockRefreshTokenRepo.On("DeleteRefreshTokensByUserId", mock.Anything, userId.String()).Return(nil)
				mockRefreshTokenRepo.On("CreateRefreshToken", mock.Anything, mock.AnythingOfType("models.RefreshToken")).Return(nil)
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

			tt.setupMocks(mockUserRepo, mockRefreshTokenRepo)

			handler := handlers.NewHandler(mockUserRepo, mockRefreshTokenRepo, mockTeamRepo, mockProjectRepo, mockTaskRepo)

			app.Post("/auth/login", handler.LogIn)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockUserRepo.AssertExpectations(t)
			mockRefreshTokenRepo.AssertExpectations(t)
		})
	}
}
