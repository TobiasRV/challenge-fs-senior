package models

import (
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"ownerId"`
}

type TeamFilter struct {
	Cursor string
	Limit  uint64
	Name   string
	UserId uuid.UUID
}

func DatabaseTeamToTeam(dbTeam database.Team) Team {
	return Team{
		ID:        dbTeam.ID,
		CreatedAt: dbTeam.CreatedAt,
		UpdatedAt: dbTeam.UpdatedAt,
		Name:      dbTeam.Name,
		OwnerID:   dbTeam.OwnerID,
	}
}
