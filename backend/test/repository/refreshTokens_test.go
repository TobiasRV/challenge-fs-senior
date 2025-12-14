package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/repository"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRefreshTokenRepository_CreateRefreshToken(t *testing.T) {
	tokenId := uuid.New()
	userId := uuid.New()
	now := time.Now().UTC()
	expiresAt := now.Add(7 * 24 * time.Hour)

	tests := []struct {
		name        string
		token       models.RefreshToken
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successfully create refresh token",
			token: models.RefreshToken{
				ID:        tokenId,
				Userid:    userId,
				Token:     "test-refresh-token-abc123",
				ExpiresAt: expiresAt,
				CreatedAt: now,
				Revoked:   false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO refresh_tokens")).
					WithArgs(tokenId, userId, "test-refresh-token-abc123", expiresAt, now, false).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name: "Create refresh token - database error",
			token: models.RefreshToken{
				ID:        tokenId,
				Userid:    userId,
				Token:     "test-refresh-token-def456",
				ExpiresAt: expiresAt,
				CreatedAt: now,
				Revoked:   false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO refresh_tokens")).
					WithArgs(tokenId, userId, "test-refresh-token-def456", expiresAt, now, false).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
		{
			name: "Create refresh token - duplicate token error",
			token: models.RefreshToken{
				ID:        tokenId,
				Userid:    userId,
				Token:     "duplicate-token",
				ExpiresAt: expiresAt,
				CreatedAt: now,
				Revoked:   false,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO refresh_tokens")).
					WithArgs(tokenId, userId, "duplicate-token", expiresAt, now, false).
					WillReturnError(sql.ErrNoRows) // simulating constraint error
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
			repo := repository.NewRefreshTokenRepository(queries, db)

			err = repo.CreateRefreshToken(context.Background(), tt.token)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRefreshTokenRepository_GetRefreshTokenByToken(t *testing.T) {
	tokenId := uuid.New()
	userId := uuid.New()
	now := time.Now().UTC()
	expiresAt := now.Add(7 * 24 * time.Hour)

	tests := []struct {
		name           string
		token          string
		mockSetup      func(sqlmock.Sqlmock)
		expectError    bool
		expectedUserId uuid.UUID
		expectedToken  string
	}{
		{
			name:  "Successfully get refresh token with user data",
			token: "valid-refresh-token",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "userid", "token", "expires_at", "created_at", "revoked",
					"userdataid", "userdatacreatedat", "userdataupdatedat",
					"userdatausername", "userdatapassword", "userdataemail", "userdatarole",
				}).
					AddRow(
						tokenId, userId, "valid-refresh-token", expiresAt, now, false,
						userId, now, now, "testuser", "hashedpassword", "test@example.com", "user",
					)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT refresh_tokens.id")).
					WithArgs("valid-refresh-token").
					WillReturnRows(rows)
			},
			expectError:    false,
			expectedUserId: userId,
			expectedToken:  "valid-refresh-token",
		},
		{
			name:  "Get refresh token - token not found",
			token: "nonexistent-token",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT refresh_tokens.id")).
					WithArgs("nonexistent-token").
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
		},
		{
			name:  "Get refresh token - database error",
			token: "error-token",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT refresh_tokens.id")).
					WithArgs("error-token").
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
		{
			name:  "Get revoked refresh token",
			token: "revoked-refresh-token",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "userid", "token", "expires_at", "created_at", "revoked",
					"userdataid", "userdatacreatedat", "userdataupdatedat",
					"userdatausername", "userdatapassword", "userdataemail", "userdatarole",
				}).
					AddRow(
						tokenId, userId, "revoked-refresh-token", expiresAt, now, true,
						userId, now, now, "testuser", "hashedpassword", "test@example.com", "user",
					)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT refresh_tokens.id")).
					WithArgs("revoked-refresh-token").
					WillReturnRows(rows)
			},
			expectError:    false,
			expectedUserId: userId,
			expectedToken:  "revoked-refresh-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			queries := database.New(db)
			repo := repository.NewRefreshTokenRepository(queries, db)

			result, err := repo.GetRefreshTokenByToken(context.Background(), tt.token)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserId, result.Userid)
				assert.Equal(t, tt.expectedToken, result.Token)
				assert.NotEmpty(t, result.UserData.Username)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRefreshTokenRepository_DeleteRefreshTokensByUserId(t *testing.T) {
	userId := uuid.New()

	tests := []struct {
		name        string
		userId      string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:   "Successfully delete refresh tokens by user ID",
			userId: userId.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM refresh_tokens")).
					WithArgs(userId).
					WillReturnResult(sqlmock.NewResult(0, 3)) // Assume 3 tokens deleted
			},
			expectError: false,
		},
		{
			name:   "Delete refresh tokens - no tokens to delete",
			userId: userId.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM refresh_tokens")).
					WithArgs(userId).
					WillReturnResult(sqlmock.NewResult(0, 0)) // No tokens deleted
			},
			expectError: false,
		},
		{
			name:        "Delete refresh tokens - invalid UUID",
			userId:      "invalid-uuid",
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
		},
		{
			name:   "Delete refresh tokens - database error",
			userId: userId.String(),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM refresh_tokens")).
					WithArgs(userId).
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
			repo := repository.NewRefreshTokenRepository(queries, db)

			err = repo.DeleteRefreshTokensByUserId(context.Background(), tt.userId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
