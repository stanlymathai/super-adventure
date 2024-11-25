package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"html-text-extraction/internal/orchestration"
	"html-text-extraction/internal/utils"
)

func main() {
	// Ask for URL or file path from the user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter URL or file path: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	input = strings.TrimSpace(input)

	outputFile := "result.json"
	tempFile := "downloaded.html"

	if utils.IsValidURL(input) {
		// Download the HTML content to a local file
		err := utils.DownloadHTML(input, tempFile)
		if err != nil {
			fmt.Printf("Error downloading HTML: %v\n", err)
			return
		}
		input = tempFile
	}

	// Pre-process the HTML with GoQuery
	cleanedHTML, err := utils.CleanAndBeautifyHTML(input)
	if err != nil {
		fmt.Printf("Error cleaning HTML: %v\n", err)
		return
	}

	// Open the cleaned HTML file
	file, err := os.Open(cleanedHTML)
	if err != nil {
		fmt.Printf("Error opening cleaned HTML file: %v\n", err)
		return
	}
	defer file.Close()

	orchestration.Orchestrate(file, outputFile)
}
