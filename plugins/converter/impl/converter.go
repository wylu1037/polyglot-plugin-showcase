// Package impl provides the implementation of the data converter plugin.
package impl

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"sort"
	"strings"
)

// ConverterImpl is the concrete implementation of the Converter interface.
type ConverterImpl struct{}

// ConvertToCSV converts JSON data to CSV format.
// Supports both single object and array of objects.
// Example: [{"name":"John","age":"30"},{"name":"Jane","age":"25"}] -> CSV with headers
func (c *ConverterImpl) ConvertToCSV(jsonData string, options map[string]string) (string, error) {
	if jsonData == "" {
		return "", errors.New("JSON data cannot be empty")
	}

	// Parse JSON data
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", fmt.Errorf("invalid JSON data: %w", err)
	}

	// Convert to array of maps
	var records []map[string]interface{}
	switch v := data.(type) {
	case map[string]interface{}:
		records = []map[string]interface{}{v}
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				records = append(records, m)
			} else {
				return "", errors.New("JSON array must contain objects")
			}
		}
	default:
		return "", errors.New("JSON data must be an object or array of objects")
	}

	if len(records) == 0 {
		return "", errors.New("no data to convert")
	}

	// Extract headers (keys from first record, sorted for consistency)
	headers := make([]string, 0, len(records[0]))
	for key := range records[0] {
		headers = append(headers, key)
	}
	sort.Strings(headers)

	// Check if custom delimiter is specified
	delimiter := ","
	if delim, ok := options["delimiter"]; ok && delim != "" {
		delimiter = delim
	}

	// Build CSV
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// Set custom delimiter if specified
	if len(delimiter) > 0 {
		writer.Comma = rune(delimiter[0])
	}

	// Write headers
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Write data rows
	for _, record := range records {
		row := make([]string, len(headers))
		for i, header := range headers {
			if val, ok := record[header]; ok {
				row[i] = fmt.Sprintf("%v", val)
			} else {
				row[i] = ""
			}
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return builder.String(), nil
}

// ConvertToTXT converts JSON data to plain text format.
// Creates a human-readable key-value format.
// Example: {"name":"John","age":"30"} -> "name: John\nage: 30"
func (c *ConverterImpl) ConvertToTXT(jsonData string, options map[string]string) (string, error) {
	if jsonData == "" {
		return "", errors.New("JSON data cannot be empty")
	}

	// Parse JSON data
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", fmt.Errorf("invalid JSON data: %w", err)
	}

	// Check format option
	format := options["format"]
	if format == "" {
		format = "key-value" // default format
	}

	var builder strings.Builder

	switch format {
	case "key-value":
		if err := c.formatKeyValue(&builder, data, 0); err != nil {
			return "", err
		}
	case "json-pretty":
		// Pretty print JSON
		prettyJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to format JSON: %w", err)
		}
		builder.Write(prettyJSON)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	return builder.String(), nil
}

// formatKeyValue recursively formats data as key-value pairs
func (c *ConverterImpl) formatKeyValue(builder *strings.Builder, data interface{}, indent int) error {
	indentStr := strings.Repeat("  ", indent)

	switch v := data.(type) {
	case map[string]interface{}:
		// Sort keys for consistent output
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			val := v[key]
			switch val.(type) {
			case map[string]interface{}, []interface{}:
				builder.WriteString(fmt.Sprintf("%s%s:\n", indentStr, key))
				if err := c.formatKeyValue(builder, val, indent+1); err != nil {
					return err
				}
			default:
				builder.WriteString(fmt.Sprintf("%s%s: %v\n", indentStr, key, val))
			}
		}

	case []interface{}:
		for i, item := range v {
			builder.WriteString(fmt.Sprintf("%s[%d]:\n", indentStr, i))
			if err := c.formatKeyValue(builder, item, indent+1); err != nil {
				return err
			}
		}

	default:
		builder.WriteString(fmt.Sprintf("%s%v\n", indentStr, v))
	}

	return nil
}

