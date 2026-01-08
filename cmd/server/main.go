package main

import (
	"log"
	"os"

	"irontrack-backend/internal/database"
	"irontrack-backend/internal/router"
)

func main() {
	database.InitDatabase()

	r := router.SetupRouter()

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
