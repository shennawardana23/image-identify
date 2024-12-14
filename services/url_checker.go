package services

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"image-identify/models"

	"github.com/sirupsen/logrus"
)

type URLCheckResult struct {
	HotelID   int
	HotelName string
	URL       string
	Status    bool
	Message   string
}

// CheckURL performs the actual HTTP request to verify if a URL is accessible
func CheckURL(url string) (bool, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return true, nil
}

// ProcessWebsites processes a list of hotels and checks their websites
func ProcessWebsites(hotels []models.Hotel, workerCount int) []URLCheckResult {
	var (
		wg           sync.WaitGroup
		successCount int64
		failureCount int64
		emptyCount   int64
	)

	logger := logrus.WithFields(logrus.Fields{
		"function": "ProcessWebsites",
		"workers":  workerCount,
	})

	results := make([]URLCheckResult, len(hotels))
	jobs := make(chan int, len(hotels))

	// Start workers
	for w := 0; w < workerCount; w++ {
		go func() {
			for j := range jobs {
				hotel := hotels[j]

				// Handle empty website links
				if hotel.WebsiteLink == "" {
					results[j] = URLCheckResult{
						HotelID:   hotel.HotelID,
						HotelName: hotel.HotelName,
						URL:       "",
						Status:    false,
						Message:   "No website link provided",
					}
					atomic.AddInt64(&emptyCount, 1)
					wg.Done()
					continue
				}

				_, err := CheckURL(hotel.WebsiteLink)
				if err != nil {
					results[j] = URLCheckResult{
						HotelID:   hotel.HotelID,
						HotelName: hotel.HotelName,
						URL:       hotel.WebsiteLink,
						Status:    false,
						Message:   err.Error(),
					}
					atomic.AddInt64(&failureCount, 1)
				} else {
					results[j] = URLCheckResult{
						HotelID:   hotel.HotelID,
						HotelName: hotel.HotelName,
						URL:       hotel.WebsiteLink,
						Status:    true,
						Message:   "OK",
					}
					atomic.AddInt64(&successCount, 1)
				}
				wg.Done()
			}
		}()
	}

	// Send jobs
	for i := range hotels {
		jobs <- i
		wg.Add(1)
	}
	close(jobs)

	wg.Wait()

	logger.WithFields(logrus.Fields{
		"total":       len(hotels),
		"success":     atomic.LoadInt64(&successCount),
		"failures":    atomic.LoadInt64(&failureCount),
		"empty_links": atomic.LoadInt64(&emptyCount),
	}).Info("Completed URL checking")

	return results
}
