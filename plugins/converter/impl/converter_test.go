package impl

import (
	"strings"
	"testing"
)

func TestConverterImpl_ConvertToCSV(t *testing.T) {
	converter := &ConverterImpl{}

	tests := []struct {
		name      string
		jsonData  string
		options   map[string]string
		wantErr   bool
		checkFunc func(string) bool
	}{
		{
			name:     "single object",
			jsonData: `{"name":"John","age":"30","city":"New York"}`,
			options:  map[string]string{},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "age,city,name") &&
					strings.Contains(result, "30,New York,John")
			},
		},
		{
			name:     "array of objects",
			jsonData: `[{"name":"John","age":"30"},{"name":"Jane","age":"25"}]`,
			options:  map[string]string{},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "age,name") &&
					strings.Contains(result, "30,John") &&
					strings.Contains(result, "25,Jane")
			},
		},
		{
			name:     "custom delimiter",
			jsonData: `{"name":"John","age":"30"}`,
			options:  map[string]string{"delimiter": ";"},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "age;name") &&
					strings.Contains(result, "30;John")
			},
		},
		{
			name:     "empty data",
			jsonData: "",
			options:  map[string]string{},
			wantErr:  true,
		},
		{
			name:     "invalid JSON",
			jsonData: `{invalid}`,
			options:  map[string]string{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ConvertToCSV(tt.jsonData, tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkFunc != nil && !tt.checkFunc(result) {
				t.Errorf("ConvertToCSV() result validation failed, got:\n%s", result)
			}
		})
	}
}

func TestConverterImpl_ConvertToTXT(t *testing.T) {
	converter := &ConverterImpl{}

	tests := []struct {
		name      string
		jsonData  string
		options   map[string]string
		wantErr   bool
		checkFunc func(string) bool
	}{
		{
			name:     "key-value format",
			jsonData: `{"name":"John","age":"30"}`,
			options:  map[string]string{"format": "key-value"},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "age: 30") &&
					strings.Contains(result, "name: John")
			},
		},
		{
			name:     "json-pretty format",
			jsonData: `{"name":"John","age":"30"}`,
			options:  map[string]string{"format": "json-pretty"},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "\"name\"") &&
					strings.Contains(result, "\"age\"")
			},
		},
		{
			name:     "nested object",
			jsonData: `{"user":{"name":"John","age":"30"},"active":true}`,
			options:  map[string]string{"format": "key-value"},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "user:") &&
					strings.Contains(result, "name: John")
			},
		},
		{
			name:     "empty data",
			jsonData: "",
			options:  map[string]string{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ConvertToTXT(tt.jsonData, tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToTXT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkFunc != nil && !tt.checkFunc(result) {
				t.Errorf("ConvertToTXT() result validation failed, got:\n%s", result)
			}
		})
	}
}

func TestConverterImpl_ConvertToHTML(t *testing.T) {
	converter := &ConverterImpl{}

	tests := []struct {
		name      string
		jsonData  string
		options   map[string]string
		wantErr   bool
		checkFunc func(string) bool
	}{
		{
			name:     "single object",
			jsonData: `{"name":"John","age":"30"}`,
			options:  map[string]string{},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "<table") &&
					strings.Contains(result, "<th>age</th>") &&
					strings.Contains(result, "<th>name</th>") &&
					strings.Contains(result, "<td>30</td>") &&
					strings.Contains(result, "<td>John</td>")
			},
		},
		{
			name:     "array of objects",
			jsonData: `[{"name":"John","age":"30"},{"name":"Jane","age":"25"}]`,
			options:  map[string]string{},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "<table") &&
					strings.Contains(result, "John") &&
					strings.Contains(result, "Jane")
			},
		},
		{
			name:     "styled table",
			jsonData: `{"name":"John"}`,
			options:  map[string]string{"styled": "true"},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "class=\"data-table\"")
			},
		},
		{
			name:     "unstyled table",
			jsonData: `{"name":"John"}`,
			options:  map[string]string{"styled": "false"},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "<table>") &&
					!strings.Contains(result, "class=")
			},
		},
		{
			name:     "full page",
			jsonData: `{"name":"John"}`,
			options:  map[string]string{"full_page": "true"},
			wantErr:  false,
			checkFunc: func(result string) bool {
				return strings.Contains(result, "<!DOCTYPE html>") &&
					strings.Contains(result, "<html>") &&
					strings.Contains(result, "</html>")
			},
		},
		{
			name:     "empty data",
			jsonData: "",
			options:  map[string]string{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ConvertToHTML(tt.jsonData, tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkFunc != nil && !tt.checkFunc(result) {
				t.Errorf("ConvertToHTML() result validation failed, got:\n%s", result)
			}
		})
	}
}

