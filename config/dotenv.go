package config

import (
	"log"

	"github.com/joho/godotenv"
)

// loadDotenv loads environment variables from a local .env file if present.
// It is safe to call multiple times; errors are non-fatal in development.
func loadDotenv() error {
	if err := godotenv.Load(); err != nil {
		// Do not fail if .env is missing; just log in development
		log.Printf(".env not loaded: %v", err)
		return err
	}
	return nil
}
