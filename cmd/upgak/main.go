package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"upgak-golang-check-up-down/internal/checker"
)

func main() {
	var (
		urls        = flag.String("urls", "", "Comma-separated list of URLs to check (required)")
		concurrent  = flag.Bool("concurrent", false, "Use goroutines for concurrent checking")
		timeoutSecs = flag.Int("timeout", 10, "Timeout in seconds for each request")
		help        = flag.Bool("help", false, "Show help message")
	)
	
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *urls == "" {
		fmt.Println("❌ Error: URLs parameter is required")
		fmt.Println("Use -help flag for usage information")
		os.Exit(1)
	}

	urlList := strings.Split(*urls, ",")
	for i, url := range urlList {
		urlList[i] = strings.TrimSpace(url)
		// Add http:// prefix if no protocol is specified
		if !strings.HasPrefix(urlList[i], "http://") && !strings.HasPrefix(urlList[i], "https://") {
			urlList[i] = "http://" + urlList[i]
		}
	}

	timeout := time.Duration(*timeoutSecs) * time.Second

	fmt.Printf("🔍 Checking %d website(s) with %s mode (timeout: %ds)\n\n", 
		len(urlList), 
		map[bool]string{true: "concurrent", false: "serial"}[*concurrent],
		*timeoutSecs)

	start := time.Now()
	var results []checker.CheckResult

	if *concurrent {
		results = checker.CheckMultipleConcurrent(urlList, timeout)
	} else {
		results = checker.CheckMultipleSerial(urlList, timeout)
	}

	// Print results
	for _, result := range results {
		checker.PrintResult(result)
	}

	totalTime := time.Since(start)
	upCount := 0
	for _, result := range results {
		if result.IsUp {
			upCount++
		}
	}

	fmt.Printf("\n📊 Summary: %d/%d websites are UP\n", upCount, len(results))
	fmt.Printf("⏱️  Total checking time: %s\n", totalTime)
}

func showHelp() {
	fmt.Println("🟢 upgak - Website Status Checker")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  upgak -urls=\"url1,url2,url3\" [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -urls string      Comma-separated list of URLs to check (required)")
	fmt.Println("  -concurrent       Use goroutines for concurrent checking (default: false)")
	fmt.Println("  -timeout int      Timeout in seconds for each request (default: 10)")
	fmt.Println("  -help             Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  upgak -urls=\"google.com,github.com,nonexistent.com\"")
	fmt.Println("  upgak -urls=\"https://google.com,https://github.com\" -concurrent")
	fmt.Println("  upgak -urls=\"google.com\" -timeout=5")
	fmt.Println()
	fmt.Println("Status Codes:")
	fmt.Println("  UP   - HTTP status codes 200-399")
	fmt.Println("  DOWN - HTTP status codes ≥400, timeouts, or connection errors")
}
