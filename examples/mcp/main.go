package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino/schema"

	"ai-bridge/pkg/bridge"
	"ai-bridge/pkg/mcp"
	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

func main() {
	ctx := context.Background()

	// 创建工具注册表
	registry := mcp.NewToolRegistry()

	// 注册示例工具
	for _, tool := range mcp.ExampleTools() {
		registry.Register(tool)
	}

	// 创建自定义工具
	customTool := mcp.NewTool(
		"get_time",
		"获取当前时间",
		mcp.CreateParameterSchema(
			map[string]interface{}{
				"timezone": mcp.CreateStringProperty("时区，例如: UTC, Asia/Shanghai"),
			},
			[]string{},
		),
		func(ctx context.Context, params map[string]interface{}) (string, error) {
			timezone, _ := params["timezone"].(string)
			if timezone == "" {
				timezone = "UTC"
			}
			return fmt.Sprintf("当前时间 (%s): 2024-01-01 12:00:00", timezone), nil
		},
	)
	registry.Register(customTool)

	// 获取所有工具
	tools := registry.ToEinoTools()
	fmt.Printf("已注册 %d 个工具\n", len(tools))

	// 使用Deepseek模型并配置工具
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		fmt.Println("请设置 DEEPSEEK_API_KEY 环境变量")
		return
	}

	client, err := bridge.NewAIClient(
		types.ProviderDeepseek,
		"deepseek-chat",
		options.WithAPIKey(apiKey),
		options.WithTools(tools...),
		options.WithTemperature(0.7),
	)
	if err != nil {
		log.Printf("创建客户端失败: %v", err)
		return
	}

	// 执行带工具的对话
	messages := []*schema.Message{
		schema.SystemMessage("你是一个有用的助手，可以使用工具来帮助用户。"),
		schema.UserMessage("请帮我计算 15 * 23 等于多少？"),
	}

	response, err := client.Chat(ctx, messages)
	if err != nil {
		log.Printf("对话失败: %v", err)
		return
	}

	fmt.Printf("模型响应: %s\n", response.Content)

	// 流式对话示例
	fmt.Println("\n=== 流式对话示例 ===")
	streamExample(ctx, client)
}

func streamExample(ctx context.Context, client types.AIBridge) {
	messages := []*schema.Message{
		schema.UserMessage("请介绍一下Go语言的特点"),
	}

	stream, err := client.ChatStream(ctx, messages)
	if err != nil {
		log.Printf("流式对话失败: %v", err)
		return
	}
	defer stream.Close()

	fmt.Print("流式响应: ")
	for {
		chunk, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Print(chunk.Content)
	}
	fmt.Println()
}
