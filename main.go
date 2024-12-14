package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db := config.InitDB()

	// Initialize Gin router
	r := gin.Default()

	// Define routes
	r.GET("/api/image", func(c *gin.Context) {
		var websites []models.Website
		if err := db.Find(&websites).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		urls := make([]string, len(websites))
		for i, website := range websites {
			urls[i] = website.ImageURL
		}

		workerCount, _ := strconv.Atoi(os.Getenv("WORKER_POOL_SIZE"))
		results := services.CheckURLs(urls, workerCount)

		// Generate CSV response
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment;filename=results.csv")

		c.Writer.Write([]byte("URL,Status,Message\n"))
		for _, result := range results {
			line := fmt.Sprintf("%s,%t,%s\n", result.URL, result.Status, result.Message)
			c.Writer.Write([]byte(line))
		}
	})

	// Start server
	port := os.Getenv("SERVER_PORT")
	r.Run(":" + port)
}
