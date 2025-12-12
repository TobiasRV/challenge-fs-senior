package handlers

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Register(r *fiber.App) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	authRoutes := v1.Group("/auth")
	authRoutes.Post("/register-admin", h.CreateUserAdmin)
	authRoutes.Post("/login", h.LogIn)
	authRoutes.Post("/refresh-token", h.RefreshToken)

	usersRoutes := v1.Group("/users")
	usersRoutes.Get("/exists-by-email", h.UserExistsByEmail)

	_ = jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler:   h.JWTErrorHandler,
		SuccessHandler: h.JWTSuccessHandler,
	})
}
