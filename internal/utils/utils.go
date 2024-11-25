package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadHTML downloads the HTML content from a given URL and saves it to a local file.
func DownloadHTML(url, filename string) error {
	fmt.Printf("Downloading HTML content from URL: %s\n", url)

	// Get the response from the URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to fetch URL. HTTP status: %s", resp.Status)
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	fmt.Printf("HTML content saved to %s\n", filename)
	return nil
}
