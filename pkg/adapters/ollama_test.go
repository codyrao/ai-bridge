package adapters

import (
	"context"
	"fmt"
	"testing"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"

	"github.com/cloudwego/eino/schema"
)

// TestOllamaQwen3Coder30b 测试本地Ollama的qwen3-coder:30b模型
func TestOllamaQwen3Coder30b(t *testing.T) {
	// 创建Ollama适配器
	// 默认使用 http://localhost:11434 作为Ollama服务地址
	adapter, err := NewOllamaAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithTemperature(0.7),
		options.WithMaxTokens(4096),
	)
	if err != nil {
		t.Fatalf("创建Ollama适配器失败: %v", err)
	}

	fmt.Printf("✓ 适配器创建成功\n")
	fmt.Printf("  厂商: %s\n", adapter.GetModelInfo().Provider)
	fmt.Printf("  模型: %s\n", adapter.GetModelInfo().Name)
	fmt.Printf("  最大Token: %d\n\n", adapter.GetModelInfo().MaxTokens)

	// 测试1: 简单的代码生成
	t.Run("代码生成测试", func(t *testing.T) {
		prompt := "请用Go语言写一个快速排序算法，并添加详细的中文注释"

		fmt.Printf("测试1 - 代码生成:\n")
		fmt.Printf("提示词: %s\n", prompt)
		fmt.Printf("等待响应...\n\n")

		resp, err := adapter.Generate(context.Background(), prompt)
		if err != nil {
			t.Errorf("代码生成失败: %v", err)
			return
		}

		fmt.Printf("✓ 响应成功:\n%s\n\n", resp)
	})

	// 测试2: 代码解释
	t.Run("代码解释测试", func(t *testing.T) {
		code := `
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}
`
		prompt := fmt.Sprintf("请解释以下Go代码的功能和原理:\n%s", code)

		fmt.Printf("测试2 - 代码解释:\n")
		fmt.Printf("提示词: %s\n", prompt)
		fmt.Printf("等待响应...\n\n")

		resp, err := adapter.Generate(context.Background(), prompt)
		if err != nil {
			t.Errorf("代码解释失败: %v", err)
			return
		}

		fmt.Printf("✓ 响应成功:\n%s\n\n", resp)
	})

	// 测试3: 多轮对话
	t.Run("多轮对话测试", func(t *testing.T) {
		fmt.Printf("测试3 - 多轮对话:\n")

		messages := []*schema.Message{
			schema.SystemMessage("你是一个专业的Go语言编程助手，擅长代码优化和最佳实践建议。"),
			schema.UserMessage("Go语言中如何优雅地处理错误？"),
		}

		fmt.Printf("系统: %s\n", messages[0].Content)
		fmt.Printf("用户: %s\n", messages[1].Content)
		fmt.Printf("等待响应...\n\n")

		resp, err := adapter.Chat(context.Background(), messages)
		if err != nil {
			t.Errorf("多轮对话失败: %v", err)
			return
		}

		fmt.Printf("助手: %s\n\n", resp.Content)

		// 第二轮对话
		messages = append(messages, resp)
		messages = append(messages, schema.UserMessage("能给我一个具体的代码示例吗？"))

		fmt.Printf("用户: %s\n", messages[len(messages)-1].Content)
		fmt.Printf("等待响应...\n\n")

		resp2, err := adapter.Chat(context.Background(), messages)
		if err != nil {
			t.Errorf("第二轮对话失败: %v", err)
			return
		}

		fmt.Printf("助手: %s\n\n", resp2.Content)
	})

	// 测试4: 流式响应
	t.Run("流式响应测试", func(t *testing.T) {
		prompt := "请写一个Go语言的HTTP服务器示例，包含路由和中间件"

		fmt.Printf("测试4 - 流式响应:\n")
		fmt.Printf("提示词: %s\n", prompt)
		fmt.Printf("等待流式响应...\n\n")

		result, err := adapter.GenerateStream(context.Background(), prompt)
		if err != nil {
			t.Errorf("流式响应失败: %v", err)
			return
		}

		fmt.Printf("✓ 流式响应完成，总长度: %d 字符\n\n", len(result))
	})

	// 测试5: 代码重构建议
	t.Run("代码重构测试", func(t *testing.T) {
		code := `
package main

import "fmt"

func main() {
    for i := 0; i < 100; i++ {
        if i % 2 == 0 {
            fmt.Println(i)
        }
    }
}
`
		prompt := fmt.Sprintf("请重构以下代码，使其更简洁、更高效:\n%s", code)

		fmt.Printf("测试5 - 代码重构:\n")
		fmt.Printf("提示词: %s\n", prompt)
		fmt.Printf("等待响应...\n\n")

		resp, err := adapter.Generate(context.Background(), prompt)
		if err != nil {
			t.Errorf("代码重构建议失败: %v", err)
			return
		}

		fmt.Printf("✓ 响应成功:\n%s\n\n", resp)
	})
}

// TestOllamaConnection 测试Ollama连接
func TestOllamaConnection(t *testing.T) {
	adapter, err := NewOllamaAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
	)
	if err != nil {
		t.Fatalf("连接Ollama失败，请确保:\n1. Ollama服务已启动 (ollama serve)\n2. qwen3-coder-30b模型已下载 (ollama pull qwen3-coder-30b)\n3. 服务运行在 http://localhost:11434\n\n错误: %v", err)
	}

	modelInfo := adapter.GetModelInfo()
	fmt.Printf("✓ Ollama连接成功\n")
	fmt.Printf("  厂商: %s\n", modelInfo.Provider)
	fmt.Printf("  模型: %s\n", modelInfo.Name)
	fmt.Printf("  描述: %s\n", modelInfo.Description)
	fmt.Printf("  最大Token: %d\n", modelInfo.MaxTokens)

	// 简单ping测试
	resp, err := adapter.Generate(context.Background(), "Hello")
	if err != nil {
		t.Errorf("简单对话测试失败: %v", err)
		return
	}

	fmt.Printf("✓ 对话测试成功\n")
	fmt.Printf("  响应: %s\n", resp)
}
