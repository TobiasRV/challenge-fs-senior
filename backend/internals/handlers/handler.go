package handlers

import "github.com/TobiasRV/challenge-fs-senior/internals/interfaces"

type Handler struct {
	validator              *Validator
	userRepository         interfaces.IUserRepository
	refreshTokenRepository interfaces.IRefreshTokenRepository
	teamRepository         interfaces.ITeamRepository
}

func NewHandler(
	ur interfaces.IUserRepository,
	rtr interfaces.IRefreshTokenRepository,
	tr interfaces.ITeamRepository,
) *Handler {
	v := NewValidator()
	return &Handler{
		validator:              v,
		userRepository:         ur,
		refreshTokenRepository: rtr,
		teamRepository:         tr,
	}
}
