package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/TobiasRV/challenge-fs-senior/docs"
	"github.com/TobiasRV/challenge-fs-senior/internals/db"
	"github.com/TobiasRV/challenge-fs-senior/internals/handlers"
	"github.com/TobiasRV/challenge-fs-senior/internals/repository"
	"github.com/TobiasRV/challenge-fs-senior/internals/router"
	"github.com/joho/godotenv"
)

// @title Challenge FS Senior API
// @version 1.0
// @description API for managing teams, projects, tasks, and users
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	godotenv.Load()

	portString := os.Getenv("SERVER_PORT")

	if portString == "" {
		log.Fatal("missing SERVER_PORT env")
	}

	r := router.New()

	queries, dbConn := db.New()

	ur := repository.NewUserRepository(queries, dbConn)
	rtr := repository.NewRefreshTokenRepository(queries, dbConn)
	tr := repository.NewTeamsRepository(queries, dbConn)
	pr := repository.NewProjectRepository(queries, dbConn)
	tsr := repository.NewTaskRepository(queries, dbConn)

	h := handlers.NewHandler(ur, rtr, tr, pr, tsr)

	h.Register(r)

	err := r.Listen(fmt.Sprintf(":%v", portString))

	if err != nil {
		log.Fatal(err)
	}
}
