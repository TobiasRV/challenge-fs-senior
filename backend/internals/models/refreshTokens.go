package models

import (
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID
	Userid    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

type RefreshTokenWithUser struct {
	ID        uuid.UUID
	Userid    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
	UserData  User
}

func DatabaseRefreshTokenToRefreshToken(dbRefreshToken database.RefreshToken) RefreshToken {
	return RefreshToken{
		ID:        dbRefreshToken.ID,
		Userid:    dbRefreshToken.Userid,
		Token:     dbRefreshToken.Token,
		ExpiresAt: dbRefreshToken.ExpiresAt,
		CreatedAt: dbRefreshToken.CreatedAt,
		Revoked:   dbRefreshToken.Revoked,
	}
}

func DatabaseRefreshTokenWithUserToRefreshTokenWithUser(dbRefreshToken database.GetRefreshTokenByTokenRow) RefreshTokenWithUser {
	return RefreshTokenWithUser{
		ID:        dbRefreshToken.ID,
		Userid:    dbRefreshToken.Userid,
		Token:     dbRefreshToken.Token,
		ExpiresAt: dbRefreshToken.ExpiresAt,
		CreatedAt: dbRefreshToken.CreatedAt,
		Revoked:   dbRefreshToken.Revoked,
		UserData: User{
			ID:        dbRefreshToken.Userdataid,
			CreatedAt: dbRefreshToken.Userdatacreatedat,
			UpdatedAt: dbRefreshToken.Userdataupdatedat,
			Username:  dbRefreshToken.Userdatausername,
			Password:  dbRefreshToken.Userdatapassword,
			Email:     dbRefreshToken.Userdataemail,
			Role:      Userroles(dbRefreshToken.Userdatarole),
		},
	}
}
