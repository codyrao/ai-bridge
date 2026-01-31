package adapters

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"

	"github.com/cloudwego/eino/schema"
)

// TestOllamaStream 测试Ollama流式返回
func TestOllamaStream(t *testing.T) {
	adapter, err := NewOllamaAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithTemperature(0.7),
		options.WithMaxTokens(512),
		options.WithTimeout(120*time.Second),
	)
	if err != nil {
		t.Fatalf("创建适配器失败: %v", err)
	}

	// 测试ChatStream
	t.Run("ChatStream", func(t *testing.T) {
		messages := []*schema.Message{
			schema.UserMessage("用一句话介绍Go语言"),
		}

		fmt.Printf("测试ChatStream:\n")
		fmt.Printf("等待流式响应...\n")

		stream, err := adapter.ChatStream(context.Background(), messages)
		if err != nil {
			t.Fatalf("ChatStream失败: %v", err)
		}
		defer stream.Close()

		var fullContent string
		chunkCount := 0

		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("接收流数据失败: %v", err)
			}
			chunkCount++
			fullContent += msg.Content
			fmt.Printf("[%d] %s\n", chunkCount, msg.Content)
		}

		fmt.Printf("\n✓ 流式响应完成\n")
		fmt.Printf("  总块数: %d\n", chunkCount)
		fmt.Printf("  总长度: %d 字符\n", len(fullContent))
		fmt.Printf("  完整内容: %s\n\n", fullContent)
	})

	// 测试GenerateStream
	t.Run("GenerateStream", func(t *testing.T) {
		prompt := "列举Go语言的3个主要优点"

		fmt.Printf("测试GenerateStream:\n")
		fmt.Printf("提示词: %s\n", prompt)
		fmt.Printf("等待流式响应...\n")

		result, err := adapter.GenerateStream(context.Background(), prompt)
		if err != nil {
			t.Fatalf("GenerateStream失败: %v", err)
		}

		fmt.Printf("\n✓ 流式响应完成\n")
		fmt.Printf("  总长度: %d 字符\n", len(result))
		fmt.Printf("  完整内容: %s\n\n", result)
	})
}

// TestStreamWithDifferentModels 测试不同模型的流式返回
func TestStreamWithDifferentModels(t *testing.T) {
	// 注意：这些测试需要相应的API密钥才能运行
	// 这里主要测试配置项是否正确传递

	testCases := []struct {
		name     string
		provider types.Provider
		model    string
		skip     bool
	}{
		{
			name:     "Ollama_qwen3-coder",
			provider: types.ProviderOllama,
			model:    "qwen3-coder:30b",
			skip:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip("跳过此测试")
			}

			var adapter types.AIBridge
			var err error

			switch tc.provider {
			case types.ProviderOllama:
				adapter, err = NewOllamaAdapter(
					tc.provider,
					tc.model,
					options.WithBaseURL("http://localhost:11434"),
					options.WithTemperature(0.7),
					options.WithMaxTokens(256),
				)
			default:
				t.Skip("需要API密钥，跳过")
				return
			}

			if err != nil {
				t.Fatalf("创建适配器失败: %v", err)
			}

			// 测试流式返回
			messages := []*schema.Message{
				schema.UserMessage("Hi"),
			}

			stream, err := adapter.ChatStream(context.Background(), messages)
			if err != nil {
				t.Fatalf("ChatStream失败: %v", err)
			}
			defer stream.Close()

			chunkCount := 0
			for {
				_, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("接收流数据失败: %v", err)
				}
				chunkCount++
			}

			fmt.Printf("✓ %s 流式返回正常，接收 %d 个数据块\n", tc.name, chunkCount)
		})
	}
}

// TestStreamConfigOptions 测试流式配置选项
func TestStreamConfigOptions(t *testing.T) {
	// 测试各种配置选项是否正确应用到流式请求
	adapter, err := NewOllamaAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithTemperature(0.5),
		options.WithTopP(0.8),
		options.WithMaxTokens(100),
		options.WithTimeout(60*time.Second),
	)
	if err != nil {
		t.Fatalf("创建适配器失败: %v", err)
	}

	info := adapter.GetModelInfo()
	fmt.Printf("流式配置测试:\n")
	fmt.Printf("  模型: %s\n", info.Name)
	fmt.Printf("  厂商: %s\n", info.Provider)

	// 执行流式请求验证配置
	messages := []*schema.Message{
		schema.UserMessage("测试流式配置"),
	}

	stream, err := adapter.ChatStream(context.Background(), messages)
	if err != nil {
		t.Fatalf("ChatStream失败: %v", err)
	}
	defer stream.Close()

	// 简单验证流式返回正常工作
	_, err = stream.Recv()
	if err != nil && err != io.EOF {
		t.Fatalf("接收流数据失败: %v", err)
	}

	fmt.Printf("✓ 流式配置选项测试通过\n")
}
