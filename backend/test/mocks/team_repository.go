package mocks

import (
	"context"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTeamRepository struct {
	mock.Mock
}

func NewMockTeamRepository() *MockTeamRepository {
	return &MockTeamRepository{}
}

func (m *MockTeamRepository) CreateTeam(ctx context.Context, team models.Team) (models.Team, error) {
	args := m.Called(ctx, team)
	return args.Get(0).(models.Team), args.Error(1)
}

func (m *MockTeamRepository) GetTeamByOwner(ctx context.Context, ownerId uuid.UUID) (bool, models.Team, error) {
	args := m.Called(ctx, ownerId)
	return args.Bool(0), args.Get(1).(models.Team), args.Error(2)
}
