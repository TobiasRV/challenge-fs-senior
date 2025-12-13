package interfaces

import (
	"context"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/google/uuid"
)

type IUserRepository interface {
	CreateUser(context.Context, models.User) (models.User, error)
	GetUserByEmail(context.Context, string) (models.User, error)
	GetUserById(context.Context, uuid.UUID) (models.User, error)
	GetUsers(context.Context, GetUserFilters) ([]models.User, error)
	UpdateUser(context.Context, UpdateUserData) (models.User, error)
	DeleteUser(context.Context, uuid.UUID) error
}

type GetUserParams struct {
	Email  string           `query:"email"`
	TeamId uuid.UUID        `query:"teamId,required, uuid"`
	Role   models.Userroles `query:"role,oneof=Admin Manager Member"`
	Cursor string           `query:"cursor"`
	Limit  uint64           `query:"limit,required"`
}

type GetUserFilters struct {
	Email           string
	TeamId          uuid.UUID
	Role            models.Userroles
	Limit           uint64
	IsFirstPage     bool
	PointsNext      bool
	CursorCreatedAt time.Time
	CursorId        uuid.UUID
}

type CreateUserPayload struct {
	Username string           `json:"username" validate:"required"`
	Email    string           `json:"email" validate:"required,email"`
	Password string           `json:"password" validate:"required"`
	Role     models.Userroles `json:"role" validate:"required,oneof=Admin Manager Member"`
	TeamId   uuid.UUID        `json:"teamId" validate:"uuid"`
}

type UpdateUserPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type UpdateUserData struct {
	Username  string
	Email     string
	UpdatedAt time.Time
	ID        uuid.UUID
}
