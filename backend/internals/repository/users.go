package repository

import (
	"context"
	"database/sql"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
)

type UserRepository struct {
	queries *database.Queries
	db      *sql.DB
}

func NewUserRepository(queries *database.Queries, db *sql.DB) *UserRepository {
	return &UserRepository{
		queries: queries,
		db:      db,
	}
}

func (ur *UserRepository) CreateUser(c context.Context, userData models.User) (models.User, error) {

	newUser, err := ur.queries.CreateUser(c, database.CreateUserParams{
		Username:  userData.Username,
		Password:  userData.Password,
		Email:     userData.Email,
		Role:      database.Userroles(userData.Role),
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	})

	if err != nil {
		return models.User{}, err
	}

	return models.DatabaseUserToUser(newUser), nil

}
