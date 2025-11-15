package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

func main() {
	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: common.Handshake,
		Plugins: map[string]plugin.Plugin{
			"desensitization": &common.PluginGRPCPlugin{},
		},
		Cmd: exec.Command("/Users/wenyanglu/Workspace/github/polyglot-plugin-showcase/host-server/bin/plugins/desensitization/desensitization_v1.0.0"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin using the plugin type name
	// This must match the key used in the plugin's Serve() call
	raw, err := rpcClient.Dispense("desensitization")
	if err != nil {
		log.Fatal(err)
	}

	// We should have a PluginInterface now! This uses the common interface.
	pluginInterface := raw.(common.PluginInterface)

	// Get plugin metadata
	metadata, err := pluginInterface.GetMetadata()
	if err != nil {
		log.Fatalf("Failed to get metadata: %v", err)
	}

	fmt.Println("=== 插件信息 ===")
	fmt.Printf("名称: %s\n", metadata.Name)
	fmt.Printf("版本: %s\n", metadata.Version)
	fmt.Printf("描述: %s\n", metadata.Description)
	fmt.Printf("协议版本: %d\n", metadata.ProtocolVersion)
	fmt.Printf("可用方法: %v\n\n", metadata.Methods)

	fmt.Println("=== 数据脱敏插件测试 ===")

	// Test DesensitizeName
	testName := "张三"
	result, err := pluginInterface.Execute("DesensitizeName", map[string]string{"data": testName})
	if err != nil {
		log.Printf("Error desensitizing name: %v", err)
	} else if !result.Success {
		log.Printf("Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("姓名脱敏:\n  原始: %s\n  脱敏: %s\n\n", testName, *result.Result)
	}

	// Test DesensitizeTelNo
	testPhone := "13812345678"
	result, err = pluginInterface.Execute("DesensitizeTelNo", map[string]string{"data": testPhone})
	if err != nil {
		log.Printf("Error desensitizing phone: %v", err)
	} else if !result.Success {
		log.Printf("Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("手机号脱敏:\n  原始: %s\n  脱敏: %s\n\n", testPhone, *result.Result)
	}

	// Test DesensitizeIDNumber
	testID := "110101199001011234"
	result, err = pluginInterface.Execute("DesensitizeIDNumber", map[string]string{"data": testID})
	if err != nil {
		log.Printf("Error desensitizing ID: %v", err)
	} else if !result.Success {
		log.Printf("Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("身份证号脱敏:\n  原始: %s\n  脱敏: %s\n\n", testID, *result.Result)
	}

	// Test DesensitizeEmail
	testEmail := "user@example.com"
	result, err = pluginInterface.Execute("DesensitizeEmail", map[string]string{"data": testEmail})
	if err != nil {
		log.Printf("Error desensitizing email: %v", err)
	} else if !result.Success {
		log.Printf("Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("邮箱脱敏:\n  原始: %s\n  脱敏: %s\n\n", testEmail, *result.Result)
	}

	// Test DesensitizeBankCard
	testCard := "6222021234567890123"
	result, err = pluginInterface.Execute("DesensitizeBankCard", map[string]string{"data": testCard})
	if err != nil {
		log.Printf("Error desensitizing bank card: %v", err)
	} else if !result.Success {
		log.Printf("Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("银行卡号脱敏:\n  原始: %s\n  脱敏: %s\n\n", testCard, *result.Result)
	}

	// Test DesensitizeAddress
	testAddress := "北京市朝阳区某某街道123号"
	result, err = pluginInterface.Execute("DesensitizeAddress", map[string]string{"data": testAddress})
	if err != nil {
		log.Printf("Error desensitizing address: %v", err)
	} else if !result.Success {
		log.Printf("Desensitize failed: %s", *result.Error)
	} else {
		fmt.Printf("地址脱敏:\n  原始: %s\n  脱敏: %s\n\n", testAddress, *result.Result)
	}

	fmt.Println("=== 测试完成 ===")

	// Check if plugin is still running
	if !client.Exited() {
		fmt.Println("\n插件进程运行正常")
	} else {
		fmt.Println("\n插件进程已退出")
	}
}
