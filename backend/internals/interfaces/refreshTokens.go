package interfaces

import (
	"context"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
)

type IRefreshTokenRepository interface {
	CreateRefreshToken(context.Context, models.RefreshToken) error
	GetRefreshTokenByToken(context.Context, string) (models.RefreshTokenWithUser, error)
}
