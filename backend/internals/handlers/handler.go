package handlers

import "github.com/TobiasRV/challenge-fs-senior/internals/interfaces"

type Handler struct {
	validator      *Validator
	userRepository interfaces.IUserRepository
}

func NewHandler(
	ur interfaces.IUserRepository,
) *Handler {
	v := NewValidator()
	return &Handler{
		validator:      v,
		userRepository: ur,
	}
}
