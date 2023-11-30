//go:build integrationDB

package tests

import (
	"homework5/tests/fixtures/postgres"
	"log"

	"github.com/joho/godotenv"
)

var (
	database *postgres.TestDb
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database = postgres.NewFromEnv()
}
