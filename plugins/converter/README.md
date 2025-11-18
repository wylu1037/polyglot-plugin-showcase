# Data Converter Plugin

A polyglot plugin for converting JSON data (typically from database queries) into various formats including CSV, TXT, and HTML.

## Features

### Supported Conversions

1. **JSON to CSV**
   - Converts JSON objects or arrays to CSV format
   - Supports custom delimiters
   - Automatically extracts headers from object keys
   - Handles missing fields gracefully

2. **JSON to TXT**
   - Two format options:
     - `key-value`: Human-readable key-value pairs with indentation
     - `json-pretty`: Pretty-printed JSON with proper formatting
   - Supports nested objects and arrays

3. **JSON to HTML**
   - Generates HTML tables from JSON data
   - Optional styling with modern CSS
   - Can generate full HTML page or just table fragment
   - Automatic HTML escaping for security

## Usage

### As a Plugin

Build the plugin:

```bash
cd plugins/converter
go build -o ../../host-server/bin/plugins/converter/converter_v1.0.0 .
```

The plugin will be automatically discovered by the host server.

### Direct Usage (Library)

```go
import "github.com/wylu1037/polyglot-plugin-showcase/plugins/converter/impl"

converter := &impl.ConverterImpl{}

// Convert to CSV
jsonData := `[{"name":"John","age":"30"},{"name":"Jane","age":"25"}]`
csv, err := converter.ConvertToCSV(jsonData, map[string]string{})

// Convert to TXT
txt, err := converter.ConvertToTXT(jsonData, map[string]string{
    "format": "key-value",
})

// Convert to HTML
html, err := converter.ConvertToHTML(jsonData, map[string]string{
    "styled": "true",
    "full_page": "true",
})
```

## API Reference

### ConvertToCSV

Converts JSON data to CSV format.

**Parameters:**
- `data`: JSON string (object or array of objects)
- `options`:
  - `delimiter`: Custom delimiter (default: `,`)

**Example:**
```json
Input: [{"name":"John","age":"30"},{"name":"Jane","age":"25"}]
Output:
age,name
30,John
25,Jane
```

### ConvertToTXT

Converts JSON data to plain text format.

**Parameters:**
- `data`: JSON string (any valid JSON)
- `options`:
  - `format`: `key-value` (default) or `json-pretty`

**Example (key-value):**
```json
Input: {"user":{"name":"John","age":"30"},"active":true}
Output:
active: true
user:
  age: 30
  name: John
```

**Example (json-pretty):**
```json
Input: {"name":"John","age":"30"}
Output:
{
  "age": "30",
  "name": "John"
}
```

### ConvertToHTML

Converts JSON data to HTML table format.

**Parameters:**
- `data`: JSON string (object or array of objects)
- `options`:
  - `styled`: `true` (default) or `false` - Apply CSS styling
  - `full_page`: `true` or `false` (default) - Generate complete HTML page

**Example:**
```json
Input: [{"name":"John","age":"30"}]
Output:
<table class="data-table">
  <thead>
    <tr>
      <th>age</th>
      <th>name</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>30</td>
      <td>John</td>
    </tr>
  </tbody>
</table>
```

## Running the Example

```bash
cd plugins/converter/example
go run main.go
```

This will demonstrate all conversion formats with sample data.

## Testing

Run the tests:

```bash
cd plugins/converter
go test ./impl/...
```

## Use Cases

1. **Database Export**: Convert database query results (JSON) to CSV for Excel import
2. **Report Generation**: Generate HTML reports from JSON data
3. **Log Processing**: Convert structured JSON logs to readable text format
4. **Data Migration**: Transform data between different format requirements
5. **API Response Formatting**: Convert API JSON responses to user-friendly formats

## Plugin Metadata

- **Name**: converter
- **Version**: 1.0.0
- **Type**: converter
- **Input Format**: JSON
- **Output Formats**: CSV, TXT, HTML

## Architecture

```
converter/
├── main.go              # Plugin entry point
├── adapter/
│   └── adapter.go       # Common plugin interface adapter
├── impl/
│   ├── converter.go     # Core conversion logic
│   └── converter_test.go # Unit tests
├── example/
│   └── main.go          # Usage examples
└── README.md
```

## Error Handling

The plugin provides detailed error messages for:
- Empty or invalid JSON data
- Unsupported data structures
- Invalid format options
- Conversion failures

All errors are returned through the standard plugin interface with descriptive messages.

