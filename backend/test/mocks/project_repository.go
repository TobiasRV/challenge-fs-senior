package mocks

import (
	"context"

	"github.com/TobiasRV/challenge-fs-senior/internals/interfaces"
	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProjectRepository struct {
	mock.Mock
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{}
}

func (m *MockProjectRepository) CreateProject(ctx context.Context, params database.CreateProjectParams) (models.Project, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *MockProjectRepository) GetProjects(ctx context.Context, filters interfaces.GetProjectsFilters) ([]interfaces.GetProjectsResponse, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]interfaces.GetProjectsResponse), args.Error(1)
}

func (m *MockProjectRepository) UpdateProject(ctx context.Context, data interfaces.UpdateProjectData) (models.Project, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *MockProjectRepository) GetProjectById(ctx context.Context, id uuid.UUID) (models.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *MockProjectRepository) GetProjectByManager(ctx context.Context, managerId uuid.UUID) (models.Project, error) {
	args := m.Called(ctx, managerId)
	return args.Get(0).(models.Project), args.Error(1)
}

func (m *MockProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
