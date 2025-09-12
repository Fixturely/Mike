package main

import (
	"fmt"
	"log"
	"mike/config"
	"mike/pkg/application"
	"mike/pkg/server"
	"net/http"
)

func main() {
	cfg := config.GetConfig()

	app, err := application.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v\n", err)
	}

	srv, err := server.New(app)
	if err != nil {
		log.Fatalf("Failed to create server: %v\n", err)
	}

	fmt.Println("Server is running on port 8080")

	// Start the server.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server encountered an error: %v\n", err)
	}
}
