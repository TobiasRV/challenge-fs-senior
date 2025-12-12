package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	_ "github.com/lib/pq"
)

func New() (queries *database.Queries, dbConn *sql.DB) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not valid")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can´t connect to database ", err)
	}

	queries = database.New(conn)

	return queries, conn
}
