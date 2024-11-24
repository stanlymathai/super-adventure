package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input-html-file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := "result.json" // Output filename is fixed

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Println("Starting text extraction...")
	sections, err := extractText(file)
	if err != nil {
		fmt.Printf("Error extracting text: %v\n", err)
		return
	}
	fmt.Println("Text extraction completed.")

	jsonData, err := json.MarshalIndent(sections, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON to file: %v\n", err)
		return
	}

	fmt.Printf("Extraction complete. Data saved to %s\n", outputFile)
}

type Section struct {
	Heading        string       `json:"heading"`
	Subsections    []Subsection `json:"subsections,omitempty"`
	Content        []string     `json:"content,omitempty"`
	ContactInfo    *string      `json:"contact_info,omitempty"`
	SubmissionLink *string      `json:"submission_link,omitempty"`
	StudyTypes     []StudyType  `json:"study_types,omitempty"`
	RequiredDocs   []string     `json:"required_documents,omitempty"`
	Guidelines     *Guidelines  `json:"guidelines,omitempty"`
	ListItems      []string     `json:"list_items,omitempty"`
}

type Subsection struct {
	Subheading string   `json:"subheading"`
	Content    []string `json:"content"`
	ListItems  []string `json:"list_items,omitempty"`
}

type StudyType struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Guidelines struct {
	Preparation []string `json:"preparation,omitempty"`
	Length      string   `json:"length,omitempty"`
}

func extractText(r io.Reader) ([]Section, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var sections []Section
	var currentSection *Section
	var currentSubheading string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		// Skip non-content nodes
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style" || n.Data == "nav" || n.Data == "footer") {
			return
		}

		// Skip elements that are not useful (e.g., Back to top, Sign In)
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "button") {
			if hasIrrelevantText(n) {
				return
			}
		}

		// Extract text from heading nodes
		if n.Type == html.ElementNode && isHeading(n.Data) {
			var buf strings.Builder
			getTextContent(n, &buf)
			heading := strings.TrimSpace(buf.String())
			if heading != "" {
				if currentSection != nil {
					sections = append(sections, *currentSection)
				}
				currentSection = &Section{Heading: heading}
				currentSubheading = ""
				fmt.Printf("Processing heading: %s\n", heading)
			}
		}

		// Extract text from subheading nodes
		if n.Type == html.ElementNode && isSubheading(n.Data) {
			var buf strings.Builder
			getTextContent(n, &buf)
			subheading := strings.TrimSpace(buf.String())
			if subheading != "" && currentSection != nil {
				currentSubheading = subheading
				currentSection.Subsections = append(currentSection.Subsections, Subsection{Subheading: subheading})
				fmt.Printf("Processing subheading: %s\n", subheading)
			}
		}

		// Extract text from paragraph nodes
		if n.Type == html.ElementNode && isParagraph(n.Data) {
			var buf strings.Builder
			getTextContent(n, &buf)
			text := strings.TrimSpace(buf.String())
			if text != "" && currentSection != nil {
				if strings.Contains(strings.ToLower(text), "contact") {
					currentSection.ContactInfo = &text
				} else if strings.Contains(strings.ToLower(text), "submit") {
					link := extractLink(n)
					currentSection.SubmissionLink = &link
				} else if strings.Contains(strings.ToLower(text), "study type") {
					currentSection.StudyTypes = append(currentSection.StudyTypes, StudyType{Type: text, Description: ""})
				} else {
					if currentSubheading != "" {
						for i := range currentSection.Subsections {
							if currentSection.Subsections[i].Subheading == currentSubheading {
								currentSection.Subsections[i].Content = append(currentSection.Subsections[i].Content, text)
								break
							}
						}
					} else {
						currentSection.Content = append(currentSection.Content, text)
					}
				}
			}
		}

		// Extract text from list items
		if n.Type == html.ElementNode && (n.Data == "ul" || n.Data == "ol") {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "li" {
					var buf strings.Builder
					getTextContent(c, &buf)
					listItem := strings.TrimSpace(buf.String())
					if listItem != "" && currentSection != nil {
						if currentSubheading != "" {
							for i := range currentSection.Subsections {
								if currentSection.Subsections[i].Subheading == currentSubheading {
									currentSection.Subsections[i].ListItems = append(currentSection.Subsections[i].ListItems, listItem)
									break
								}
							}
						} else if strings.Contains(strings.ToLower(listItem), "required document") {
							currentSection.RequiredDocs = append(currentSection.RequiredDocs, listItem)
						} else {
							currentSection.ListItems = append(currentSection.ListItems, listItem)
						}
					}
				}
			}
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	fmt.Println("Traversing HTML document...")
	traverse(doc)
	fmt.Println("Traversal completed.")

	if currentSection != nil {
		sections = append(sections, *currentSection)
	}

	return sections, nil
}

func isHeading(tag string) bool {
	headingTags := map[string]bool{
		"h1": true,
		"h2": true,
		"h3": true,
	}
	return headingTags[tag]
}

func isSubheading(tag string) bool {
	subheadingTags := map[string]bool{
		"h4": true,
		"h5": true,
		"h6": true,
	}
	return subheadingTags[tag]
}

func isParagraph(tag string) bool {
	paragraphTags := map[string]bool{
		"p":    true,
		"span": true,
	}
	return paragraphTags[tag]
}

func getTextContent(n *html.Node, buf *strings.Builder) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getTextContent(c, buf)
	}
}

func hasIrrelevantText(n *html.Node) bool {
	var buf strings.Builder
	getTextContent(n, &buf)
	text := strings.TrimSpace(strings.ToLower(buf.String()))
	irrelevantTexts := []string{"back to top", "sign in", "full info"}
	for _, ir := range irrelevantTexts {
		if strings.Contains(text, ir) {
			return true
		}
	}
	return false
}

func extractLink(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
