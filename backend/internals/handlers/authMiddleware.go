package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/TobiasRV/challenge-fs-senior/internals/utils"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) JWTErrorHandler(c *fiber.Ctx, err error) error {
	if err == jwtware.ErrJWTMissingOrMalformed {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorString("missing or malformed JWT"))
	}
	return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorString("Invalid or expired API Key"))
}

func (h *Handler) JWTSuccessHandler(c *fiber.Ctx) error {

	cookieToken := c.Cookies("jwt")
	var tokenString string
	if cookieToken != "" {
		tokenString = cookieToken
	} else {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorString("missing or malformed JWT"))
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorString("missing or malformed JWT"))
		}

		tokenString = tokenParts[1]
	}

	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		c.ClearCookie()
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorString("missing or malformed JWT"))
	}

	userId := token.Claims.(jwt.MapClaims)["userId"].(string)
	userRole := token.Claims.(jwt.MapClaims)["userRole"].(string)
	c.Locals("userId", userId)
	c.Locals("userRole", userRole)

	return c.Next()
}
