package database

import (
	"log"
	"os"

	"irontrack-backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection.
// If dsn is empty, it reads from DATABASE_PATH env var or uses default.
func ConnectDatabase(dsn string) {
	if dsn == "" {
		// Check for DATABASE_PATH environment variable (for Render persistent disk)
		dsn = os.Getenv("DATABASE_PATH")
		if dsn == "" {
			dsn = "irontrack.db"
		}
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the schema
	log.Println("Migrating database schema...")
	err = DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.ExerciseDefinition{},
		&models.WorkoutPlan{},
		&models.PlanExercise{},
		&models.WorkoutLog{},
		&models.LogExercise{},
		&models.LogSet{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed.")
}

func InitDatabase() {
	ConnectDatabase("")
}
