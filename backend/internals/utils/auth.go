package utils

import (
	"os"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CustomerClaims struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
	jwt.Claims
}

func GenerateJWTToken(user models.User) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	method := jwt.SigningMethodHS256
	claims := jwt.MapClaims{
		"userId":   user.ID,
		"userRole": user.Role,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	token, err := jwt.NewWithClaims(method, claims).SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateJWTRefreshToken(user models.User) (refreshToken string, expiresAt time.Time, err error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	method := jwt.SigningMethodHS256
	expiresAt = time.Now().UTC().Add(time.Hour * 168) // 7 days
	claims := jwt.MapClaims{
		"userId": user.ID,
		"exp":    expiresAt.Unix(),
	}

	token, err := jwt.NewWithClaims(method, claims).SignedString(jwtSecret)
	if err != nil {
		return "", time.Now(), err
	}

	return token, expiresAt, nil
}

func GenerateCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: !c.IsFromLocal(),
		Secure:   !c.IsFromLocal(),
		MaxAge:   3600 * 24 * 7, // 7 days
	})
}
