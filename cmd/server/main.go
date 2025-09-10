package main

import (
	"log"
	"os"

	"logistics-api/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

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
