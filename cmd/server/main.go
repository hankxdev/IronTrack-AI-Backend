package main

import (
	"log"
	"os"

	"irontrack-backend/internal/database"
	"irontrack-backend/internal/router"

	"github.com/joho/godotenv"
)

var (
	GitCommit = "none"
	BuildTime = "unknown"
)

func main() {
	// Load .env file (ignore error in production where env vars may be set directly)
	_ = godotenv.Load()

	database.InitDatabase()

	r := router.SetupRouter(GitCommit, BuildTime)

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
