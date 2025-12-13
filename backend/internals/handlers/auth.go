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

func (h *Handler) LogIn(c *fiber.Ctx) error {
	payload := struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err := h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	existingUser, err := h.userRepository.GetUserByEmail(c.Context(), payload.Email)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusConflict).JSON(utils.ErrorString("Invalid username or password"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorString("Error getting data from database"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(payload.Password)); err != nil {
		return c.Status(fiber.StatusConflict).JSON(utils.ErrorString("Invalid credentials"))
	}

	token, err := utils.GenerateJWTToken(existingUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	refreshToken, expiresAt, err := utils.GenerateJWTRefreshToken(existingUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	h.refreshTokenRepository.CreateRefreshToken(c.Context(), models.RefreshToken{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: expiresAt,
		Userid:    existingUser.ID,
		Token:     refreshToken,
		Revoked:   false,
	})

	utils.GenerateCookie(c, token)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  token,
		"refreshToken": refreshToken,
		"user":         existingUser,
	})
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	payload := struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err := h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	existingToken, err := h.refreshTokenRepository.GetRefreshTokenByToken(c.Context(), payload.RefreshToken)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorString("invalid refresh token"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	newAccessToken, err := utils.GenerateJWTToken(existingToken.UserData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(200).JSON(fiber.Map{
		"accessToken": newAccessToken,
	})

}
