package adapters

import (
	"context"
	"fmt"
	"testing"
	"time"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
	"github.com/cloudwego/eino/schema"
)

// TestOllamaQuickTest 快速测试Ollama qwen3-coder:30b模型
func TestOllamaQuickTest(t *testing.T) {
	// 创建Ollama适配器
	adapter, err := NewOllamaAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithTemperature(0.7),
		options.WithMaxTokens(1024),
		options.WithTimeout(180*time.Second),
	)
	if err != nil {
		t.Fatalf("创建Ollama适配器失败: %v", err)
	}

	fmt.Printf("✓ 适配器创建成功\n")
	fmt.Printf("  厂商: %s\n", adapter.GetModelInfo().Provider)
	fmt.Printf("  模型: %s\n", adapter.GetModelInfo().Name)
	fmt.Printf("  最大Token: %d\n\n", adapter.GetModelInfo().MaxTokens)

	// 测试1: 简单问候
	t.Run("简单问候", func(t *testing.T) {
		prompt := "你好，请简短介绍一下自己"
		
		fmt.Printf("测试1 - 简单问候:\n")
		fmt.Printf("提示词: %s\n", prompt)
		fmt.Printf("等待响应...\n")

		start := time.Now()
		resp, err := adapter.Generate(context.Background(), prompt)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("失败: %v", err)
			return
		}

		fmt.Printf("✓ 响应成功 (耗时: %.2f秒):\n%s\n\n", elapsed.Seconds(), resp)
	})

	// 测试2: 简单代码问题
	t.Run("简单代码问题", func(t *testing.T) {
		prompt := "Go语言中map和slice的区别是什么？请简要回答。"
		
		fmt.Printf("测试2 - 简单代码问题:\n")
		fmt.Printf("提示词: %s\n", prompt)
		fmt.Printf("等待响应...\n")

		start := time.Now()
		resp, err := adapter.Generate(context.Background(), prompt)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("失败: %v", err)
			return
		}

		fmt.Printf("✓ 响应成功 (耗时: %.2f秒):\n%s\n\n", elapsed.Seconds(), resp)
	})

	// 测试3: 代码补全
	t.Run("代码补全", func(t *testing.T) {
		code := `func add(a, b int) int {`
		prompt := fmt.Sprintf("请补全以下Go函数:\n%s", code)
		
		fmt.Printf("测试3 - 代码补全:\n")
		fmt.Printf("提示词: %s\n", prompt)
		fmt.Printf("等待响应...\n")

		start := time.Now()
		resp, err := adapter.Generate(context.Background(), prompt)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("失败: %v", err)
			return
		}

		fmt.Printf("✓ 响应成功 (耗时: %.2f秒):\n%s\n\n", elapsed.Seconds(), resp)
	})

	// 测试4: 多轮对话
	t.Run("多轮对话", func(t *testing.T) {
		fmt.Printf("测试4 - 多轮对话:\n")

		messages := []*schema.Message{
			schema.UserMessage("什么是Go语言中的goroutine？"),
		}

		fmt.Printf("用户: %s\n", messages[0].Content)
		fmt.Printf("等待响应...\n")

		start := time.Now()
		resp, err := adapter.Chat(context.Background(), messages)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("失败: %v", err)
			return
		}

		fmt.Printf("助手 (耗时: %.2f秒): %s\n\n", elapsed.Seconds(), resp.Content)

		// 第二轮
		messages = append(messages, resp)
		messages = append(messages, schema.UserMessage("它和线程有什么区别？"))

		fmt.Printf("用户: %s\n", messages[len(messages)-1].Content)
		fmt.Printf("等待响应...\n")

		start = time.Now()
		resp2, err := adapter.Chat(context.Background(), messages)
		elapsed = time.Since(start)

		if err != nil {
			t.Errorf("失败: %v", err)
			return
		}

		fmt.Printf("助手 (耗时: %.2f秒): %s\n\n", elapsed.Seconds(), resp2.Content)
	})
}

// TestOllamaModelInfo 查看模型信息
func TestOllamaModelInfo(t *testing.T) {
	adapter, err := NewOllamaAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
	)
	if err != nil {
		t.Fatalf("创建适配器失败: %v", err)
	}

	info := adapter.GetModelInfo()
	fmt.Printf("模型信息:\n")
	fmt.Printf("  名称: %s\n", info.Name)
	fmt.Printf("  厂商: %s\n", info.Provider)
	fmt.Printf("  描述: %s\n", info.Description)
	fmt.Printf("  最大Token: %d\n", info.MaxTokens)
}
