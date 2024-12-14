package utils

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"image-identify/services"
)

func GenerateCSVResponse(c *gin.Context, results []services.URLCheckResult) error {
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write headers
	if err := writer.Write([]string{"Hotel ID", "Hotel Name", "Website URL", "Status"}); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Count empty website links
	emptyCount := 0
	for _, result := range results {
		if result.URL == "" {
			emptyCount++
		}
	}

	// Write results
	for _, result := range results {
		var status string
		if result.Status {
			status = `Success`
		} else {
			status = `Failed`
		}

		if err := writer.Write([]string{
			fmt.Sprintf("%d", result.HotelID),
			result.HotelName,
			result.URL,
			status,
		}); err != nil {
			return fmt.Errorf("error writing CSV row: %v", err)
		}
	}

	// Write empty line before summary
	if err := writer.Write([]string{"", "", "", ""}); err != nil {
		return fmt.Errorf("error writing empty line: %v", err)
	}

	// Write summary
	if err := writer.Write([]string{
		"Summary",
		fmt.Sprintf("Total Hotels: %d", len(results)),
		fmt.Sprintf("Empty Website Links: %d", emptyCount),
		fmt.Sprintf("Check Date: %s", time.Now().Format("2006-01-02 15:04:05")),
	}); err != nil {
		return fmt.Errorf("error writing summary: %v", err)
	}

	return nil
}
