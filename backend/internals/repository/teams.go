package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
)

type TeamsRepository struct {
	queries *database.Queries
	db      *sql.DB
}

func NewTeamsRepository(queries *database.Queries, db *sql.DB) *TeamsRepository {
	return &TeamsRepository{
		queries: queries,
		db:      db,
	}
}

func (tr *TeamsRepository) CreateTeam(c context.Context, teamData models.Team) (models.Team, error) {

	newTeam, err := tr.queries.CreateTeam(c, database.CreateTeamParams{
		CreatedAt: teamData.CreatedAt,
		UpdatedAt: teamData.UpdatedAt,
		Name:      teamData.Name,
		OwnerID:   teamData.OwnerID,
	})

	if err != nil {
		return models.Team{}, err
	}

	return models.DatabaseTeamToTeam(newTeam), nil
}
func (tr *TeamsRepository) GetTeamByOwner(c context.Context, ownerId uuid.UUID) (exists bool, team models.Team, err error) {

	t, err := tr.queries.GetTeamByOwner(c, ownerId)

	if errors.Is(err, sql.ErrNoRows) {
		return false, models.Team{}, nil
	}

	if err != nil {
		return false, models.Team{}, err
	}

	return true, models.DatabaseTeamToTeam(t), nil

}
