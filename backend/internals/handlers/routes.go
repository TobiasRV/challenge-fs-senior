package handlers

import "github.com/gofiber/fiber/v2"

func (h *Handler) Register(r *fiber.App) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	authRoutes := v1.Group("/auth")
	authRoutes.Post("/register-admin", h.CreateUserAdmin)
}
