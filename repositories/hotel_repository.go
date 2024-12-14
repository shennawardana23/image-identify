package repositories

import (
	"fmt"
	"image-identify/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HotelRepository interface {
	FetchHotelsWithWebsites() ([]models.Hotel, error)
}

type hotelRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{
		db:     db,
		logger: logrus.New(),
	}
}

func (r *hotelRepository) FetchHotelsWithWebsites() ([]models.Hotel, error) {
	var hotels []models.Hotel

	query := r.db.Debug().
		Joins("LEFT JOIN tb_brands Brand ON tb_hotels.brand_id = Brand.brand_id").
		Where("LOWER(Brand.brand_name) != ?", "powered by archi").
		Order("tb_hotels.hotel_id ASC")

	result := query.Find(&hotels)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error":    result.Error.Error(),
			"function": "FetchHotelsWithWebsites",
			"query":    query.Statement.SQL.String(),
			"params":   query.Statement.Vars,
		}).Error("Database query failed")
		return nil, fmt.Errorf("failed to fetch hotels: %v", result.Error)
	}

	return hotels, nil
}
