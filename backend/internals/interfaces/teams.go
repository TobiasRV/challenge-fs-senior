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

type CreateTeamRequest struct {
	Name string `json:"name" example:"Development Team"`
}

type TeamExistsResponse struct {
	Exists bool         `json:"exists" example:"true"`
	Team   *models.Team `json:"team"`
}
