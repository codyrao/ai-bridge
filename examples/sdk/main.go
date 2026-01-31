package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"ai-bridge/pkg/bridge"
	"ai-bridge/pkg/types"

	"github.com/cloudwego/eino/schema"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("     AI Bridge SDK 使用示例")
	fmt.Println("========================================")
	fmt.Println()

	// 方式1：使用Config结构配置SDK
	demoWithConfig()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println()

	// 方式2：从环境变量加载配置
	demoWithEnvConfig()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println()

	// 方式3：使用高级客户端Option
	demoWithAdvancedOptions()
}

// demoWithConfig 使用Config结构配置
func demoWithConfig() {
	fmt.Println("【示例1：使用Config结构配置】")
	fmt.Println()

	// 创建SDK配置
	config := &bridge.SDKConfig{
		GPT: bridge.ProviderConfig{
			APIKey:      os.Getenv("GPT_API_KEY"),
			Temperature: 0.7,
			MaxTokens:   2048,
		},
		QWen: bridge.ProviderConfig{
			APIKey:      os.Getenv("QWEN_API_KEY"),
			Temperature: 0.8,
			MaxTokens:   1024,
		},
		Deepseek: bridge.ProviderConfig{
			APIKey:    os.Getenv("DEEPSEEK_API_KEY"),
			MaxTokens: 4096,
		},
		Ollama: bridge.ProviderConfig{
			BaseURL:     "http://localhost:11434",
			Temperature: 0.7,
			MaxTokens:   512,
		},
	}

	// 创建SDK实例
	sdk := bridge.NewSDK(config)

	// 测试Ollama本地模型（如果有的话）
	ollamaClient, err := sdk.CreateClient(types.ProviderOllama, "qwen3-coder:30b")
	if err != nil {
		log.Printf("创建Ollama客户端失败: %v\n", err)
	} else {
		fmt.Println("✓ Ollama客户端创建成功")
		resp, err := ollamaClient.Generate(context.Background(), "你好，请简短介绍自己")
		if err != nil {
			log.Printf("Ollama调用失败: %v\n", err)
		} else {
			fmt.Printf("Ollama响应: %s\n\n", resp)
		}
	}

	// 测试其他模型（需要API Key）
	testModels := []struct {
		name     string
		provider types.Provider
		model    string
	}{
		{"GPT-3.5", types.ProviderGPT, "gpt-3.5-turbo"},
		{"通义千问", types.ProviderQWen, "qwen-turbo"},
		{"Deepseek", types.ProviderDeepseek, "deepseek-chat"},
	}

	for _, tm := range testModels {
		client, err := sdk.CreateClient(tm.provider, tm.model)
		if err != nil {
			fmt.Printf("✗ %s: %v\n", tm.name, err)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		resp, err := client.Generate(ctx, "你好，请用一句话介绍自己")
		cancel()

		if err != nil {
			fmt.Printf("✗ %s调用失败: %v\n", tm.name, err)
			continue
		}

		fmt.Printf("✓ %s: %s\n", tm.name, resp)
	}
}

// demoWithEnvConfig 从环境变量加载配置
func demoWithEnvConfig() {
	fmt.Println("【示例2：从环境变量加载配置】")
	fmt.Println()

	// 从环境变量自动加载配置
	config := bridge.ConfigFromEnv()

	// 显示加载的配置
	fmt.Println("已加载的配置:")
	if config.GPT.APIKey != "" {
		fmt.Println("  ✓ GPT: 已配置")
	} else {
		fmt.Println("  ✗ GPT: 未配置 (GPT_API_KEY)")
	}

	if config.QWen.APIKey != "" {
		fmt.Println("  ✓ QWen: 已配置")
	} else {
		fmt.Println("  ✗ QWen: 未配置 (QWEN_API_KEY)")
	}

	if config.Deepseek.APIKey != "" {
		fmt.Println("  ✓ Deepseek: 已配置")
	} else {
		fmt.Println("  ✗ Deepseek: 未配置 (DEEPSEEK_API_KEY)")
	}

	fmt.Printf("  ✓ Ollama: %s\n", config.Ollama.BaseURL)

	// 创建SDK实例
	sdk := bridge.NewSDK(config)

	// 获取支持的厂商列表
	fmt.Println()
	fmt.Println("支持的厂商列表:")
	providers := bridge.GetProviders()
	for _, p := range providers {
		models := bridge.GetModels(p)
		fmt.Printf("  - %s: %d个模型\n", p, len(models))
	}

	// 尝试使用Ollama
	fmt.Println()
	fmt.Println("测试Ollama连接:")
	client, err := sdk.CreateClient(types.ProviderOllama, "qwen3-coder:30b")
	if err != nil {
		fmt.Printf("  ✗ 失败: %v\n", err)
		return
	}

	info := client.GetModelInfo()
	fmt.Printf("  ✓ 连接成功\n")
	fmt.Printf("    模型: %s\n", info.Name)
	fmt.Printf("    厂商: %s\n", info.Provider)
	fmt.Printf("    描述: %s\n", info.Description)
}

// demoWithAdvancedOptions 使用高级客户端Option
func demoWithAdvancedOptions() {
	fmt.Println("【示例3：使用高级客户端Option】")
	fmt.Println()

	// 创建SDK配置
	config := &bridge.SDKConfig{
		Ollama: bridge.ProviderConfig{
			BaseURL:     "http://localhost:11434",
			Temperature: 0.7,
			MaxTokens:   256,
		},
	}

	sdk := bridge.NewSDK(config)

	// 创建高级客户端
	client, err := sdk.CreateSDKClient(types.ProviderOllama, "qwen3-coder:30b")
	if err != nil {
		log.Printf("创建高级客户端失败: %v\n", err)
		return
	}

	// 示例1: 简单生成（默认流式）
	fmt.Println("1. 简单生成（默认流式）:")
	resp, err := client.Generate(context.Background(), "Hello")
	if err != nil {
		log.Printf("  失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n\n", resp)
	}

	// 示例2: 带历史记录的生成
	fmt.Println("2. 带历史记录的生成:")
	history := []*schema.Message{
		schema.UserMessage("你好"),
		schema.AssistantMessage("你好！有什么可以帮助你的？", nil),
	}
	resp, err = client.Generate(context.Background(),
		"今天天气如何？",
		bridge.WithHistory(history),
	)
	if err != nil {
		log.Printf("  失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n\n", resp)
	}

	// 示例3: 禁用流式
	fmt.Println("3. 禁用流式:")
	resp, err = client.Generate(context.Background(),
		"Hello",
		bridge.WithStream(false),
	)
	if err != nil {
		log.Printf("  失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n\n", resp)
	}

	// 示例4: 设置超时
	fmt.Println("4. 设置超时（10秒）:")
	resp, err = client.Generate(context.Background(),
		"Hello",
		bridge.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Printf("  失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n\n", resp)
	}

	// 示例5: 流式生成（实时输出）
	fmt.Println("5. 流式生成（实时输出）:")
	stream, err := client.GenerateStream(context.Background(), "Hello")
	if err != nil {
		log.Printf("  失败: %v\n", err)
	} else {
		defer stream.Close()
		fmt.Print("  响应: ")
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("  接收错误: %v\n", err)
				break
			}
			fmt.Print(msg.Content)
		}
		fmt.Println()
		fmt.Println()
	}

	// 示例6: Chat方法（默认流式）
	fmt.Println("6. Chat方法（默认流式）:")
	messages := []*schema.Message{
		schema.UserMessage("Hello"),
	}
	result, err := client.Chat(context.Background(), messages)
	if err != nil {
		log.Printf("  失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n", result.Content)
		fmt.Printf("  是否流式: %v\n\n", result.Stream)
	}

	// 示例7: Chat方法（禁用流式）
	fmt.Println("7. Chat方法（禁用流式）:")
	result, err = client.Chat(context.Background(), messages, bridge.WithStream(false))
	if err != nil {
		log.Printf("  失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n", result.Content)
		fmt.Printf("  是否流式: %v\n\n", result.Stream)
	}
}
