package handlers

import (
	"database/sql"
	"errors"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/interfaces"
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

	_, err = h.userRepository.GetUserByEmail(c.Context(), payload.Email)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(utils.ErrorString("user already exists"))
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

func (h *Handler) CreateUser(c *fiber.Ctx) error {

	userRole := c.Locals("userRole")

	if userRole != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	payload := interfaces.CreateUserPayload{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err := h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	_, err = h.userRepository.GetUserByEmail(c.Context(), payload.Email)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(utils.ErrorString("user already exists"))
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
		TeamId:    payload.TeamId,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": newUser,
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

func (h *Handler) GetUsers(c *fiber.Ctx) error {

	userRole := c.Locals("userRole")

	if userRole != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	queryParams := interfaces.GetUserParams{}

	if err := c.QueryParser(&queryParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(err))
	}

	cursor := queryParams.Cursor
	isFirstPage := true
	var cursorCreatedAt time.Time
	var cursorId uuid.UUID
	pointsNext := false
	if cursor != "" {
		decodedCursor, err := utils.DecodeCursor(cursor)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
		}
		cursorCreatedAt, err = time.Parse(time.RFC3339Nano, decodedCursor["created_at"].(string))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
		}
		cursorId, err = uuid.Parse(decodedCursor["id"].(string))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
		}

		pointsNext = decodedCursor["points_next"] == true

		isFirstPage = false
	}

	users, err := h.userRepository.GetUsers(c.Context(), interfaces.GetUserFilters{
		Email:           queryParams.Email,
		TeamId:          queryParams.TeamId,
		Limit:           queryParams.Limit,
		IsFirstPage:     isFirstPage,
		CursorCreatedAt: cursorCreatedAt,
		CursorId:        cursorId,
		PointsNext:      pointsNext,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	hasPagination := len(users) > int(queryParams.Limit)

	if hasPagination {
		if cursor == "" || pointsNext {
			users = users[:int(queryParams.Limit)]
		} else {
			users = users[len(users)-int(queryParams.Limit):]
		}
	}

	if len(users) == 0 {
		pager := utils.GeneratePager(utils.Cursor{}, utils.Cursor{})
		return c.Status(200).JSON(fiber.Map{
			"data":       users,
			"pagination": pager,
		})
	}

	var nextCursor utils.Cursor
	var prevCursor utils.Cursor

	if cursor == "" {
		if hasPagination {
			nextCursor = utils.CreateCursor(users[len(users)-1].ID, users[len(users)-1].CreatedAt, true)
		}
	} else {
		if pointsNext {

			if hasPagination {
				nextCursor = utils.CreateCursor(users[len(users)-1].ID, users[len(users)-1].CreatedAt, true)
			}

			prevCursor = utils.CreateCursor(users[0].ID, users[0].CreatedAt, false)
		} else {
			nextCursor = utils.CreateCursor(users[len(users)-1].ID, users[len(users)-1].CreatedAt, true)

			if hasPagination {
				prevCursor = utils.CreateCursor(users[0].ID, users[0].CreatedAt, false)
			}
		}
	}

	pager := utils.GeneratePager(nextCursor, prevCursor)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       users,
		"pagination": pager,
	})
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {

	userRole := c.Locals("userRole")

	if userRole != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	userId := c.Params("id")

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("invalid id"))
	}

	userUUID, err := uuid.Parse(userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	currentUser, err := h.userRepository.GetUserById(c.Context(), userUUID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorString("user not found"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	payload := interfaces.UpdateUserPayload{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err = h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	if currentUser.Email != payload.Email {
		_, err = h.userRepository.GetUserByEmail(c.Context(), payload.Email)

		if err == nil {
			return c.Status(fiber.StatusConflict).JSON(utils.ErrorString("user already exists"))
		}
	}

	newUser, err := h.userRepository.UpdateUser(c.Context(), interfaces.UpdateUserData{
		Username:  payload.Username,
		Email:     payload.Email,
		UpdatedAt: time.Now().UTC(),
		ID:        userUUID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(newUser)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {

	userRole := c.Locals("userRole")

	if userRole != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	userId := c.Params("id")

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("invalid id"))
	}

	userUUID, err := uuid.Parse(userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err = h.userRepository.DeleteUser(c.Context(), userUUID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"deleted": true,
	})
}
