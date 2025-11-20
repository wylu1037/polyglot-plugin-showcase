package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

const RUN_PATH = "/Users/wenyanglu/Workspace/github/polyglot-plugin-showcase/host-server/bin/plugins/builtin/data-processing/desensitization/v1.0.0/darwin_arm64/plugin"

func main() {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: common.Handshake,
		Plugins: map[string]plugin.Plugin{
			"desensitization": &common.PluginGRPCPlugin{},
		},
		Cmd: exec.Command(RUN_PATH),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	raw, err := rpcClient.Dispense("desensitization")
	if err != nil {
		log.Fatal(err)
	}

	pluginInterface := raw.(common.PluginInterface)
	metadata, err := pluginInterface.GetMetadata()
	if err != nil {
		log.Fatalf("❌ Failed to get metadata: %v", err)
	}

	fmt.Println("🔌 Plugin Information:")
	fmt.Printf("Name: %s\n", metadata.Name)
	fmt.Printf("Version: %s\n", metadata.Version)
	fmt.Printf("Description: %s\n", metadata.Description)
	fmt.Printf("Protocol Version: %d\n", metadata.ProtocolVersion)
	fmt.Printf("Available Methods: %v\n\n", metadata.Methods)

	fmt.Println("🧪 Data Desensitization Plugin Test:")

	testName := "John Doe"
	result, err := pluginInterface.Execute("DesensitizeName", map[string]string{"data": testName})
	if err != nil {
		log.Printf("❌ Error desensitizing name: %v", err)
	} else if !result.Success {
		log.Printf("❌ Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("👤 Name: \n  Original: %s\n  Desensitized: %s\n\n", testName, *result.Result)
	}

	testPhone := "13812345678"
	result, err = pluginInterface.Execute("DesensitizeTelNo", map[string]string{"data": testPhone})
	if err != nil {
		log.Printf("❌ Error desensitizing phone: %v", err)
	} else if !result.Success {
		log.Printf("❌ Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("📞 Phone Number: \n  Original: %s\n  Desensitized: %s\n\n", testPhone, *result.Result)
	}

	testID := "110101199001011234"
	result, err = pluginInterface.Execute("DesensitizeIDNumber", map[string]string{"data": testID})
	if err != nil {
		log.Printf("❌ Error desensitizing ID: %v", err)
	} else if !result.Success {
		log.Printf("❌ Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("🆔 ID Number: \n  Original: %s\n  Desensitized: %s\n\n", testID, *result.Result)
	}

	testEmail := "user@example.com"
	result, err = pluginInterface.Execute("DesensitizeEmail", map[string]string{"data": testEmail})
	if err != nil {
		log.Printf("❌ Error desensitizing email: %v", err)
	} else if !result.Success {
		log.Printf("❌ Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("📧 Email: \n  Original: %s\n  Desensitized: %s\n\n", testEmail, *result.Result)
	}

	testCard := "6222021234567890123"
	result, err = pluginInterface.Execute("DesensitizeBankCard", map[string]string{"data": testCard})
	if err != nil {
		log.Printf("❌ Error desensitizing bank card: %v", err)
	} else if !result.Success {
		log.Printf("❌ Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("💳 Bank Card: \n  Original: %s\n  Desensitized: %s\n\n", testCard, *result.Result)
	}

	testAddress := "123 Some Street, Chaoyang District, Beijing"
	result, err = pluginInterface.Execute("DesensitizeAddress", map[string]string{"data": testAddress})
	if err != nil {
		log.Printf("❌ Error desensitizing address: %v", err)
	} else if !result.Success {
		log.Printf("❌ Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("🏠 Address: \n  Original: %s\n  Desensitized: %s\n\n", testAddress, *result.Result)
	}

	fmt.Println("✅ All tests completed.")

	if !client.Exited() {
		fmt.Println("\n🚀 Plugin process is running normally.")
	} else {
		fmt.Println("\n👋 Plugin process has exited.")
	}
}
