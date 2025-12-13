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

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler:   h.JWTErrorHandler,
		SuccessHandler: h.JWTSuccessHandler,
	})
	usersRoutes.Use(jwtMiddleware)
	usersRoutes.Get("/", h.GetUsers)
	usersRoutes.Post("/", h.CreateUser)
	usersRoutes.Put("/:id", h.UpdateUser)
	usersRoutes.Delete("/:id", h.DeleteUser)

	teamsRoutes := v1.Group("/teams", jwtMiddleware)
	teamsRoutes.Post("/", h.CreateTeam)
	teamsRoutes.Get("/by-owner", h.GetTeamByOwner)

	projectRoutes := v1.Group("/projects", jwtMiddleware)
	projectRoutes.Get("/", h.GetProjects)
	projectRoutes.Post("/", h.CreateProject)
	projectRoutes.Put("/:id", h.UpdateProject)
	projectRoutes.Delete("/:id", h.DeleteProject)

	taskRoutes := v1.Group("/tasks", jwtMiddleware)
	taskRoutes.Post("/", h.CreateTask)
	taskRoutes.Get("/", h.GetTasks)
	taskRoutes.Put("/:id", h.UpdateTask)
	taskRoutes.Delete("/:id", h.DeleteTask)

}
