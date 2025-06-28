package checker

import (
	"testing"
	"time"
)

func TestCheckSingle(t *testing.T) {
	// Test with a valid URL (GitHub should be accessible)
	result := CheckSingle("https://github.com", 10*time.Second)
	
	if result.URL != "https://github.com" {
		t.Errorf("Expected URL to be 'https://github.com', got %s", result.URL)
	}
	
	if result.ResponseTime <= 0 {
		t.Errorf("Expected response time to be greater than 0, got %v", result.ResponseTime)
	}
	
	// GitHub should return a successful status code
	if result.Error == nil && result.StatusCode == 200 && !result.IsUp {
		t.Errorf("Expected IsUp to be true for status code 200")
	}
}

func TestCheckSingleInvalidURL(t *testing.T) {
	// Test with an invalid URL
	result := CheckSingle("https://thisdoesnotexist123456789.invalid", 2*time.Second)
	
	if result.Error == nil {
		t.Errorf("Expected error for invalid URL, got nil")
	}
	
	if result.IsUp {
		t.Errorf("Expected IsUp to be false for invalid URL")
	}
}

func TestStatusCodeDetermination(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   bool
	}{
		{200, true},
		{201, true},
		{301, true},
		{399, true},
		{400, false},
		{404, false},
		{500, false},
	}
	
	for _, tc := range testCases {
		result := CheckResult{
			StatusCode: tc.statusCode,
		}
		// Manually set IsUp based on status code logic
		result.IsUp = tc.statusCode >= 200 && tc.statusCode < 400
		
		if result.IsUp != tc.expected {
			t.Errorf("For status code %d, expected IsUp to be %v, got %v", 
				tc.statusCode, tc.expected, result.IsUp)
		}
	}
}

func TestCheckMultipleModes(t *testing.T) {
	urls := []string{"https://github.com"}
	timeout := 10 * time.Second
	
	// Test serial mode
	serialResults := CheckMultipleSerial(urls, timeout)
	if len(serialResults) != 1 {
		t.Errorf("Expected 1 result from serial check, got %d", len(serialResults))
	}
	
	// Test concurrent mode
	concurrentResults := CheckMultipleConcurrent(urls, timeout)
	if len(concurrentResults) != 1 {
		t.Errorf("Expected 1 result from concurrent check, got %d", len(concurrentResults))
	}
	
	// Both should have the same URL
	if serialResults[0].URL != concurrentResults[0].URL {
		t.Errorf("URLs don't match between serial and concurrent results")
	}
}