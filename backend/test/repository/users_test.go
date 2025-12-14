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

func TestUserRepository_CreateUser(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name        string
		userData    models.User
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectedId  uuid.UUID
	}{
		{
			name: "Successfully create user with team",
			userData: models.User{
				Username:  "testuser",
				Password:  "hashedpassword",
				Email:     "test@example.com",
				Role:      models.UserrolesMember,
				TeamId:    teamId,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "testuser", "hashedpassword", "test@example.com", "Member", teamId)
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users")).
					WithArgs(now, now, "testuser", "hashedpassword", "test@example.com", database.UserrolesMember, sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			expectError: false,
			expectedId:  userId,
		},
		{
			name: "Successfully create user without team (Admin)",
			userData: models.User{
				Username:  "admin",
				Password:  "hashedpassword",
				Email:     "admin@example.com",
				Role:      models.UserrolesAdmin,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "admin", "hashedpassword", "admin@example.com", "Admin", nil)
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users")).
					WithArgs(now, now, "admin", "hashedpassword", "admin@example.com", database.UserrolesAdmin, sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			expectError: false,
			expectedId:  userId,
		},
		{
			name: "Database error on create",
			userData: models.User{
				Username:  "testuser",
				Password:  "hashedpassword",
				Email:     "test@example.com",
				Role:      models.UserrolesMember,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
		{
			name: "Duplicate email error",
			userData: models.User{
				Username:  "testuser",
				Password:  "hashedpassword",
				Email:     "existing@example.com",
				Role:      models.UserrolesMember,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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
			repo := repository.NewUserRepository(queries, db)

			result, err := repo.CreateUser(context.Background(), tt.userData)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedId, result.ID)
				assert.Equal(t, tt.userData.Email, result.Email)
				assert.Equal(t, tt.userData.Username, result.Username)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name          string
		email         string
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		expectedEmail string
	}{
		{
			name:  "Successfully get user by email",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "testuser", "hashedpassword", "test@example.com", "Member", teamId)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, username, password, email, role, team_id FROM users")).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedEmail: "test@example.com",
		},
		{
			name:  "User not found",
			email: "notfound@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, username, password, email, role, team_id FROM users")).
					WithArgs("notfound@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
		},
		{
			name:  "Database connection error",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, username, password, email, role, team_id FROM users")).
					WithArgs("test@example.com").
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
			repo := repository.NewUserRepository(queries, db)

			result, err := repo.GetUserByEmail(context.Background(), tt.email)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedEmail, result.Email)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetUserById(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name        string
		userId      uuid.UUID
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectedId  uuid.UUID
	}{
		{
			name:   "Successfully get user by ID",
			userId: userId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "testuser", "hashedpassword", "test@example.com", "Member", teamId)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, username, password, email, role, team_id FROM users")).
					WithArgs(userId).
					WillReturnRows(rows)
			},
			expectError: false,
			expectedId:  userId,
		},
		{
			name:   "User not found by ID",
			userId: userId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, username, password, email, role, team_id FROM users")).
					WithArgs(userId).
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
			repo := repository.NewUserRepository(queries, db)

			result, err := repo.GetUserById(context.Background(), tt.userId)

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

func TestUserRepository_GetUsers(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name          string
		filters       interfaces.GetUserFilters
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		expectedCount int
	}{
		{
			name: "Get users with team filter",
			filters: interfaces.GetUserFilters{
				TeamId:      teamId,
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "testuser", "hashedpassword", "test@example.com", "Member", teamId)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get users with email filter",
			filters: interfaces.GetUserFilters{
				Email:       "test",
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "testuser", "hashedpassword", "test@example.com", "Member", teamId)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get users with role filter",
			filters: interfaces.GetUserFilters{
				Role:        "Admin",
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "admin", "hashedpassword", "admin@example.com", "Admin", nil)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get users - empty result",
			filters: interfaces.GetUserFilters{
				TeamId:      uuid.New(),
				Limit:       10,
				IsFirstPage: true,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 0,
		},
		{
			name: "Get users with pagination - next page",
			filters: interfaces.GetUserFilters{
				TeamId:          teamId,
				Limit:           10,
				IsFirstPage:     false,
				PointsNext:      true,
				CursorCreatedAt: now.Add(-time.Hour),
				CursorId:        uuid.New(),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "testuser", "hashedpassword", "test@example.com", "Member", teamId)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Get users with pagination - previous page",
			filters: interfaces.GetUserFilters{
				TeamId:          teamId,
				Limit:           10,
				IsFirstPage:     false,
				PointsNext:      false,
				CursorCreatedAt: now.Add(time.Hour),
				CursorId:        uuid.New(),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "testuser", "hashedpassword", "test@example.com", "Member", teamId)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "Database error",
			filters: interfaces.GetUserFilters{
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
			repo := repository.NewUserRepository(queries, db)

			result, err := repo.GetUsers(context.Background(), tt.filters)

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

func TestUserRepository_UpdateUser(t *testing.T) {
	userId := uuid.New()
	teamId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name             string
		updateData       interfaces.UpdateUserData
		mockSetup        func(sqlmock.Sqlmock)
		expectError      bool
		expectedUsername string
	}{
		{
			name: "Successfully update user",
			updateData: interfaces.UpdateUserData{
				ID:        userId,
				Username:  "updateduser",
				Email:     "updated@example.com",
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "password", "email", "role", "team_id"}).
					AddRow(userId, now, now, "updateduser", "hashedpassword", "updated@example.com", "Member", teamId)
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE users")).
					WithArgs("updateduser", "updated@example.com", now, userId).
					WillReturnRows(rows)
			},
			expectError:      false,
			expectedUsername: "updateduser",
		},
		{
			name: "Update user - not found",
			updateData: interfaces.UpdateUserData{
				ID:        userId,
				Username:  "updateduser",
				Email:     "updated@example.com",
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE users")).
					WithArgs("updateduser", "updated@example.com", now, userId).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
		},
		{
			name: "Update user - database error",
			updateData: interfaces.UpdateUserData{
				ID:        userId,
				Username:  "updateduser",
				Email:     "updated@example.com",
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE users")).
					WithArgs("updateduser", "updated@example.com", now, userId).
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
			repo := repository.NewUserRepository(queries, db)

			result, err := repo.UpdateUser(context.Background(), tt.updateData)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUsername, result.Username)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	userId := uuid.New()

	tests := []struct {
		name        string
		userId      uuid.UUID
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:   "Successfully delete user",
			userId: userId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users")).
					WithArgs(userId).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectError: false,
		},
		{
			name:   "Delete user - database error",
			userId: userId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users")).
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
			repo := repository.NewUserRepository(queries, db)

			err = repo.DeleteUser(context.Background(), tt.userId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
