package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;type:text" json:"id"`
	Email     string         `gorm:"uniqueIndex;type:text" json:"email"`
	Password  string         `gorm:"type:text" json:"-"` // Stored as hash
	Name      string         `gorm:"type:text" json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Plans     []WorkoutPlan        `gorm:"foreignKey:UserID" json:"plans,omitempty"`
	Logs      []WorkoutLog         `gorm:"foreignKey:UserID" json:"logs,omitempty"`
	Profile   UserProfile          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile,omitempty"`
	Exercises []ExerciseDefinition `gorm:"foreignKey:UserID" json:"exercises,omitempty"`
}

type UserProfile struct {
	UserID          string `gorm:"primaryKey;type:text" json:"userId"`
	Gender          string `json:"gender"`
	Age             string `json:"age"`
	Height          string `json:"height"`
	Weight          string `json:"weight"`
	MainGoal        string `json:"mainGoal"`
	WorkoutDuration string `json:"workoutDuration"`
	ExperienceLevel string `json:"experienceLevel"`
	WeightUnit      string `json:"weightUnit"` // 'kg' or 'lbs'
}

type ExerciseDefinition struct {
	ID     string  `gorm:"primaryKey;type:text" json:"id"`
	UserID *string `gorm:"index;type:text" json:"userId"` // Can be null for system exercises? For now, let's copy system exercises to user or make them global.
	// For simplicity, we'll store all user custom exercises here. Defaults can be handled by a seeder or frontend logic (or a "is_system" flag).
	Name         string `gorm:"type:text" json:"name"`
	MuscleGroup  string `gorm:"type:text" json:"muscleGroup"`
	Instructions string `gorm:"type:text" json:"instructions,omitempty"`
}

type WorkoutPlan struct {
	ID            string    `gorm:"primaryKey;type:text" json:"id"`
	UserID        string    `gorm:"index;type:text" json:"userId"`
	Name          string    `gorm:"type:text" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	TargetGoal    string    `gorm:"type:text" json:"targetGoal"`
	IsAiGenerated bool      `json:"isAiGenerated"`
	CreatedAt     time.Time `json:"createdAt"`

	Exercises []PlanExercise `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"exercises"`
}

type PlanExercise struct {
	ID           uint   `gorm:"primaryKey" json:"-"`
	PlanID       string `gorm:"index;type:text" json:"-"`
	Name         string `gorm:"type:text" json:"name"`
	DefaultSets  int    `json:"defaultSets"`
	DefaultReps  int    `json:"defaultReps"`
	MuscleGroup  string `json:"muscleGroup,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

type WorkoutLog struct {
	ID              string    `gorm:"primaryKey;type:text" json:"id"`
	UserID          string    `gorm:"index;type:text" json:"userId"`
	Date            time.Time `gorm:"index" json:"date"`
	DurationMinutes int       `json:"durationMinutes"`
	PlanName        string    `json:"planName,omitempty"`

	Exercises []LogExercise `gorm:"foreignKey:LogID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"exercises"`
}

type LogExercise struct {
	ID           string `gorm:"primaryKey;type:text" json:"id"`
	LogID        string `gorm:"index;type:text" json:"-"`
	Name         string `gorm:"type:text" json:"name"`
	MuscleGroup  string `json:"muscleGroup,omitempty"`
	Instructions string `json:"instructions,omitempty"`

	Sets []LogSet `gorm:"foreignKey:LogExerciseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sets"`
}

type LogSet struct {
	ID            string  `gorm:"primaryKey;type:text" json:"id"`
	LogExerciseID string  `gorm:"index;type:text" json:"-"`
	Weight        float64 `json:"weight"`
	Reps          int     `json:"reps"`
	Completed     bool    `json:"completed"`
}
