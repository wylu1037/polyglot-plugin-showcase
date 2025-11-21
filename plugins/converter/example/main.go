package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

const RUN_PATH = "/Users/wenyanglu/Workspace/github/polyglot-plugin-showcase/host-server/bin/plugins/builtin/data-processing/converter/v1.0.0/darwin_arm64/plugin"

func main() {
	// Get the plugin binary path
	pluginPath := os.Getenv("PLUGIN_PATH")
	if pluginPath == "" {
		pluginPath = RUN_PATH
	}

	// Make path absolute
	absPath, err := filepath.Abs(pluginPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// Check if plugin exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("Plugin binary not found at %s. Please build it first using 'make plugin-converter'", absPath)
	}

	// Create plugin client
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: common.Handshake,
		Plugins: map[string]plugin.Plugin{
			"converter": &common.PluginGRPCPlugin{},
		},
		Cmd:              exec.Command(absPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		AutoMTLS:         true,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatalf("Failed to create RPC client: %v", err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("converter")
	if err != nil {
		log.Fatalf("Failed to dispense plugin: %v", err)
	}

	// Cast to the plugin interface
	converterPlugin := raw.(common.PluginInterface)

	// Get metadata
	fmt.Println("=== Plugin Metadata ===")
	metadata, err := converterPlugin.GetMetadata()
	if err != nil {
		log.Fatalf("Failed to get metadata: %v", err)
	}
	fmt.Printf("Name: %s\n", metadata.Name)
	fmt.Printf("Version: %s\n", metadata.Version)
	fmt.Printf("Description: %s\n", metadata.Description)
	fmt.Printf("Methods: %v\n", metadata.Methods)
	fmt.Printf("Capabilities: %v\n\n", metadata.Capabilities)

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
	csvResult, err := converterPlugin.Execute("ConvertToCSV", map[string]string{"data": jsonData})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if csvResult.Success {
		fmt.Printf("%s\n\n", *csvResult.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *csvResult.Error)
	}

	// Convert to CSV with custom delimiter
	fmt.Println("=== Convert to CSV (with semicolon delimiter) ===")
	csvWithDelimiter, err := converterPlugin.Execute("ConvertToCSV", map[string]string{
		"data":      jsonData,
		"delimiter": ";",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if csvWithDelimiter.Success {
		fmt.Printf("%s\n\n", *csvWithDelimiter.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *csvWithDelimiter.Error)
	}

	// Convert to TXT (key-value format)
	fmt.Println("=== Convert to TXT (key-value format) ===")
	txtResult, err := converterPlugin.Execute("ConvertToTXT", map[string]string{
		"data":   jsonData,
		"format": "key-value",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if txtResult.Success {
		fmt.Printf("%s\n\n", *txtResult.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *txtResult.Error)
	}

	// Convert to TXT (pretty JSON format)
	fmt.Println("=== Convert to TXT (pretty JSON format) ===")
	txtPretty, err := converterPlugin.Execute("ConvertToTXT", map[string]string{
		"data":   jsonData,
		"format": "json-pretty",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if txtPretty.Success {
		fmt.Printf("%s\n\n", *txtPretty.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *txtPretty.Error)
	}

	// Convert to HTML (styled table)
	fmt.Println("=== Convert to HTML (styled table) ===")
	htmlResult, err := converterPlugin.Execute("ConvertToHTML", map[string]string{
		"data":   jsonData,
		"styled": "true",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if htmlResult.Success {
		fmt.Printf("%s\n\n", *htmlResult.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *htmlResult.Error)
	}

	// Convert to HTML (full page)
	fmt.Println("=== Convert to HTML (full page) ===")
	paramsMap := map[string]string{
		"data":      jsonData,
		"styled":    "true",
		"full_page": "true",
	}
	htmlFullPage, err := converterPlugin.Execute("ConvertToHTML", paramsMap)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if htmlFullPage.Success {
		fmt.Printf("%s\n\n", *htmlFullPage.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *htmlFullPage.Error)
	}

	// Example with single object
	singleObject := `{"id": "1", "name": "John Doe", "email": "john@example.com"}`
	fmt.Println("=== Single Object to CSV ===")
	singleCSV, err := converterPlugin.Execute("ConvertToCSV", map[string]string{"data": singleObject})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if singleCSV.Success {
		fmt.Printf("%s\n\n", *singleCSV.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *singleCSV.Error)
	}

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
	nestedTXT, err := converterPlugin.Execute("ConvertToTXT", map[string]string{
		"data":   nestedObject,
		"format": "key-value",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if nestedTXT.Success {
		fmt.Printf("%s\n\n", *nestedTXT.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *nestedTXT.Error)
	}

	fmt.Println("=== All examples completed ===")
}
