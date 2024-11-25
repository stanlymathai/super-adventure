package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const MaxMemorySize = 10 * 1024 * 1024 // 10 MB

// GetHTMLReader determines if the input is a URL and returns an appropriate reader.
func GetHTMLReader(input string) (io.Reader, error) {
	if !IsValidURL(input) {
		return nil, fmt.Errorf("Invalid URL: %s", input)
	}

	fmt.Println("Fetching HTML from URL...")

	// Create a new HTTP request with headers
	req, err := http.NewRequest("GET", input, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	// Set headers to emulate a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", input)

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error fetching URL: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch URL. HTTP status: %s", resp.Status)
	}

	// Determine content length and proceed accordingly
	if resp.ContentLength > MaxMemorySize {
		fmt.Println("Large file detected, saving to temporary file...")
		tempFile, err := os.CreateTemp("", "large_html_*.tmp")
		if err != nil {
			return nil, fmt.Errorf("Error creating temporary file: %v", err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(tempFile, resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error saving large HTML to temporary file: %v", err)
		}

		_, err = tempFile.Seek(0, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("Error seeking in temporary file: %v", err)
		}
		return tempFile, nil
	}

	// Stream smaller responses directly
	return resp.Body, nil
}

// IsValidURL checks if the given input is a valid URL by attempting an HTTP GET request.
func IsValidURL(input string) bool {
	req, err := http.NewRequest("HEAD", input, nil)
	if err != nil {
		return false
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// WriteToFile writes data to a file in a formatted JSON structure.
func WriteToFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %v", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing JSON to file: %v", err)
	}

	return nil
}
