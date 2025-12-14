package mocks

import (
	"context"

	"github.com/TobiasRV/challenge-fs-senior/internals/interfaces"
	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{}
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, params database.CreateTasksParams) (models.Task, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTasks(ctx context.Context, filters interfaces.GetTasksFilters) ([]interfaces.GetTasksResponse, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]interfaces.GetTasksResponse), args.Error(1)
}

func (m *MockTaskRepository) GetTaskById(ctx context.Context, id uuid.UUID) (models.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(ctx context.Context, data interfaces.UpdateTaskData) (models.Task, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
