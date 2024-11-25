package processing

import (
    "fmt"
    "html-text-extraction/internal/extraction"
    "strings"
)

type IntermediateSection struct {
    Tag      string                `json:"tag"`
    Content  string                `json:"content,omitempty"`
    Children []IntermediateSection `json:"children,omitempty"`
}

func TransformToIntermediateSections(rawSections []extraction.RawSection) []IntermediateSection {
    var intermediateSections []IntermediateSection

    for _, raw := range rawSections {
        section := IntermediateSection{
            Tag:     raw.Tag,
            Content: strings.TrimSpace(raw.Content),
        }

        // Recursively transform children
        if len(raw.Children) > 0 {
            section.Children = TransformToIntermediateSections(raw.Children)
        }

        // Only append if it has meaningful content or children
        if section.Content != "" || len(section.Children) > 0 {
            intermediateSections = append(intermediateSections, section)
        }
    }

    return intermediateSections
}

type Section struct {
    Heading     string    `json:"heading,omitempty"`
    Content     string    `json:"content,omitempty"`
    Subsections []Section `json:"subsections,omitempty"`
    ListItems   []string  `json:"list_items,omitempty"`
}

// New function to post-process intermediate sections
func PostProcessIntermediateSections(intermediateSections []IntermediateSection) []Section {
    var sections []Section

    for _, intermediate := range intermediateSections {
        section := Section{}

        if isHeading(intermediate.Tag) && intermediate.Content != "" {
            fmt.Printf("Processing heading: %s\n", intermediate.Content)
            section.Heading = intermediate.Content
        } else if isParagraph(intermediate.Tag) && intermediate.Content != "" {
            fmt.Printf("Processing paragraph at top-level: %s\n", intermediate.Content)
            section.Content = intermediate.Content
        } else if len(intermediate.Children) > 0 {
            // Recursively process children if current tag is not directly processed
            processIntermediateChildren(intermediate.Children, &section)
        }

        if isSectionValid(section) {
            sections = append(sections, section)
        } else {
            fmt.Printf("Skipping invalid section with tag: %s\n", intermediate.Tag)
        }
    }

    return sections
}

func processIntermediateChildren(children []IntermediateSection, currentSection *Section) {
    for _, child := range children {
        if isHeading(child.Tag) && child.Content != "" {
            fmt.Printf("Processing subheading: %s\n", child.Content)
            subsection := Section{Heading: child.Content}
            processIntermediateChildren(child.Children, &subsection)
            if isSectionValid(subsection) {
                currentSection.Subsections = append(currentSection.Subsections, subsection)
            }
        } else if isParagraph(child.Tag) && child.Content != "" {
            fmt.Printf("Adding paragraph content: %s\n", child.Content)
            currentSection.Content = strings.TrimSpace(currentSection.Content + " " + child.Content)
        } else if isList(child.Tag) {
            fmt.Println("Processing list items")
            var listItems []string
            for _, listItem := range child.Children {
                if listItem.Tag == "li" && listItem.Content != "" {
                    fmt.Printf("Adding list item: %s\n", listItem.Content)
                    listItems = append(listItems, listItem.Content)
                }
            }
            if len(listItems) > 0 {
                currentSection.ListItems = append(currentSection.ListItems, listItems...)
            }
        } else {
            // Recursively process other children
            processIntermediateChildren(child.Children, currentSection)
        }
    }
}

func isSectionValid(section Section) bool {
    return section.Heading != "" || section.Content != "" || len(section.Subsections) > 0 || len(section.ListItems) > 0
}

func isHeading(tag string) bool {
    headingTags := map[string]bool{
        "h1": true,
        "h2": true,
        "h3": true,
        "h4": true,
        "h5": true,
        "h6": true,
    }
    return headingTags[tag]
}

func isParagraph(tag string) bool {
    paragraphTags := map[string]bool{
        "p": true,
    }
    return paragraphTags[tag]
}

func isList(tag string) bool {
    listTags := map[string]bool{
        "ul": true,
        "ol": true,
    }
    return listTags[tag]
}
