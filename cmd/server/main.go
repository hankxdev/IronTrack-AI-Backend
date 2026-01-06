package main

import (
	"log"

	"irontrack-backend/internal/database"
	"irontrack-backend/internal/router"
)

func main() {
	database.InitDatabase()

	r := router.SetupRouter()

	log.Println("Server running on port 8080")
	r.Run(":8080")
}
