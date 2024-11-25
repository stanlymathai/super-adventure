package preprocessor

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CleanHTML uses goquery to clean up the HTML and remove unnecessary elements.
func CleanHTML(reader io.Reader) (io.Reader, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("Error parsing HTML: %v", err)
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
	cleanedHTML, err := doc.Html()
	if err != nil {
		return nil, fmt.Errorf("Error converting cleaned HTML to string: %v", err)
	}

	// Create a reader from the cleaned HTML string
	return strings.NewReader(cleanedHTML), nil
}
