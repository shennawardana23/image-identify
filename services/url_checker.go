package services

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type URLCheckResult struct {
	URL     string
	Status  bool
	Message string
}

func CheckURLs(urls []string, workerCount int) []URLCheckResult {
	var (
		wg           sync.WaitGroup
		successCount int64
		failureCount int64
	)

	results := make([]URLCheckResult, len(urls))
	jobs := make(chan int, len(urls))

	// Start workers
	for w := 0; w < workerCount; w++ {
		go worker(&wg, jobs, urls, &results, &successCount, &failureCount)
	}

	// Send jobs
	for i := range urls {
		jobs <- i
		wg.Add(1)
	}
	close(jobs)

	wg.Wait()

	log.Printf("Checked %d URLs. Success: %d, Failure: %d",
		len(urls),
		atomic.LoadInt64(&successCount),
		atomic.LoadInt64(&failureCount))

	return results
}

func worker(wg *sync.WaitGroup, jobs <-chan int, urls []string, results *[]URLCheckResult,
	successCount, failureCount *int64) {
	for j := range jobs {
		status, err := checkURL(urls[j])
		if err != nil {
			(*results)[j] = URLCheckResult{
				URL:     urls[j],
				Status:  false,
				Message: err.Error(),
			}
			atomic.AddInt64(failureCount, 1)
		} else {
			(*results)[j] = URLCheckResult{
				URL:     urls[j],
				Status:  true,
				Message: "OK",
			}
			atomic.AddInt64(successCount, 1)
		}
		wg.Done()
	}
}

func checkURL(url string) (bool, error) {
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
