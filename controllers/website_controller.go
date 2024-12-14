package controllers

import (
	"fmt"
	"image-identify/repositories"
	"image-identify/services"
	"image-identify/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WebsiteController struct {
	hotelRepo repositories.HotelRepository
	logger    *logrus.Logger
}

func NewWebsiteController(hotelRepo repositories.HotelRepository) *WebsiteController {
	return &WebsiteController{
		hotelRepo: hotelRepo,
		logger:    logrus.New(),
	}
}

func (c *WebsiteController) CheckWebsites(ctx *gin.Context) {
	logger := c.logger.WithFields(logrus.Fields{
		"controller": "WebsiteController",
		"action":     "CheckWebsites",
	})

	// Fetch hotels from repository
	hotels, err := c.hotelRepo.FetchHotelsWithWebsites()
	if err != nil {
		logger.WithError(err).Error("Failed to fetch hotels")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch hotels",
		})
		return
	}

	if len(hotels) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No hotels found with valid website links",
		})
		return
	}

	// Process websites and get results
	results := services.ProcessWebsites(hotels, 10)

	// Set headers for file download - ONLY HERE
	filename := fmt.Sprintf("website_check_results_%s.csv", time.Now().Format("2006-01-02"))
	ctx.Writer.Header().Set("Content-Type", "text/csv")
	ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

	// Generate CSV response
	if err := utils.GenerateCSVResponse(ctx, results); err != nil {
		logger.WithError(err).Error("Failed to generate CSV response")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate CSV response",
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"total_processed": len(results),
	}).Info("Successfully processed websites and generated CSV")
}
