package extraction

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

type RawSection struct {
	Tag        string            `json:"tag"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Content    string            `json:"content,omitempty"`
	Children   []RawSection      `json:"children,omitempty"`
}

func ExtractRawSections(r io.Reader) ([]RawSection, error) {
	// Load HTML document into GoQuery
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("Error loading HTML with goquery: %v", err)
	}

	// Clean up the document by removing non-content nodes
	doc.Find("script, style, nav, footer, header, noscript").Remove()

	// Traverse the cleaned document to extract meaningful content
	var sections []RawSection
	doc.Find("body").Children().Each(func(index int, sel *goquery.Selection) {
		section := traverseSelection(sel)
		if section.Content != "" || len(section.Children) > 0 {
			sections = append(sections, section)
		}
	})

	return sections, nil
}

// traverseSelection recursively traverses a GoQuery selection to extract content
func traverseSelection(sel *goquery.Selection) RawSection {
	tag := goquery.NodeName(sel)
	if tag == "" {
		return RawSection{}
	}

	section := RawSection{
		Tag:        tag,
		Attributes: extractAttributes(sel),
	}

	// Extract content if it's a text node or meaningful element
	content := sel.Text()
	if content != "" {
		section.Content = content
	}

	// Recursively process child elements
	sel.Children().Each(func(index int, childSel *goquery.Selection) {
		childSection := traverseSelection(childSel)
		if childSection.Content != "" || len(childSection.Children) > 0 {
			section.Children = append(section.Children, childSection)
		}
	})

	return section
}

// extractAttributes extracts attributes from a goquery selection
func extractAttributes(sel *goquery.Selection) map[string]string {
	attributes := make(map[string]string)
	for _, attr := range sel.Nodes[0].Attr {
		attributes[attr.Key] = attr.Val
	}
	return attributes
}