// ConvertToHTML converts JSON data to HTML table format.
// Creates a styled HTML table with headers and data rows.
// Example: [{"name":"John","age":"30"}] -> HTML table
func (c *ConverterImpl) ConvertToHTML(jsonData string, options map[string]string) (string, error) {
	if jsonData == "" {
		return "", errors.New("JSON data cannot be empty")
	}

	// Parse JSON data
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", fmt.Errorf("invalid JSON data: %w", err)
	}

	// Convert to array of maps
	var records []map[string]interface{}
	switch v := data.(type) {
	case map[string]interface{}:
		records = []map[string]interface{}{v}
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				records = append(records, m)
			} else {
				return "", errors.New("JSON array must contain objects")
			}
		}
	default:
		return "", errors.New("JSON data must be an object or array of objects")
	}

	if len(records) == 0 {
		return "", errors.New("no data to convert")
	}

	// Extract headers (keys from first record, sorted for consistency)
	headers := make([]string, 0, len(records[0]))
	for key := range records[0] {
		headers = append(headers, key)
	}
	sort.Strings(headers)

	// Check if styled option is enabled (default: true)
	styled := true
	if style, ok := options["styled"]; ok && style == "false" {
		styled = false
	}

	// Build HTML
	var builder strings.Builder

	// Add DOCTYPE and basic structure if full_page option is enabled
	if fullPage, ok := options["full_page"]; ok && fullPage == "true" {
		builder.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
		builder.WriteString("  <meta charset=\"UTF-8\">\n")
		builder.WriteString("  <title>Data Table</title>\n")
		if styled {
			c.writeHTMLStyles(&builder)
		}
		builder.WriteString("</head>\n<body>\n")
	}

	// Start table
	if styled {
		builder.WriteString("<table class=\"data-table\">\n")
	} else {
		builder.WriteString("<table>\n")
	}

	// Write headers
	builder.WriteString("  <thead>\n    <tr>\n")
	for _, header := range headers {
		builder.WriteString(fmt.Sprintf("      <th>%s</th>\n", html.EscapeString(header)))
	}
	builder.WriteString("    </tr>\n  </thead>\n")

	// Write data rows
	builder.WriteString("  <tbody>\n")
	for _, record := range records {
		builder.WriteString("    <tr>\n")
		for _, header := range headers {
			val := ""
			if v, ok := record[header]; ok {
				val = fmt.Sprintf("%v", v)
			}
			builder.WriteString(fmt.Sprintf("      <td>%s</td>\n", html.EscapeString(val)))
		}
		builder.WriteString("    </tr>\n")
	}
	builder.WriteString("  </tbody>\n")

	// Close table
	builder.WriteString("</table>\n")

	// Close HTML structure if full_page option is enabled
	if fullPage, ok := options["full_page"]; ok && fullPage == "true" {
		builder.WriteString("</body>\n</html>")
	}

	return builder.String(), nil
}

// writeHTMLStyles writes CSS styles for the HTML table
func (c *ConverterImpl) writeHTMLStyles(builder *strings.Builder) {
	builder.WriteString("  <style>\n")
	builder.WriteString("    .data-table {\n")
	builder.WriteString("      border-collapse: collapse;\n")
	builder.WriteString("      width: 100%;\n")
	builder.WriteString("      font-family: Arial, sans-serif;\n")
	builder.WriteString("      box-shadow: 0 2px 4px rgba(0,0,0,0.1);\n")
	builder.WriteString("    }\n")
	builder.WriteString("    .data-table th {\n")
	builder.WriteString("      background-color: #4CAF50;\n")
	builder.WriteString("      color: white;\n")
	builder.WriteString("      padding: 12px;\n")
	builder.WriteString("      text-align: left;\n")
	builder.WriteString("      font-weight: bold;\n")
	builder.WriteString("    }\n")
	builder.WriteString("    .data-table td {\n")
	builder.WriteString("      padding: 10px;\n")
	builder.WriteString("      border-bottom: 1px solid #ddd;\n")
	builder.WriteString("    }\n")
	builder.WriteString("    .data-table tr:hover {\n")
	builder.WriteString("      background-color: #f5f5f5;\n")
	builder.WriteString("    }\n")
	builder.WriteString("    .data-table tr:nth-child(even) {\n")
	builder.WriteString("      background-color: #f9f9f9;\n")
	builder.WriteString("    }\n")
	builder.WriteString("  </style>\n")
}
