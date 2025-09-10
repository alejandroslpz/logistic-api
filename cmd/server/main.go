package main

import (
	"log"

	"logistics-api/internal/app"
)

func main() {
	// Initialize container with dependency injection
	container, err := app.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application container: %v", err)
	}

	// Start the application
	if err := container.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
