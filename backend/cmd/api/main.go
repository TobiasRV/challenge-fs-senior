package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TobiasRV/challenge-fs-senior/internals/db"
	"github.com/TobiasRV/challenge-fs-senior/internals/handlers"
	"github.com/TobiasRV/challenge-fs-senior/internals/repository"
	"github.com/TobiasRV/challenge-fs-senior/internals/router"
	"github.com/joho/godotenv"
)

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
