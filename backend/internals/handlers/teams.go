package handlers

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/TobiasRV/challenge-fs-senior/internals/interfaces"
	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateTeam godoc
// @Summary Create a new team
// @Description Create a new team (Admin only)
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body interfaces.CreateTeamRequest true "Team creation data"
// @Success 201 {object} models.Team
// @Failure 400 {object} utils.ErrorResponse "Validation error"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} utils.ErrorResponse "User not found"
// @Failure 409 {object} utils.ErrorResponse "User already has a team"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /teams [post]
func (h *Handler) CreateTeam(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	userRole := c.Locals("userRole")
	if userRole != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	payload := struct {
		Name string `json:"name" validate:"required"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err := h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	userUUID, err := uuid.Parse(userId.(string))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	_, err = h.userRepository.GetUserById(c.Context(), userUUID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorString("user not found"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	alreadyExists, _, err := h.teamRepository.GetTeamByOwner(c.Context(), userUUID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	if alreadyExists {
		return c.Status(fiber.StatusConflict).JSON(utils.ErrorString("user already has a team with that name"))
	}

	team, err := h.teamRepository.CreateTeam(c.Context(), models.Team{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      payload.Name,
		OwnerID:   userUUID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(team)
}

// GetTeamByOwner godoc
// @Summary Get team by owner
// @Description Get the team owned by the current user (Admin only)
// @Tags Teams
// @Accept json
// @Produce json
// @Success 200 {object} interfaces.TeamExistsResponse
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin only"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /teams/by-owner [get]
func (h *Handler) GetTeamByOwner(c *fiber.Ctx) error {

	userId := c.Locals("userId")
	userRole := c.Locals("userRole")

	if userRole != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	userUUID, err := uuid.Parse(userId.(string))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	exists, team, err := h.teamRepository.GetTeamByOwner(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	if !exists {
		return c.Status(200).JSON(fiber.Map{
			"exists": exists,
			"team":   nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"exists": exists,
		"team":   team,
	})
}
