package interfaces

import (
	"context"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
)

type IUserRepository interface {
	CreateUser(context.Context, models.User) (models.User, error)
	GetUserByEmail(context.Context, string) (models.User, error)
}
