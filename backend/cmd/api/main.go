package main

import (
	"fmt"
	"log"
	"os"

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

	err := r.Listen(fmt.Sprintf(":%v", portString))

	if err != nil {
		log.Fatal(err)
	}
}
