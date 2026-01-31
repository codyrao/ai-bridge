package main

import (
	"context"
	"fmt"
	"log"

	"ai-bridge/pkg/bridge"
	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

func main() {
	ctx := context.Background()

	// Ollama本地模型示例
	// 默认连接到 http://localhost:11434
	// 可以使用 WithBaseURL 指定其他地址

	fmt.Println("=== Ollama本地模型示例 ===")

	// 示例1: 使用默认地址
	fmt.Println("\n--- 示例1: 使用llama2模型 ---")
	ollamaExample1(ctx)

	// 示例2: 使用自定义地址
	fmt.Println("\n--- 示例2: 使用自定义Ollama地址 ---")
	ollamaExample2(ctx)

	// 示例3: 使用其他模型
	fmt.Println("\n--- 示例3: 使用mistral模型 ---")
	ollamaExample3(ctx)
}

func ollamaExample1(ctx context.Context) {
	// 创建Ollama客户端，使用默认地址 localhost:11434
	client, err := bridge.NewAIClient(
		types.ProviderOllama,
		"llama2",
		options.WithMaxTokens(512),
		options.WithTemperature(0.7),
	)
	if err != nil {
		log.Printf("创建Ollama客户端失败: %v", err)
		return
	}

	// 生成文本
	response, err := client.Generate(ctx, "请用中文介绍一下自己")
	if err != nil {
		log.Printf("生成失败: %v", err)
		return
	}

	fmt.Printf("Ollama (llama2) 响应:\n%s\n", response)
}

func ollamaExample2(ctx context.Context) {
	// 创建Ollama客户端，使用自定义地址
	client, err := bridge.NewAIClient(
		types.ProviderOllama,
		"llama2",
		options.WithBaseURL("http://192.168.1.100:11434"), // 自定义Ollama服务器地址
		options.WithMaxTokens(256),
		options.WithTemperature(0.5),
	)
	if err != nil {
		log.Printf("创建Ollama客户端失败: %v", err)
		return
	}

	fmt.Printf("已创建连接到 http://192.168.1.100:11434 的Ollama客户端\n")
	fmt.Printf("模型: %s\n", client.GetModelInfo().Name)
}

func ollamaExample3(ctx context.Context) {
	// 使用其他Ollama模型
	client, err := bridge.NewAIClient(
		types.ProviderOllama,
		"mistral",
		options.WithMaxTokens(1024),
		options.WithTemperature(0.8),
	)
	if err != nil {
		log.Printf("创建Ollama客户端失败: %v", err)
		return
	}

	// 生成文本
	response, err := client.Generate(ctx, "What is the capital of France?")
	if err != nil {
		log.Printf("生成失败: %v", err)
		return
	}

	fmt.Printf("Ollama (mistral) 响应:\n%s\n", response)
}
