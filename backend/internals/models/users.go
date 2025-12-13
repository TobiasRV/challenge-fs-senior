package models

import (
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
)

type Userroles string

const (
	UserrolesAdmin   Userroles = "Admin"
	UserrolesManager Userroles = "Manager"
	UserrolesMember  Userroles = "Member"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Role      Userroles `json:"role"`
	TeamId    uuid.UUID `json:"teamId"`
}

func DatabaseUserToUser(dbUser database.User) User {
	uuid := uuid.UUID{}
	if dbUser.TeamID.Valid {
		uuid = dbUser.TeamID.UUID
	}

	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Username:  dbUser.Username,
		Password:  dbUser.Password,
		Email:     dbUser.Email,
		Role:      Userroles(dbUser.Role),
		TeamId:    uuid,
	}
}

func DatabaseUsersToUsers(dbUsers []database.User) []User {
	res := []User{}
	for _, u := range dbUsers {
		res = append(res, DatabaseUserToUser(u))
	}

	return res
}
