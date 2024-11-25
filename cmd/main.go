package main

import (
	"fmt"
	"os"

	"html-text-extraction/internal/orchestration"
	"html-text-extraction/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input-html-file|url>")
		return
	}

	input := os.Args[1]
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

	// Open the input file (either a downloaded HTML or a local file)
	file, err := os.Open(input)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	orchestration.Orchestrate(file, outputFile)
}
