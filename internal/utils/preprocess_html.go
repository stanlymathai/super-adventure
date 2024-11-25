package utils

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// CleanAndBeautifyHTML uses goquery to clean up the HTML and remove unnecessary elements, saving the cleaned HTML to a new file.
func CleanAndBeautifyHTML(inputFile string) (string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return "", fmt.Errorf("Error opening input HTML file: %v", err)
	}
	defer file.Close()

	// Load the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return "", fmt.Errorf("Error parsing HTML: %v", err)
	}

	// Remove unwanted tags
	doc.Find("script, style, noscript, iframe, svg, link, meta").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	// Remove unwanted attributes based on values that indicate ads or tracking
	doc.Find("[class], [id]").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		id, _ := s.Attr("id")
		if strings.Contains(strings.ToLower(class), "advertisement") || strings.Contains(strings.ToLower(id), "advertisement") {
			s.Remove()
		}
	})

	// Convert the cleaned HTML document back to a string
	htmlContent, err := doc.Html()
	if err != nil {
		return "", fmt.Errorf("Error getting HTML content from document: %v", err)
	}

	var buffer bytes.Buffer
	htmlWriter := html.NewTokenizer(strings.NewReader(htmlContent))
	for {
		tt := htmlWriter.Next()
		if tt == html.ErrorToken {
			break
		}
		buffer.Write(htmlWriter.Raw())
	}

	// Write cleaned HTML to a new file
	cleanedFileName := "cleaned.html"
	err = os.WriteFile(cleanedFileName, buffer.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("Error writing cleaned HTML to file: %v", err)
	}

	return cleanedFileName, nil
}
