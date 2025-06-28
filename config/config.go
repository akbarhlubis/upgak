package config

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// Config holds all configuration options for upgak
type Config struct {
	URLs         []string
	Interval     time.Duration
	Count        int
	Silent       bool
	NotifyOnUp   bool
	BatchFile    string
	AddURL       string
	ShowList     bool
	OneTime      bool
}

// ParseFlags parses command line arguments and returns a Config
func ParseFlags() (*Config, error) {
	config := &Config{}

	var intervalSec int
	var intervalSet bool
	flag.IntVar(&intervalSec, "interval", 60, "Check interval in seconds")
	flag.IntVar(&config.Count, "count", 0, "Number of checks to perform (0 = infinite)")
	flag.BoolVar(&config.Silent, "silent", false, "Disable notifications")
	flag.BoolVar(&config.Silent, "no-notify", false, "Disable notifications (alias for --silent)")
	flag.BoolVar(&config.NotifyOnUp, "notify-on-up", false, "Send notification when site comes back up")
	flag.StringVar(&config.BatchFile, "batch", "", "File containing URLs to check")
	flag.StringVar(&config.AddURL, "add", "", "Add URL to monitoring")
	flag.BoolVar(&config.ShowList, "list", false, "Show list of monitored URLs")

	// Parse flags first to determine if interval was explicitly set
	flag.Parse()

	// Check if interval was explicitly set by comparing with default
	for _, arg := range os.Args[1:] {
		if arg == "--interval" || strings.HasPrefix(arg, "--interval=") {
			intervalSet = true
			break
		}
	}

	config.Interval = time.Duration(intervalSec) * time.Second

	// Get URLs from command line arguments
	args := flag.Args()
	if len(args) > 0 {
		config.URLs = args
	}

	// Determine if this is a one-time check
	// One-time if: no explicit interval set AND (count is 0 OR count is 1)
	config.OneTime = !intervalSet && (config.Count == 0 || config.Count == 1)

	return config, nil
}

// LoadBatchFile loads URLs from a batch file (supports .txt, .json, .csv)
func (c *Config) LoadBatchFile() error {
	if c.BatchFile == "" {
		return nil
	}

	data, err := ioutil.ReadFile(c.BatchFile)
	if err != nil {
		return fmt.Errorf("failed to read batch file: %v", err)
	}

	ext := strings.ToLower(c.BatchFile)
	
	if strings.HasSuffix(ext, ".json") {
		return c.loadJSONUrls(data)
	} else if strings.HasSuffix(ext, ".csv") {
		return c.loadCSVUrls(data)
	} else {
		// Default to text file format
		return c.loadTextUrls(data)
	}
}

func (c *Config) loadTextUrls(data []byte) error {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			c.URLs = append(c.URLs, line)
		}
	}
	return nil
}

func (c *Config) loadJSONUrls(data []byte) error {
	var urls []string
	if err := json.Unmarshal(data, &urls); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}
	c.URLs = append(c.URLs, urls...)
	return nil
}

func (c *Config) loadCSVUrls(data []byte) error {
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}
	
	for i, record := range records {
		if len(record) > 0 {
			url := strings.TrimSpace(record[0])
			// Skip header row (if it contains "url" or "http" patterns)
			if i == 0 && (url == "url" || url == "URL" || !strings.Contains(url, "://")) {
				continue
			}
			if url != "" && !strings.HasPrefix(url, "#") {
				c.URLs = append(c.URLs, url)
			}
		}
	}
	return nil
}

// PrintUsage prints usage information
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <url1> [url2...]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s https://example.com                    # One-time check\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s https://example.com --interval 60     # Check every 60 seconds\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --batch urls.txt --interval 120       # Batch check every 2 minutes\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --add https://newsite.com             # Add URL to monitoring\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s --list                                # Show monitored URLs\n", os.Args[0])
}
