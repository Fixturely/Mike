package main

import (
	"log"
	"mike/utils"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run db/migrations/main.go <up|down>")
	}

	direction := os.Args[1]
	switch direction {
	case "up":
		log.Println("Running up migrations")
	case "down":
		log.Println("Running down migrations")
	default:
		log.Fatal("Invalid direction, must be up or down")
	}
	db := utils.GetDatabase()

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Executing migrations failed: %v", err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Rolling back migrations failed: %v", err)
		}
	}

	log.Println("Migrations completed successfully")
}
