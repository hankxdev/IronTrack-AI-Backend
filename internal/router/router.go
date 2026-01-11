package router

import (
	"irontrack-backend/internal/auth"
	"irontrack-backend/internal/handlers"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Trusted Proxies Setup
	// For Render.com and other cloud platforms, we trust their reverse proxies
	// This can be configured via TRUSTED_PROXIES env var (comma-separated CIDRs)
	// If not set, we don't trust any proxy (safest for security)
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		r.SetTrustedProxies(strings.Split(trustedProxies, ","))
	} else {
		// Don't trust any proxies by default (most secure)
		r.SetTrustedProxies(nil)
	}

	// CORS Setup
	config := cors.DefaultConfig()

	// Read allowed origins from environment variable
	// Format: ALLOWED_ORIGINS=https://example.com,https://app.example.com
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		config.AllowOrigins = strings.Split(allowedOrigins, ",")
	} else {
		// Default to allow all for development
		config.AllowAllOrigins = true
	}

	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))
	r.Use(DevelopmentLogger())

	// Health check endpoint (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "irontrack-backend",
		})
	})

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

		// Admin-only routes
		admin := api.Group("/admin")
		admin.Use(auth.AuthMiddleware())
		admin.Use(auth.AdminMiddleware())
		{
			admin.GET("/summary", handlers.AdminSummary)

			// Users
			admin.GET("/users", handlers.AdminListUsers)
			admin.POST("/users", handlers.AdminCreateUser)
			admin.PUT("/users/:id", handlers.AdminUpdateUser)
			admin.DELETE("/users/:id", handlers.AdminDeleteUser)

			// Plans
			admin.GET("/plans", handlers.AdminListPlans)
			admin.POST("/plans", handlers.AdminCreatePlan)
			admin.DELETE("/plans/:id", handlers.AdminDeletePlan)

			// Exercises
			admin.GET("/exercises", handlers.AdminListExercises)
			admin.POST("/exercises", handlers.AdminCreateExercise)
			admin.DELETE("/exercises/:id", handlers.AdminDeleteExercise)

			// AI Requests log
			admin.GET("/ai-requests", handlers.AdminListAIRequests)

			admin.POST("/exercises/bulk", handlers.BulkUploadExercises)
		}
	}
	return r
}
