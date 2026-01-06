package router

import (
	"irontrack-backend/internal/auth"
	"irontrack-backend/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Setup
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Routes
	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		protected := api.Group("/")
		protected.Use(auth.AuthMiddleware())
		{
			protected.GET("/me", handlers.GetMe)

			// Plans
			protected.GET("/plans", handlers.GetPlans)
			protected.POST("/plans", handlers.CreatePlan)
			protected.DELETE("/plans/:id", handlers.DeletePlan)

			// Logs
			protected.GET("/logs", handlers.GetLogs)
			protected.POST("/logs", handlers.CreateLog)

			// Exercises
			protected.GET("/exercises", handlers.GetExercises)
			protected.POST("/exercises", handlers.CreateExercise)
			protected.DELETE("/exercises/:id", handlers.DeleteExercise)

			// Profile
			protected.GET("/profile", handlers.GetProfile)
			protected.POST("/profile", handlers.SaveProfile)

			// AI
			protected.POST("/ai/generate-plan", handlers.GenerateWorkoutPlan)
			protected.POST("/ai/generate-report", handlers.GenerateProgressReport)
		}
	}
	return r
}
