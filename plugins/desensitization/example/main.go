package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/desensitization"
)

func main() {
	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: desensitization.Handshake,
		Plugins:         desensitization.PluginMap,
		Cmd:             exec.Command("../../../host-server/bin/plugins/desensitization"),
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

	// Request the plugin
	raw, err := rpcClient.Dispense("desensitizer")
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Desensitizer now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	desensitizer := raw.(desensitization.Desensitizer)

	fmt.Println("=== 数据脱敏插件测试 ===\n")

	// Test DesensitizeName
	testName := "张三"
	result, err := desensitizer.DesensitizeName(testName)
	if err != nil {
		log.Printf("Error desensitizing name: %v", err)
	} else {
		fmt.Printf("姓名脱敏:\n  原始: %s\n  脱敏: %s\n\n", testName, result)
	}

	// Test DesensitizeTelNo
	testPhone := "13812345678"
	result, err = desensitizer.DesensitizeTelNo(testPhone)
	if err != nil {
		log.Printf("Error desensitizing phone: %v", err)
	} else {
		fmt.Printf("手机号脱敏:\n  原始: %s\n  脱敏: %s\n\n", testPhone, result)
	}

	// Test DesensitizeIDNumber
	testID := "110101199001011234"
	result, err = desensitizer.DesensitizeIDNumber(testID)
	if err != nil {
		log.Printf("Error desensitizing ID: %v", err)
	} else {
		fmt.Printf("身份证号脱敏:\n  原始: %s\n  脱敏: %s\n\n", testID, result)
	}

	// Test DesensitizeEmail
	testEmail := "user@example.com"
	result, err = desensitizer.DesensitizeEmail(testEmail)
	if err != nil {
		log.Printf("Error desensitizing email: %v", err)
	} else {
		fmt.Printf("邮箱脱敏:\n  原始: %s\n  脱敏: %s\n\n", testEmail, result)
	}

	// Test DesensitizeBankCard
	testCard := "6222021234567890123"
	result, err = desensitizer.DesensitizeBankCard(testCard)
	if err != nil {
		log.Printf("Error desensitizing bank card: %v", err)
	} else {
		fmt.Printf("银行卡号脱敏:\n  原始: %s\n  脱敏: %s\n\n", testCard, result)
	}

	// Test DesensitizeAddress
	testAddress := "北京市朝阳区某某街道123号"
	result, err = desensitizer.DesensitizeAddress(testAddress)
	if err != nil {
		log.Printf("Error desensitizing address: %v", err)
	} else {
		fmt.Printf("地址脱敏:\n  原始: %s\n  脱敏: %s\n\n", testAddress, result)
	}

	fmt.Println("=== 测试完成 ===")

	// Check if plugin is still running
	if !client.Exited() {
		fmt.Println("\n插件进程运行正常")
	} else {
		exitErr := client.ReattachConfig()
		fmt.Printf("\n插件进程已退出: %v\n", exitErr)
		os.Exit(1)
	}
}
