package repository

import (
	"context"
	"database/sql"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
)

type RefreshTokenRepository struct {
	queries *database.Queries
	db      *sql.DB
}

func NewRefreshTokenRepository(queries *database.Queries, db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		queries: queries,
		db:      db,
	}
}

func (rtr *RefreshTokenRepository) CreateRefreshToken(c context.Context, refreshTokenData models.RefreshToken) error {
	err := rtr.queries.CreateRefreshToken(c, database.CreateRefreshTokenParams(refreshTokenData))

	if err != nil {
		return err
	}

	return nil
}

func (rtr *RefreshTokenRepository) GetRefreshTokenByToken(c context.Context, token string) (models.RefreshTokenWithUser, error) {
	refreshToken, err := rtr.queries.GetRefreshTokenByToken(c, token)

	if err != nil {
		return models.RefreshTokenWithUser{}, err
	}

	return models.DatabaseRefreshTokenWithUserToRefreshTokenWithUser(refreshToken), nil
}
