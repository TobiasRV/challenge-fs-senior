package mocks

import (
	"context"
	"database/sql"

	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// IQueries defines the interface for database queries used by repositories
type IQueries interface {
	// User methods
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUserByEmail(ctx context.Context, email string) (database.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (database.User, error)
	UpdateUser(ctx context.Context, arg database.UpdateUserParams) (database.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Team methods
	CreateTeam(ctx context.Context, arg database.CreateTeamParams) (database.Team, error)
	GetTeamByOwner(ctx context.Context, ownerID uuid.UUID) (database.Team, error)

	// Project methods
	CreateProject(ctx context.Context, arg database.CreateProjectParams) (database.Project, error)
	GetProjectById(ctx context.Context, id uuid.UUID) (database.Project, error)
	GetProjectByManager(ctx context.Context, managerID uuid.UUID) (database.Project, error)
	UpdateProject(ctx context.Context, arg database.UpdateProjectParams) (database.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error

	// Task methods
	CreateTasks(ctx context.Context, arg database.CreateTasksParams) (database.Task, error)
	GetTaskById(ctx context.Context, id uuid.UUID) (database.Task, error)
	UpdateTask(ctx context.Context, arg database.UpdateTaskParams) (database.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error

	// RefreshToken methods
	CreateRefreshToken(ctx context.Context, arg database.CreateRefreshTokenParams) error
	GetRefreshTokenByToken(ctx context.Context, token string) (database.GetRefreshTokenByTokenRow, error)
	DeleteRefreshTokensByUserId(ctx context.Context, userid uuid.UUID) error
	UpdaterefreshToken(ctx context.Context, arg database.UpdaterefreshTokenParams) error
}

// MockQueries mocks the SQLC database.Queries struct
type MockQueries struct {
	mock.Mock
}

func NewMockQueries() *MockQueries {
	return &MockQueries{}
}

// Verify MockQueries implements IQueries
var _ IQueries = (*MockQueries)(nil)

// User methods
func (m *MockQueries) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockQueries) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockQueries) GetUserById(ctx context.Context, id uuid.UUID) (database.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockQueries) UpdateUser(ctx context.Context, arg database.UpdateUserParams) (database.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockQueries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Team methods
func (m *MockQueries) CreateTeam(ctx context.Context, arg database.CreateTeamParams) (database.Team, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.Team), args.Error(1)
}

func (m *MockQueries) GetTeamByOwner(ctx context.Context, ownerID uuid.UUID) (database.Team, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(database.Team), args.Error(1)
}

// Project methods
func (m *MockQueries) CreateProject(ctx context.Context, arg database.CreateProjectParams) (database.Project, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.Project), args.Error(1)
}

func (m *MockQueries) GetProjectById(ctx context.Context, id uuid.UUID) (database.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.Project), args.Error(1)
}

func (m *MockQueries) GetProjectByManager(ctx context.Context, managerID uuid.UUID) (database.Project, error) {
	args := m.Called(ctx, managerID)
	return args.Get(0).(database.Project), args.Error(1)
}

func (m *MockQueries) UpdateProject(ctx context.Context, arg database.UpdateProjectParams) (database.Project, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.Project), args.Error(1)
}

func (m *MockQueries) DeleteProject(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Task methods
func (m *MockQueries) CreateTasks(ctx context.Context, arg database.CreateTasksParams) (database.Task, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.Task), args.Error(1)
}

func (m *MockQueries) GetTaskById(ctx context.Context, id uuid.UUID) (database.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.Task), args.Error(1)
}

func (m *MockQueries) UpdateTask(ctx context.Context, arg database.UpdateTaskParams) (database.Task, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.Task), args.Error(1)
}

func (m *MockQueries) DeleteTask(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// RefreshToken methods
func (m *MockQueries) CreateRefreshToken(ctx context.Context, arg database.CreateRefreshTokenParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQueries) GetRefreshTokenByToken(ctx context.Context, token string) (database.GetRefreshTokenByTokenRow, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(database.GetRefreshTokenByTokenRow), args.Error(1)
}

func (m *MockQueries) DeleteRefreshTokensByUserId(ctx context.Context, userid uuid.UUID) error {
	args := m.Called(ctx, userid)
	return args.Error(0)
}

func (m *MockQueries) UpdaterefreshToken(ctx context.Context, arg database.UpdaterefreshTokenParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// MockDB is a mock for sql.DB operations used by squirrel queries
type MockDB struct {
	mock.Mock
}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	callArgs := m.Called(query, args)
	if callArgs.Get(0) == nil {
		return nil, callArgs.Error(1)
	}
	return callArgs.Get(0).(*sql.Rows), callArgs.Error(1)
}
