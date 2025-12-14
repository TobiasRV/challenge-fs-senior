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

func TestTeamsRepository_CreateTeam(t *testing.T) {
	teamId := uuid.New()
	ownerId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name         string
		teamData     models.Team
		mockSetup    func(sqlmock.Sqlmock)
		expectError  bool
		expectedId   uuid.UUID
		expectedName string
	}{
		{
			name: "Successfully create team",
			teamData: models.Team{
				Name:      "Test Team",
				OwnerID:   ownerId,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "owner_id"}).
					AddRow(teamId, now, now, "Test Team", ownerId)
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO teams")).
					WithArgs(now, now, "Test Team", ownerId).
					WillReturnRows(rows)
			},
			expectError:  false,
			expectedId:   teamId,
			expectedName: "Test Team",
		},
		{
			name: "Database error on create",
			teamData: models.Team{
				Name:      "Test Team",
				OwnerID:   ownerId,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO teams")).
					WithArgs(now, now, "Test Team", ownerId).
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
		{
			name: "Duplicate owner error",
			teamData: models.Team{
				Name:      "Another Team",
				OwnerID:   ownerId,
				CreatedAt: now,
				UpdatedAt: now,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO teams")).
					WithArgs(now, now, "Another Team", ownerId).
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
			repo := repository.NewTeamsRepository(queries, db)

			result, err := repo.CreateTeam(context.Background(), tt.teamData)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedId, result.ID)
				assert.Equal(t, tt.expectedName, result.Name)
				assert.Equal(t, tt.teamData.OwnerID, result.OwnerID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTeamsRepository_GetTeamByOwner(t *testing.T) {
	teamId := uuid.New()
	ownerId := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name           string
		ownerId        uuid.UUID
		mockSetup      func(sqlmock.Sqlmock)
		expectError    bool
		expectedExists bool
		expectedId     uuid.UUID
	}{
		{
			name:    "Successfully get team by owner - team exists",
			ownerId: ownerId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "owner_id"}).
					AddRow(teamId, now, now, "Test Team", ownerId)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, name, owner_id FROM teams")).
					WithArgs(ownerId).
					WillReturnRows(rows)
			},
			expectError:    false,
			expectedExists: true,
			expectedId:     teamId,
		},
		{
			name:    "Team not found - returns exists: false",
			ownerId: ownerId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, name, owner_id FROM teams")).
					WithArgs(ownerId).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:    false,
			expectedExists: false,
		},
		{
			name:    "Database connection error",
			ownerId: ownerId,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, created_at, updated_at, name, owner_id FROM teams")).
					WithArgs(ownerId).
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
			repo := repository.NewTeamsRepository(queries, db)

			exists, team, err := repo.GetTeamByOwner(context.Background(), tt.ownerId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExists, exists)
				if tt.expectedExists {
					assert.Equal(t, tt.expectedId, team.ID)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
