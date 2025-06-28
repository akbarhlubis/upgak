package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"upgak-golang-check-up-down/config"
	"upgak-golang-check-up-down/internal/checker"
	"upgak-golang-check-up-down/internal/notification"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Handle special commands first
	if cfg.ShowList {
		// TODO: Implement persistent URL list storage
		fmt.Println("Feature not implemented yet: persistent URL list")
		return
	}

	if cfg.AddURL != "" {
		// TODO: Implement adding URL to persistent monitoring
		fmt.Printf("Feature not implemented yet: adding URL %s to monitoring\n", cfg.AddURL)
		return
	}

	// Load URLs from batch file if specified
	if err := cfg.LoadBatchFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading batch file: %v\n", err)
		os.Exit(1)
	}

	// Validate that we have URLs to check
	if len(cfg.URLs) == 0 {
		fmt.Fprintf(os.Stderr, "No URLs provided to check\n\n")
		config.PrintUsage()
		os.Exit(1)
	}

	// Create notifier
	notifier := notification.NewNotifier(cfg.Silent)

	// Run checks
	if cfg.OneTime {
		runOneTimeCheck(cfg.URLs, notifier)
	} else {
		runPeriodicCheck(cfg, notifier)
	}
}

func runOneTimeCheck(urls []string, notifier *notification.Notifier) {
	var wg sync.WaitGroup
	
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := checker.CheckWebsite(u)
			fmt.Println(result.String())
			
			if !result.IsUp {
				message := "Connection failed"
				if result.Error == nil {
					message = fmt.Sprintf("HTTP %d", result.StatusCode)
				} else {
					message = result.Error.Error()
				}
				notifier.SendDownNotification(u, message)
			}
		}(url)
	}
	
	wg.Wait()
}

func runPeriodicCheck(cfg *config.Config, notifier *notification.Notifier) {
	// Track previous states for notification purposes
	previousStates := make(map[string]bool)
	
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()
	
	checkCount := 0
	
	// Perform initial check
	performCheck(cfg.URLs, notifier, previousStates, cfg.NotifyOnUp)
	checkCount++
	
	fmt.Printf("Monitoring %d URL(s) every %v. Press Ctrl+C to stop.\n", 
		len(cfg.URLs), cfg.Interval)
	
	for {
		select {
		case <-ticker.C:
			performCheck(cfg.URLs, notifier, previousStates, cfg.NotifyOnUp)
			checkCount++
			
			// Check if we've reached the count limit
			if cfg.Count > 0 && checkCount >= cfg.Count {
				fmt.Printf("Completed %d checks. Exiting.\n", checkCount)
				return
			}
			
		case <-sigChan:
			fmt.Printf("\nReceived interrupt signal. Completed %d checks. Exiting.\n", checkCount)
			return
		}
	}
}

func performCheck(urls []string, notifier *notification.Notifier, previousStates map[string]bool, notifyOnUp bool) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := checker.CheckWebsite(u)
			fmt.Println(result.String())
			
			mu.Lock()
			previousUp, existed := previousStates[u]
			currentUp := result.IsUp
			previousStates[u] = currentUp
			mu.Unlock()
			
			// Send notifications based on state changes
			if !currentUp && (existed && previousUp || !existed) {
				// Site went down or is down on first check
				message := "Connection failed"
				if result.Error == nil {
					message = fmt.Sprintf("HTTP %d", result.StatusCode)
				} else {
					message = result.Error.Error()
				}
				notifier.SendDownNotification(u, message)
			} else if currentUp && existed && !previousUp && notifyOnUp {
				// Site came back up
				notifier.SendUpNotification(u)
			}
		}(url)
	}
	
	wg.Wait()
}
