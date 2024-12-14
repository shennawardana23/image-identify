package main

import (
	"image-identify/config"
	"image-identify/controllers"
	"image-identify/repositories"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	// Initialize DB
	db := config.InitDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
		return
	}

	// Initialize repository
	hotelRepo := repositories.NewHotelRepository(db)

	// Initialize controller
	websiteController := controllers.NewWebsiteController(hotelRepo)

	// Setup router
	router := gin.Default()

	// Routes
	api := router.Group("/api")
	{
		api.GET("/link-checker", websiteController.CheckWebsites)
	}

	// Get port from env or use default
	port := getEnvWithDefault("SERVER_PORT", "8080")
	router.Run(":" + port)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
