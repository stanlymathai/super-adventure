package orchestration

import (
	"fmt"
	"html-text-extraction/internal/extraction"
	"html-text-extraction/internal/preprocessor"
	"html-text-extraction/internal/processing"
	"html-text-extraction/internal/utils"
	"io"
)

func Orchestrate(reader io.Reader, outputFile string) {
	// Step 1: Pre-Treat the HTML Content
	fmt.Println("Starting HTML pre-treatment...")
	cleanReader, err := preprocessor.CleanHTML(reader)
	if err != nil {
		fmt.Printf("Error during HTML pre-treatment: %v\n", err)
		return
	}

	// Step 2: Extract Raw Sections
	fmt.Println("Starting raw text extraction...")
	rawSections, err := extraction.ExtractRawSections(cleanReader)
	if err != nil {
		fmt.Printf("Error extracting text: %v\n", err)
		return
	}

	// Debugging output to verify rawSections content
	if len(rawSections) == 0 {
		fmt.Println("Warning: rawSections is empty after extraction. Check extraction logic.")
	} else {
		fmt.Printf("rawSections extraction success. Total sections extracted: %d\n", len(rawSections))
	}

	// Save raw output for further inspection
	if err := utils.WriteToFile("raw_result.json", rawSections); err != nil {
		fmt.Printf("Error writing raw JSON to file: %v\n", err)
		return
	}

	// Step 3: Transform to Intermediate Sections
	fmt.Println("Starting transformation to intermediate sections...")
	intermediateSections := processing.TransformToIntermediateSections(rawSections)

	if len(intermediateSections) == 0 {
		fmt.Println("Warning: intermediateSections is empty after transformation. Check transformation logic.")
	} else {
		fmt.Printf("intermediateSections transformation success. Total intermediate sections: %d\n", len(intermediateSections))
	}

	// Step 4: Post-Process Intermediate Sections to Create Structured Output
	fmt.Println("Starting post-processing...")
	structuredSections := processing.PostProcessIntermediateSections(intermediateSections)

	// Debugging output to verify structuredSections content
	if len(structuredSections) == 0 {
		fmt.Println("Warning: structuredSections is empty after post-processing. Check processing logic.")
	} else {
		fmt.Printf("structuredSections post-processing success. Total structured sections: %d\n", len(structuredSections))
	}

	// Save structured output to the final file
	if err := utils.WriteToFile(outputFile, structuredSections); err != nil {
		fmt.Printf("Error writing structured JSON to file: %v\n", err)
		return
	}

	fmt.Printf("Extraction complete. Data saved to %s\n", outputFile)
}
