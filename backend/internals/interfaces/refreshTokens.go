package interfaces

import (
	"context"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
)

type IRefreshTokenRepository interface {
	CreateRefreshToken(context.Context, models.RefreshToken) error
	GetRefreshTokenByToken(context.Context, string) (models.RefreshTokenWithUser, error)
	DeleteRefreshTokensByUserId(context.Context, string) error
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	AccessToken  string      `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string      `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User         models.User `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type LogoutResponse struct {
	Message string `json:"message" example:"Successfully logged out"`
}
