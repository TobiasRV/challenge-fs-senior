package handlers

import (
	"database/sql"
	"errors"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUserAdmin(c *fiber.Ctx) error {

	payload := struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		// Role     models.Userroles `json:"role" validate:"required,oneof=Admin Manager Member"`
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
		Role:      models.UserrolesAdmin,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	token, err := utils.GenerateJWTToken(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	refreshToken, expiresAt, err := utils.GenerateJWTRefreshToken(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err = h.refreshTokenRepository.CreateRefreshToken(c.Context(), models.RefreshToken{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: expiresAt,
		Userid:    newUser.ID,
		Token:     refreshToken,
		Revoked:   false,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	utils.GenerateCookie(c, token)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"accessToken":  token,
		"refreshToken": refreshToken,
		"user":         newUser,
	})
}

func (h *Handler) UserExistsByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("email is required"))
	}

	_, err := h.userRepository.GetUserByEmail(c.Context(), email)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"exists": false,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exists": true,
	})
}
