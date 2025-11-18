package main

import (
	"fmt"
	"log"

	"github.com/wylu1037/polyglot-plugin-showcase/plugins/converter/impl"
)

func main() {
	converter := &impl.ConverterImpl{}

	jsonData := `[
		{"id": "1", "name": "John Doe", "email": "john@example.com", "age": "30", "city": "New York"},
		{"id": "2", "name": "Jane Smith", "email": "jane@example.com", "age": "25", "city": "Los Angeles"},
		{"id": "3", "name": "Bob Johnson", "email": "bob@example.com", "age": "35", "city": "Chicago"}
	]`

	fmt.Println("=== Original JSON Data ===")
	fmt.Println(jsonData)
	fmt.Println()

	// Convert to CSV
	fmt.Println("=== Convert to CSV ===")
	csvResult, err := converter.ConvertToCSV(jsonData, map[string]string{})
	if err != nil {
		log.Fatalf("CSV conversion failed: %v", err)
	}
	fmt.Println(csvResult)
	fmt.Println()

	// Convert to CSV with custom delimiter
	fmt.Println("=== Convert to CSV (with semicolon delimiter) ===")
	csvWithDelimiter, err := converter.ConvertToCSV(jsonData, map[string]string{"delimiter": ";"})
	if err != nil {
		log.Fatalf("CSV conversion with delimiter failed: %v", err)
	}
	fmt.Println(csvWithDelimiter)
	fmt.Println()

	// Convert to TXT (key-value format)
	fmt.Println("=== Convert to TXT (key-value format) ===")
	txtResult, err := converter.ConvertToTXT(jsonData, map[string]string{"format": "key-value"})
	if err != nil {
		log.Fatalf("TXT conversion failed: %v", err)
	}
	fmt.Println(txtResult)
	fmt.Println()

	// Convert to TXT (pretty JSON format)
	fmt.Println("=== Convert to TXT (pretty JSON format) ===")
	txtPretty, err := converter.ConvertToTXT(jsonData, map[string]string{"format": "json-pretty"})
	if err != nil {
		log.Fatalf("TXT pretty conversion failed: %v", err)
	}
	fmt.Println(txtPretty)
	fmt.Println()

	// Convert to HTML (styled table)
	fmt.Println("=== Convert to HTML (styled table) ===")
	htmlResult, err := converter.ConvertToHTML(jsonData, map[string]string{"styled": "true"})
	if err != nil {
		log.Fatalf("HTML conversion failed: %v", err)
	}
	fmt.Println(htmlResult)
	fmt.Println()

	// Convert to HTML (full page)
	fmt.Println("=== Convert to HTML (full page) ===")
	htmlFullPage, err := converter.ConvertToHTML(jsonData, map[string]string{
		"styled":    "true",
		"full_page": "true",
	})
	if err != nil {
		log.Fatalf("HTML full page conversion failed: %v", err)
	}
	fmt.Println(htmlFullPage)
	fmt.Println()

	// Example with single object
	singleObject := `{"id": "1", "name": "John Doe", "email": "john@example.com"}`
	fmt.Println("=== Single Object to CSV ===")
	singleCSV, err := converter.ConvertToCSV(singleObject, map[string]string{})
	if err != nil {
		log.Fatalf("Single object CSV conversion failed: %v", err)
	}
	fmt.Println(singleCSV)

	// Example with nested object
	nestedObject := `{
		"user": {
			"name": "John Doe",
			"age": "30"
		},
		"active": true,
		"roles": ["admin", "user"]
	}`
	fmt.Println("=== Nested Object to TXT ===")
	nestedTXT, err := converter.ConvertToTXT(nestedObject, map[string]string{"format": "key-value"})
	if err != nil {
		log.Fatalf("Nested object TXT conversion failed: %v", err)
	}
	fmt.Println(nestedTXT)
}
