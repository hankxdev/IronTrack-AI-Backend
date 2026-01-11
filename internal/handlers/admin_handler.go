package handlers

import (
	"net/http"
	"time"

	"irontrack-backend/internal/database"
	"irontrack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AdminSummary(c *gin.Context) {
	var userCount, planCount, exerciseCount, aiCount int64

	if err := database.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user count"})
		return
	}
	if err := database.DB.Model(&models.WorkoutPlan{}).Count(&planCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load plan count"})
		return
	}
	if err := database.DB.Model(&models.ExerciseDefinition{}).Count(&exerciseCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load exercise count"})
		return
	}
	if err := database.DB.Model(&models.AIRequestLog{}).Count(&aiCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load AI request count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":      userCount,
		"plans":      planCount,
		"exercises":  exerciseCount,
		"aiRequests": aiCount,
	})
}

func AdminListUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Order("created_at desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

type AdminCreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	IsAdmin  bool   `json:"isAdmin"`
}

func AdminCreateUser(c *gin.Context) {
	var req AdminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		IsAdmin:   req.IsAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

type AdminUpdateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	IsAdmin  *bool   `json:"isAdmin"`
}

func AdminUpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}
	user.UpdatedAt = time.Now()

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func AdminDeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := database.DB.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func AdminListPlans(c *gin.Context) {
	var plans []models.WorkoutPlan
	if err := database.DB.Preload("Exercises").Order("created_at desc").Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load plans"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func AdminCreatePlan(c *gin.Context) {
	var plan models.WorkoutPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if plan.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	if plan.ID == "" {
		plan.ID = uuid.New().String()
	}
	if plan.CreatedAt.IsZero() {
		plan.CreatedAt = time.Now()
	}

	if err := database.DB.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan"})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func AdminDeletePlan(c *gin.Context) {
	planID := c.Param("id")
	if err := database.DB.Where("id = ?", planID).Delete(&models.WorkoutPlan{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted"})
}

func AdminListExercises(c *gin.Context) {
	var exercises []models.ExerciseDefinition
	if err := database.DB.Order("name asc").Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load exercises"})
		return
	}
	c.JSON(http.StatusOK, exercises)
}

type AdminExerciseRequest struct {
	Name         string  `json:"name" binding:"required"`
	MuscleGroup  string  `json:"muscleGroup"`
	Instructions string  `json:"instructions"`
	UserID       *string `json:"userId"`
	IsGlobal     bool    `json:"isGlobal"`
}

func AdminCreateExercise(c *gin.Context) {
	var req AdminExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exercise := models.ExerciseDefinition{
		ID:           uuid.New().String(),
		Name:         req.Name,
		MuscleGroup:  req.MuscleGroup,
		Instructions: req.Instructions,
		IsGlobal:     req.IsGlobal,
	}

	if !req.IsGlobal {
		if req.UserID == nil || *req.UserID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required for non-global exercises"})
			return
		}
		var user models.User
		if err := database.DB.Select("id").Where("id = ?", *req.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		exercise.UserID = req.UserID
	} else {
		exercise.UserID = nil
	}

	if err := database.DB.Create(&exercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

func AdminDeleteExercise(c *gin.Context) {
	exerciseID := c.Param("id")
	if err := database.DB.Where("id = ?", exerciseID).Delete(&models.ExerciseDefinition{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted"})
}

func AdminListAIRequests(c *gin.Context) {
	var logs []models.AIRequestLog
	if err := database.DB.Order("created_at desc").Limit(200).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load AI requests"})
		return
	}
	c.JSON(http.StatusOK, logs)
}
