package mocks

import (
	"context"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/stretchr/testify/mock"
)

type MockRefreshTokenRepository struct {
	mock.Mock
}

func NewMockRefreshTokenRepository() *MockRefreshTokenRepository {
	return &MockRefreshTokenRepository{}
}

func (m *MockRefreshTokenRepository) CreateRefreshToken(ctx context.Context, token models.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetRefreshTokenByToken(ctx context.Context, token string) (models.RefreshTokenWithUser, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(models.RefreshTokenWithUser), args.Error(1)
}

func (m *MockRefreshTokenRepository) DeleteRefreshTokensByUserId(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
