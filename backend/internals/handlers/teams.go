package handlers

import (
	"database/sql"
	"errors"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) CreateTeam(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	userRole := c.Locals("userRole")
	if userRole != "Admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorString("unauthorized"))
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
