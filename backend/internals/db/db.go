package db

import (
	"database/sql"
	"log"
	"os"

	internalsql "github.com/TobiasRV/challenge-fs-senior/internals/sql"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
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

	// Run migrations
	if err := runMigrations(conn); err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	// Run seed data only in dev environment
	if os.Getenv("ENVIRONMENT") == "dev" {
		if err := runSeed(conn); err != nil {
			log.Fatal("Failed to run seed: ", err)
		}
	}

	queries = database.New(conn)

	return queries, conn
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(internalsql.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "schemas"); err != nil {
		return err
	}

	return nil
}

func runSeed(db *sql.DB) error {
	goose.SetBaseFS(internalsql.EmbedSeed)
	goose.SetTableName("goose_db_seed_version")

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "seed"); err != nil {
		return err
	}

	return nil
}
