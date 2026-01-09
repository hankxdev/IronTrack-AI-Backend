package database

import (
	"fmt"
	"log"
	"os"

	"irontrack-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection.
// If dsn is empty, it reads from environment variables or uses defaults.
func ConnectDatabase(dsn string) {
	if dsn == "" {
		// Check for DATABASE_URL first (for services like Render that provide full connection string)
		dsn = os.Getenv("DATABASE_URL")
		
		// If DATABASE_URL not set, build DSN from individual env vars
		if dsn == "" {
			host := os.Getenv("DB_HOST")
			if host == "" {
				host = "localhost"
			}
			
			port := os.Getenv("DB_PORT")
			if port == "" {
				port = "5432"
			}
			
			user := os.Getenv("DB_USER")
			if user == "" {
				user = "postgres"
			}
			
			password := os.Getenv("DB_PASSWORD")
			if password == "" {
				password = ""
			}
			
			dbname := os.Getenv("DB_NAME")
			if dbname == "" {
				dbname = "irontrack"
			}
			
			sslmode := os.Getenv("DB_SSLMODE")
			if sslmode == "" {
				sslmode = "disable"
			}
			
			// Build DSN string
			if password != "" {
				dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
					host, port, user, password, dbname, sslmode)
			} else {
				dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s TimeZone=UTC",
					host, port, user, dbname, sslmode)
			}
		}
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
