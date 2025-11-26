package main

import (
	"log"
	"os"

	"hr-pepper-server/database"
	"hr-pepper-server/handlers"
	"hr-pepper-server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Get database configuration from environment
	dbConfig := database.GetConfigFromEnv()

	// Initialize MySQL database
	if err := database.Initialize(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create Gin router
	r := gin.Default()

	// Configure CORS for Pepper robot and dashboard access
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Apply middleware
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "hr-pepper-server", "database": "mysql"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Questions endpoints
		questions := api.Group("/questions")
		{
			questions.GET("", handlers.GetQuestions)
			questions.GET("/:id", handlers.GetQuestion)
			questions.POST("", handlers.CreateQuestion)
			questions.PUT("/:id", handlers.UpdateQuestion)
			questions.DELETE("/:id", handlers.DeleteQuestion)
		}

		// Interviews endpoints
		interviews := api.Group("/interviews")
		{
			interviews.GET("", handlers.GetInterviews)
			interviews.GET("/:id", handlers.GetInterview)
			interviews.POST("", handlers.CreateInterview)
			interviews.PUT("/:id", handlers.UpdateInterview)
			interviews.DELETE("/:id", handlers.DeleteInterview)
			interviews.GET("/:id/report", handlers.GetInterviewReport)
		}

		// Responses endpoint
		api.POST("/responses", handlers.SubmitResponse)

		// Analytics endpoint
		api.GET("/analytics", handlers.GetAnalytics)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("HR Pepper Interview Server starting on port %s", port)
	log.Printf("Database: MySQL @ %s:%s/%s", dbConfig.Host, dbConfig.Port, dbConfig.Database)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
