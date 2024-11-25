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

	// Open the input file (either a downloaded HTML or a local file)
	file, err := os.Open(input)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	orchestration.Orchestrate(file, outputFile)
}
