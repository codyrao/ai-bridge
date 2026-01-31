package bridge

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"ai-bridge/pkg/types"

	"github.com/cloudwego/eino/schema"
)

// TestSDKClientOptions 测试SDK客户端Option功能
func TestSDKClientOptions(t *testing.T) {
	// 创建SDK配置
	config := &SDKConfig{
		Ollama: ProviderConfig{
			BaseURL:     "http://localhost:11434",
			Temperature: 0.7,
			MaxTokens:   256,
		},
	}

	sdk := NewSDK(config)

	// 创建高级客户端
	client, err := sdk.CreateSDKClient(types.ProviderOllama, "qwen3-coder:30b")
	if err != nil {
		t.Skipf("Ollama未启动，跳过测试: %v", err)
		return
	}

	fmt.Println("\n【SDK客户端Option测试】")
	fmt.Println()

	// 测试1: 简单生成（默认流式）
	t.Run("简单生成", func(t *testing.T) {
		fmt.Println("测试1: 简单生成（默认流式）")
		resp, err := client.Generate(context.Background(), "Hi")
		if err != nil {
			t.Errorf("生成失败: %v", err)
			return
		}
		fmt.Printf("  ✓ 响应: %s\n\n", resp)
	})

	// 测试2: 带历史记录的生成
	t.Run("带历史记录", func(t *testing.T) {
		fmt.Println("测试2: 带历史记录的生成")
		history := []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("你好！有什么可以帮助你的？", nil),
		}
		resp, err := client.Generate(context.Background(),
			"今天天气如何？",
			WithHistory(history),
		)
		if err != nil {
			t.Errorf("生成失败: %v", err)
			return
		}
		fmt.Printf("  ✓ 响应: %s\n\n", resp)
	})

	// 测试3: 禁用流式
	t.Run("禁用流式", func(t *testing.T) {
		fmt.Println("测试3: 禁用流式")
		resp, err := client.Generate(context.Background(),
			"Hello",
			WithStream(false),
		)
		if err != nil {
			t.Errorf("生成失败: %v", err)
			return
		}
		fmt.Printf("  ✓ 响应: %s\n\n", resp)
	})

	// 测试4: 设置超时
	t.Run("设置超时", func(t *testing.T) {
		fmt.Println("测试4: 设置超时（5秒）")
		resp, err := client.Generate(context.Background(),
			"Hello",
			WithTimeout(5*time.Second),
		)
		if err != nil {
			t.Errorf("生成失败: %v", err)
			return
		}
		fmt.Printf("  ✓ 响应: %s\n\n", resp)
	})

	// 测试5: 组合Option
	t.Run("组合Option", func(t *testing.T) {
		fmt.Println("测试5: 组合Option（历史+禁用流式+超时）")
		history := []*schema.Message{
			schema.SystemMessage("你是一个简洁的助手。"),
		}
		resp, err := client.Generate(context.Background(),
			"Hello",
			WithHistory(history),
			WithStream(false),
			WithTimeout(10*time.Second),
		)
		if err != nil {
			t.Errorf("生成失败: %v", err)
			return
		}
		fmt.Printf("  ✓ 响应: %s\n\n", resp)
	})

	// 测试6: 流式生成
	t.Run("流式生成", func(t *testing.T) {
		fmt.Println("测试6: 流式生成")
		stream, err := client.GenerateStream(context.Background(), "Hi")
		if err != nil {
			t.Errorf("创建流失败: %v", err)
			return
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
				t.Errorf("接收流数据失败: %v", err)
				return
			}
			chunkCount++
			fullContent += msg.Content
		}

		fmt.Printf("  ✓ 接收 %d 个数据块\n", chunkCount)
		fmt.Printf("  ✓ 完整内容: %s\n\n", fullContent)
	})

	// 测试7: Chat方法（默认流式）
	t.Run("Chat默认流式", func(t *testing.T) {
		fmt.Println("测试7: Chat方法（默认流式）")
		messages := []*schema.Message{
			schema.UserMessage("Hello"),
		}
		result, err := client.Chat(context.Background(), messages)
		if err != nil {
			t.Errorf("对话失败: %v", err)
			return
		}
		fmt.Printf("  ✓ 响应: %s\n", result.Content)
		fmt.Printf("  ✓ 是否流式: %v\n\n", result.Stream)
	})

	// 测试8: Chat方法（禁用流式）
	t.Run("Chat禁用流式", func(t *testing.T) {
		fmt.Println("测试8: Chat方法（禁用流式）")
		messages := []*schema.Message{
			schema.UserMessage("Hello"),
		}
		result, err := client.Chat(context.Background(), messages, WithStream(false))
		if err != nil {
			t.Errorf("对话失败: %v", err)
			return
		}
		fmt.Printf("  ✓ 响应: %s\n", result.Content)
		fmt.Printf("  ✓ 是否流式: %v\n\n", result.Stream)
	})

	// 测试9: ChatStream方法
	t.Run("ChatStream", func(t *testing.T) {
		fmt.Println("测试9: ChatStream方法")
		messages := []*schema.Message{
			schema.UserMessage("Hello"),
		}
		stream, err := client.ChatStream(context.Background(), messages)
		if err != nil {
			t.Errorf("创建流失败: %v", err)
			return
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
				t.Errorf("接收流数据失败: %v", err)
				return
			}
			chunkCount++
			fullContent += msg.Content
		}

		fmt.Printf("  ✓ 接收 %d 个数据块\n", chunkCount)
		fmt.Printf("  ✓ 完整内容: %s\n\n", fullContent)
	})
}

// TestClientConfigDefaults 测试客户端配置默认值
func TestClientConfigDefaults(t *testing.T) {
	cfg := DefaultClientConfig()

	if cfg.Stream != true {
		t.Error("默认Stream应该为true")
	}
	if cfg.Timeout != 60*time.Second {
		t.Errorf("默认Timeout应该为60s，实际为%v", cfg.Timeout)
	}
	if cfg.History != nil {
		t.Error("默认History应该为nil")
	}

	fmt.Println("✓ 客户端配置默认值测试通过")
}

// TestClientOptions 测试ClientOption功能
func TestClientOptions(t *testing.T) {
	// 测试WithHistory
	history := []*schema.Message{
		schema.UserMessage("test"),
	}
	cfg := DefaultClientConfig()
	WithHistory(history)(cfg)
	if len(cfg.History) != 1 {
		t.Error("WithHistory失败")
	}

	// 测试WithStream
	cfg = DefaultClientConfig()
	WithStream(false)(cfg)
	if cfg.Stream != false {
		t.Error("WithStream失败")
	}

	// 测试WithTimeout
	cfg = DefaultClientConfig()
	WithTimeout(30 * time.Second)(cfg)
	if cfg.Timeout != 30*time.Second {
		t.Error("WithTimeout失败")
	}

	// 测试多个Option组合
	cfg = DefaultClientConfig()
	WithHistory(history)(cfg)
	WithStream(false)(cfg)
	WithTimeout(30 * time.Second)(cfg)

	if len(cfg.History) != 1 || cfg.Stream != false || cfg.Timeout != 30*time.Second {
		t.Error("多个Option组合失败")
	}

	fmt.Println("✓ ClientOption功能测试通过")
}
