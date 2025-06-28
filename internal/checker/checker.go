package checker

import (
	"fmt"
	"net/http"
	"time"
)

// CheckResult contains the result of a website check
type CheckResult struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	IsUp         bool
	Error        error
	Timestamp    time.Time
}

// CheckWebsite performs a single HTTP check on the given URL
func CheckWebsite(url string) *CheckResult {
	start := time.Now()
	result := &CheckResult{
		URL:       url,
		Timestamp: start,
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
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
	result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 400

	return result
}

// String returns a formatted string representation of the check result
func (cr *CheckResult) String() string {
	status := "DOWN"
	if cr.IsUp {
		status = "UP"
	}

	if cr.Error != nil {
		return fmt.Sprintf("[%s] %s - %s (%dms) - ERROR: %v",
			cr.Timestamp.Format("15:04:05"),
			cr.URL,
			status,
			cr.ResponseTime.Milliseconds(),
			cr.Error)
	}

	return fmt.Sprintf("[%s] %s - %s (HTTP %d, %dms)",
		cr.Timestamp.Format("15:04:05"),
		cr.URL,
		status,
		cr.StatusCode,
		cr.ResponseTime.Milliseconds())
}
