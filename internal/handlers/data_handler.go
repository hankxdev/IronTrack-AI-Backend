package handlers

import (
	"net/http"

	"irontrack-backend/internal/database"
	"irontrack-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// --- Plans ---

func GetPlans(c *gin.Context) {
	userID := c.GetString("userID")
	var plans []models.WorkoutPlan
	if err := database.DB.Preload("Exercises").Where("user_id = ?", userID).Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plans"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func CreatePlan(c *gin.Context) {
	userID := c.GetString("userID")
	var plan models.WorkoutPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan.UserID = userID
	// Ensure ID is set if not provided? Frontend usually generates UUIDs, but backend can enforce.
	// We will trust frontend provided ID or generate one if missing logic is added, but Gorm handles insertion.
	// Ideally we should overwrite ID if we want to ensure uniqueness via backend, but let's assume UUID from FE or simple checks.
	// Actually, better to just let DB handle it or validate.

	if err := database.DB.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan"})
		return
	}
	c.JSON(http.StatusCreated, plan)
}

func DeletePlan(c *gin.Context) {
	userID := c.GetString("userID")
	planID := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).Delete(&models.WorkoutPlan{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted"})
}

// --- Logs ---

func GetLogs(c *gin.Context) {
	userID := c.GetString("userID")
	var logs []models.WorkoutLog
	// Preload nested structure
	if err := database.DB.Preload("Exercises.Sets").Where("user_id = ?", userID).Order("date desc").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func CreateLog(c *gin.Context) {
	userID := c.GetString("userID")
	var log models.WorkoutLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.UserID = userID

	if err := database.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create log"})
		return
	}
	c.JSON(http.StatusCreated, log)
}

// --- Exercises ---

func GetExercises(c *gin.Context) {
	userID := c.GetString("userID")
	var exercises []models.ExerciseDefinition
	
	// Fetch global exercises OR exercises created by this user
	if err := database.DB.Where("is_global = ? OR user_id = ?", true, userID).Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}
	c.JSON(http.StatusOK, exercises)
}

func CreateExercise(c *gin.Context) {
	userID := c.GetString("userID")
	var exercise models.ExerciseDefinition
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	exercise.UserID = &userID
	exercise.IsGlobal = false // User exercises are never global

	if err := database.DB.Create(&exercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}
	c.JSON(http.StatusCreated, exercise)
}

func DeleteExercise(c *gin.Context) {
	userID := c.GetString("userID")
	exerciseID := c.Param("id")
	
	// Only allow deleting user's own exercises (not global ones)
	if err := database.DB.Where("id = ? AND user_id = ?", exerciseID, userID).Delete(&models.ExerciseDefinition{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted"})
}

// BulkUploadExercises allows admins to upload global exercises
func BulkUploadExercises(c *gin.Context) {
	var exercises []models.ExerciseDefinition
	if err := c.ShouldBindJSON(&exercises); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	if len(exercises) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No exercises provided"})
		return
	}

	// Mark all exercises as global and set UserID to nil
	for i := range exercises {
		exercises[i].IsGlobal = true
		exercises[i].UserID = nil
	}

	// Bulk insert exercises
	if err := database.DB.Create(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload exercises", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Global exercises uploaded successfully",
		"count":   len(exercises),
		"data":    exercises,
	})
}

// --- Profile ---

func GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	var profile models.UserProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		// Return empty or default
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func SaveProfile(c *gin.Context) {
	userID := c.GetString("userID")
	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile.UserID = userID

	// Upsert
	if err := database.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile"})
		return
	}
	c.JSON(http.StatusOK, profile)
}
