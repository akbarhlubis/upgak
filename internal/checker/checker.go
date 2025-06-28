package checker

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// CheckResult represents the result of checking a website
type CheckResult struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	IsUp         bool
	Error        error
}

// CheckSingle checks a single URL and returns the result
func CheckSingle(url string, timeout time.Duration) CheckResult {
	start := time.Now()
	result := CheckResult{
		URL: url,
	}

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	result.ResponseTime = time.Since(start)

	if err != nil {
		result.Error = err
		result.IsUp = false
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	// Website is up if status code is 200-399
	result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400

	return result
}

// CheckMultipleSerial checks multiple URLs sequentially
func CheckMultipleSerial(urls []string, timeout time.Duration) []CheckResult {
	results := make([]CheckResult, len(urls))
	
	for i, url := range urls {
		results[i] = CheckSingle(url, timeout)
	}
	
	return results
}

// CheckMultipleConcurrent checks multiple URLs concurrently using goroutines
func CheckMultipleConcurrent(urls []string, timeout time.Duration) []CheckResult {
	results := make([]CheckResult, len(urls))
	var wg sync.WaitGroup
	
	for i, url := range urls {
		wg.Add(1)
		go func(index int, u string) {
			defer wg.Done()
			results[index] = CheckSingle(u, timeout)
		}(i, url)
	}
	
	wg.Wait()
	return results
}

// PrintResult prints a formatted result for a single check
func PrintResult(result CheckResult) {
	status := "DOWN"
	statusIcon := "❌"
	if result.IsUp {
		status = "UP"
		statusIcon = "✅"
	}
	
	if result.Error != nil {
		fmt.Printf("%s %s - %s - ERROR: %v - Response time: %d ms\n", 
			statusIcon, result.URL, status, result.Error, result.ResponseTime.Milliseconds())
	} else {
		fmt.Printf("%s %s - %s - HTTP %d - Response time: %d ms\n", 
			statusIcon, result.URL, status, result.StatusCode, result.ResponseTime.Milliseconds())
	}
}
