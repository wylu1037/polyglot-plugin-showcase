package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

const RUN_PATH = "/Users/wenyanglu/Workspace/github/polyglot-plugin-showcase/host-server/bin/plugins/dpanonymizer/dpanonymizer_v1.0.0"

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
		log.Fatalf("Plugin binary not found at %s. Please build it first using 'make plugin-dpanonymizer'", absPath)
	}

	// Create plugin client
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: common.Handshake,
		Plugins: map[string]plugin.Plugin{
			"dpanonymizer": &common.PluginGRPCPlugin{},
		},
		Cmd:              exec.Command(absPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatalf("Failed to create RPC client: %v", err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("dpanonymizer")
	if err != nil {
		log.Fatalf("Failed to dispense plugin: %v", err)
	}

	// Cast to the plugin interface
	dpPlugin := raw.(common.PluginInterface)

	// Get metadata
	fmt.Println("=== Plugin Metadata ===")
	metadata, err := dpPlugin.GetMetadata()
	if err != nil {
		log.Fatalf("Failed to get metadata: %v", err)
	}
	fmt.Printf("Name: %s\n", metadata.Name)
	fmt.Printf("Version: %s\n", metadata.Version)
	fmt.Printf("Description: %s\n", metadata.Description)
	fmt.Printf("Methods: %v\n", metadata.Methods)
	fmt.Printf("Capabilities: %v\n\n", metadata.Capabilities)

	// Example 1: Add Laplace Noise
	fmt.Println("=== Example 1: Add Laplace Noise ===")
	fmt.Println("Original value: 100.0")
	fmt.Println("Epsilon: 1.0 (privacy budget)")
	fmt.Println("Sensitivity: 1.0")

	laplaceResult, err := dpPlugin.Execute("AddLaplaceNoise", map[string]string{
		"value":       "100.0",
		"epsilon":     "1.0",
		"sensitivity": "1.0",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if laplaceResult.Success {
		fmt.Printf("Noisy value: %s\n\n", *laplaceResult.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *laplaceResult.Error)
	}

	// Example 2: Add Gaussian Noise
	fmt.Println("=== Example 2: Add Gaussian Noise ===")
	fmt.Println("Original value: 100.0")
	fmt.Println("Epsilon: 1.0, Delta: 1e-5")
	fmt.Println("Sensitivity: 1.0")

	gaussianResult, err := dpPlugin.Execute("AddGaussianNoise", map[string]string{
		"value":       "100.0",
		"epsilon":     "1.0",
		"delta":       "0.00001",
		"sensitivity": "1.0",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if gaussianResult.Success {
		fmt.Printf("Noisy value: %s\n\n", *gaussianResult.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *gaussianResult.Error)
	}

	// Example 3: Differentially Private Count
	fmt.Println("=== Example 3: Differentially Private Count ===")
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	valuesJSON, _ := json.Marshal(values)
	fmt.Printf("Counting %d values\n", len(values))
	fmt.Println("Epsilon: 1.0, Delta: 1e-5")

	countResult, err := dpPlugin.Execute("DPCount", map[string]string{
		"values":                     string(valuesJSON),
		"epsilon":                    "1.0",
		"delta":                      "0.00001",
		"max_partitions_contributed": "1",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if countResult.Success {
		fmt.Printf("Noisy count: %s (actual: %d)\n\n", *countResult.Result, len(values))
	} else {
		fmt.Printf("Error: %s\n\n", *countResult.Error)
	}

	// Example 4: Differentially Private Sum
	fmt.Println("=== Example 4: Differentially Private Sum ===")
	sumValues := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
	sumValuesJSON, _ := json.Marshal(sumValues)
	actualSum := 0.0
	for _, v := range sumValues {
		actualSum += v
	}
	fmt.Printf("Summing values: %v (actual sum: %.2f)\n", sumValues, actualSum)
	fmt.Println("Epsilon: 1.0, Delta: 1e-5")
	fmt.Println("Bounds: [0, 100]")

	sumResult, err := dpPlugin.Execute("DPSum", map[string]string{
		"values":                     string(sumValuesJSON),
		"epsilon":                    "1.0",
		"delta":                      "0.00001",
		"lower_bound":                "0.0",
		"upper_bound":                "100.0",
		"max_partitions_contributed": "1",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if sumResult.Success {
		fmt.Printf("Noisy sum: %s\n\n", *sumResult.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *sumResult.Error)
	}

	// Example 5: Differentially Private Mean
	fmt.Println("=== Example 5: Differentially Private Mean ===")
	meanValues := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
	meanValuesJSON, _ := json.Marshal(meanValues)
	actualMean := 0.0
	for _, v := range meanValues {
		actualMean += v
	}
	actualMean /= float64(len(meanValues))
	fmt.Printf("Computing mean of: %v (actual mean: %.2f)\n", meanValues, actualMean)
	fmt.Println("Epsilon: 1.0, Delta: 1e-5")
	fmt.Println("Bounds: [0, 100]")

	meanResult, err := dpPlugin.Execute("DPMean", map[string]string{
		"values":                     string(meanValuesJSON),
		"epsilon":                    "1.0",
		"delta":                      "0.00001",
		"lower_bound":                "0.0",
		"upper_bound":                "100.0",
		"max_partitions_contributed": "1",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if meanResult.Success {
		fmt.Printf("Noisy mean: %s\n\n", *meanResult.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *meanResult.Error)
	}

	// Example 6: Differentially Private Variance
	fmt.Println("=== Example 6: Differentially Private Variance ===")
	varianceValues := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
	varianceValuesJSON, _ := json.Marshal(varianceValues)

	// Calculate actual variance
	mean := 0.0
	for _, v := range varianceValues {
		mean += v
	}
	mean /= float64(len(varianceValues))
	actualVariance := 0.0
	for _, v := range varianceValues {
		actualVariance += (v - mean) * (v - mean)
	}
	actualVariance /= float64(len(varianceValues))

	fmt.Printf("Computing variance of: %v (actual variance: %.2f)\n", varianceValues, actualVariance)
	fmt.Println("Epsilon: 1.0, Delta: 1e-5")
	fmt.Println("Bounds: [0, 100]")

	varianceResult, err := dpPlugin.Execute("DPVariance", map[string]string{
		"values":                     string(varianceValuesJSON),
		"epsilon":                    "1.0",
		"delta":                      "0.00001",
		"lower_bound":                "0.0",
		"upper_bound":                "100.0",
		"max_partitions_contributed": "1",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else if varianceResult.Success {
		fmt.Printf("Noisy variance: %s\n\n", *varianceResult.Result)
	} else {
		fmt.Printf("Error: %s\n\n", *varianceResult.Error)
	}

	fmt.Println("=== All examples completed ===")
}
