package handlers

import (
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUserAdmin(c *fiber.Ctx) error {

	payload := struct {
		Username string           `json:"username" validate:"required"`
		Password string           `json:"password" validate:"required"`
		Email    string           `json:"email" validate:"required,email"`
		Role     models.Userroles `json:"role" validate:"required,oneof=Admin Manager Member"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err := h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	newUser, err := h.userRepository.CreateUser(c.Context(), models.User{
		Username:  payload.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Email:     payload.Email,
		Role:      payload.Role,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": newUser,
	})
}
