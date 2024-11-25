package extraction

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type RawSection struct {
	Tag        string            `json:"tag"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Content    string            `json:"content,omitempty"`
	Children   []RawSection      `json:"children,omitempty"`
}

func ExtractRawSections(r io.Reader) ([]RawSection, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	// Find the first meaningful element to start extraction from
	meaningfulRoot := findFirstMeaningfulElement(doc)
	if meaningfulRoot == nil {
		fmt.Println("Error: Could not find a meaningful root element.")
		return nil, nil
	}

	fmt.Printf("Found meaningful root: <%s>\n", meaningfulRoot.Data)
	fmt.Println("Starting traversal of meaningful HTML document...")

	root := traverseRaw(meaningfulRoot)

	if len(root.Children) == 0 {
		fmt.Println("Warning: No children extracted from the meaningful root element. Check HTML structure or tag filtering.")
	} else {
		fmt.Printf("Extracted %d top-level children.\n", len(root.Children))
	}

	return root.Children, nil
}

func traverseRaw(n *html.Node) RawSection {
	// Initialize section to hold the node's data
	section := RawSection{
		Tag:        n.Data,
		Attributes: extractAttributes(n),
	}

	// Extract content if it's a text node
	if n.Type == html.TextNode {
		text := strings.TrimSpace(strings.ReplaceAll(n.Data, "\n", " "))
		if text != "" {
			section.Content = text
			fmt.Printf("Extracted text: %s\n", text)
		}
	}

	// Traverse all relevant children nodes
	if n.Type == html.ElementNode {
		fmt.Printf("Processing element: <%s>\n", n.Data)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if isNonContentNode(c) {
				fmt.Printf("Skipping non-content node: <%s>\n", c.Data)
				continue
			}
			childSection := traverseRaw(c)
			if childSection.Tag != "" || childSection.Content != "" || len(childSection.Children) > 0 {
				section.Children = append(section.Children, childSection)
			}
		}
	}

	return section
}

// Helper function to find the first meaningful element, skipping any document-level wrappers
func findFirstMeaningfulElement(n *html.Node) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			return c
		}
	}
	return nil
}

// Helper function to skip non-content nodes such as scripts, styles, and comments
func isNonContentNode(n *html.Node) bool {
	if n.Type == html.CommentNode || n.Type == html.DoctypeNode {
		return true
	}

	if n.Type == html.ElementNode {
		nonContentTags := map[string]bool{
			"script":   true,
			"style":    true,
			"meta":     true,
			"link":     true,
			"noscript": true,
		}
		return nonContentTags[n.Data]
	}

	return false
}

func extractAttributes(n *html.Node) map[string]string {
	if len(n.Attr) == 0 {
		return nil
	}
	attributes := make(map[string]string)
	for _, attr := range n.Attr {
		attributes[attr.Key] = attr.Val
	}
	return attributes
}
