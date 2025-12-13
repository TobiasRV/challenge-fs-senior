package handlers

import "github.com/TobiasRV/challenge-fs-senior/internals/interfaces"

type Handler struct {
	validator              *Validator
	userRepository         interfaces.IUserRepository
	refreshTokenRepository interfaces.IRefreshTokenRepository
	teamRepository         interfaces.ITeamRepository
	projectRepository      interfaces.IProjectRepository
	taskRepository         interfaces.ITaskRepository
}

func NewHandler(
	ur interfaces.IUserRepository,
	rtr interfaces.IRefreshTokenRepository,
	tr interfaces.ITeamRepository,
	pr interfaces.IProjectRepository,
	tsr interfaces.ITaskRepository,
) *Handler {
	v := NewValidator()
	return &Handler{
		validator:              v,
		userRepository:         ur,
		refreshTokenRepository: rtr,
		teamRepository:         tr,
		projectRepository:      pr,
		taskRepository:         tsr,
	}
}
