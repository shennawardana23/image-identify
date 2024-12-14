package handlers

import (
	"fmt"
	"net/http"

	"image-identify/models"
	"image-identify/services"
	"image-identify/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ImageHandler struct {
	db *gorm.DB
}

func NewImageHandler(db *gorm.DB) *ImageHandler {
	return &ImageHandler{db: db}
}

func (h *ImageHandler) CheckWebsites(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"handler": "CheckWebsites",
		"file":    "image_handler.go",
	})

	var hotels []models.Hotel

	// Query hotels with brand information, excluding "powered by archi" and order by hotel_id
	query := h.db.Debug().
		Joins("LEFT JOIN tb_brands Brand ON tb_hotels.brand_id = Brand.brand_id").
		Where("LOWER(Brand.brand_name) != ?", "powered by archi").
		Order("tb_hotels.hotel_id ASC")

	result := query.Find(&hotels)

	if result.Error != nil {
		logger.WithFields(logrus.Fields{
			"error":    result.Error.Error(),
			"function": "CheckWebsites",
			"line":     fmt.Sprintf("%s:%d", "image_handler.go", 40),
			"query":    query.Statement.SQL.String(),
			"params":   query.Statement.Vars,
			"duration": result.RowsAffected,
		}).Error("Database query failed")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch hotels from database",
			"details": result.Error.Error(),
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"hotels_found": len(hotels),
	}).Info("Successfully fetched hotels")

	if len(hotels) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No hotels found with valid website links",
		})
		return
	}

	// Process websites and get results
	results := services.ProcessWebsites(hotels, 10)

	// Generate and send Excel response
	if err := utils.GenerateCSVResponse(c, results); err != nil {
		logger.WithError(err).Error("Failed to generate Excel response")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Excel response",
		})
		return
	}

	logger.WithFields(logrus.Fields{
		"total_processed": len(results),
	}).Info("Successfully generated Excel response")
}
