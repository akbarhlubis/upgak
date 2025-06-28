package main

import (
	"bufio"
	"fmt"
	"os"
	"upgak/internal/checker"
)

func main() {
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("How many URLs do you want to check? ")
	var n int
	fmt.Scanln(&n)

	for i := 0; i < n; i++ {
		fmt.Printf("Enter URL #%d: ", i+1)
		scanner.Scan()
		url := scanner.Text()
		urls = append(urls, url)
	}

	fmt.Println("Done adding URLs. Checking status...")

	for _, u := range urls {
		ok := checker.IsUp(u)
		status := "❌ Down"
		if ok {
			status = "✅ Up"
		}
		fmt.Printf("%s is %s\n", u, status)
	}

	fmt.Println("Done checking all URLs")
}
