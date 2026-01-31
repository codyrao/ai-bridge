package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"ai-bridge/pkg/bridge"
	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

func main() {
	ctx := context.Background()

	// 示例1: 使用Deepseek模型
	fmt.Println("=== 示例1: Deepseek ===")
	deepseekExample(ctx)

	// 示例2: 使用OpenAI GPT模型
	fmt.Println("\n=== 示例2: GPT ===")
	gptExample(ctx)

	// 示例3: 使用通义千问
	fmt.Println("\n=== 示例3: QWen ===")
	qwenExample(ctx)

	// 示例4: 列出所有支持的模型
	fmt.Println("\n=== 示例4: 支持的模型列表 ===")
	listAllModels()
}

func deepseekExample(ctx context.Context) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		fmt.Println("请设置 DEEPSEEK_API_KEY 环境变量")
		return
	}

	// 创建Deepseek客户端，使用多种Option配置
	client, err := bridge.NewAIClient(
		types.ProviderDeepseek,
		"deepseek-chat",
		options.WithAPIKey(apiKey),
		options.WithTemperature(0.8),
		options.WithMaxTokens(1024),
		options.WithTimeout(60),
	)
	if err != nil {
		log.Printf("创建Deepseek客户端失败: %v", err)
		return
	}

	// 生成文本
	response, err := client.Generate(ctx, "请用一句话介绍Go语言")
	if err != nil {
		log.Printf("生成失败: %v", err)
		return
	}

	fmt.Printf("Deepseek响应: %s\n", response)
}

func gptExample(ctx context.Context) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("请设置 OPENAI_API_KEY 环境变量")
		return
	}

	// 创建GPT客户端
	client, err := bridge.NewAIClient(
		types.ProviderGPT,
		"gpt-3.5-turbo",
		options.WithAPIKey(apiKey),
		options.WithTemperature(0.7),
		options.WithMaxTokens(512),
	)
	if err != nil {
		log.Printf("创建GPT客户端失败: %v", err)
		return
	}

	// 生成文本
	response, err := client.Generate(ctx, "Hello, how are you?")
	if err != nil {
		log.Printf("生成失败: %v", err)
		return
	}

	fmt.Printf("GPT响应: %s\n", response)
}

func qwenExample(ctx context.Context) {
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		fmt.Println("请设置 QWEN_API_KEY 环境变量")
		return
	}

	// 创建通义千问客户端
	client, err := bridge.NewAIClient(
		types.ProviderQWen,
		"qwen-turbo",
		options.WithAPIKey(apiKey),
		options.WithTemperature(0.7),
		options.WithMaxTokens(1024),
	)
	if err != nil {
		log.Printf("创建QWen客户端失败: %v", err)
		return
	}

	// 生成文本
	response, err := client.Generate(ctx, "请介绍一下人工智能的发展历史")
	if err != nil {
		log.Printf("生成失败: %v", err)
		return
	}

	fmt.Printf("QWen响应: %s\n", response)
}

func listAllModels() {
	providers := bridge.GetProviders()
	for _, provider := range providers {
		models := bridge.GetModels(provider)
		fmt.Printf("\n厂商: %s\n", provider)
		for _, model := range models {
			fmt.Printf("  - %s: %s (MaxTokens: %d)\n", model.Name, model.Description, model.MaxTokens)
		}
	}
}
