package interfaces

import (
	"context"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/google/uuid"
)

type ITeamRepository interface {
	CreateTeam(context.Context, models.Team) (models.Team, error)
	GetTeamByOwner(context.Context, uuid.UUID) (exists bool, team models.Team, err error)
}
