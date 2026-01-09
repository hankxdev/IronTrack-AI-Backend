package database

import (
	"log"
	"irontrack-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection.
// If dsn is empty, it uses a default local postgres connection string.
func ConnectDatabase(dsn string) {
	if dsn == "" {
		dsn = "host=localhost user=hankmendix dbname=irontrack port=5432 sslmode=disable TimeZone=UTC"
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
