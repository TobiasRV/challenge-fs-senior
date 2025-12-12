package handlers

import "github.com/TobiasRV/challenge-fs-senior/internals/interfaces"

type Handler struct {
	validator              *Validator
	userRepository         interfaces.IUserRepository
	refreshTokenRepository interfaces.IRefreshTokenRepository
}

func NewHandler(
	ur interfaces.IUserRepository,
	rtr interfaces.IRefreshTokenRepository,
) *Handler {
	v := NewValidator()
	return &Handler{
		validator:              v,
		userRepository:         ur,
		refreshTokenRepository: rtr,
	}
}
