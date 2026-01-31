package types

import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// Config AI客户端配置
type Config struct {
	// APIKey API密钥
	APIKey string

	// BaseURL 基础URL（可选，用于自定义端点）
	BaseURL string

	// Timeout 请求超时时间
	Timeout time.Duration

	// MaxRetries 最大重试次数
	MaxRetries int

	// Temperature 温度参数
	Temperature float32

	// TopP 核采样参数
	TopP float32

	// MaxTokens 最大生成token数
	MaxTokens int

	// Tools MCP工具列表
	Tools []tool.BaseTool

	// Stream 是否使用流式响应
	Stream bool

	// ExtraHeaders 额外请求头
	ExtraHeaders map[string]string

	// Organization OpenAI组织ID（仅OpenAI使用）
	Organization string

	// Proxy 代理地址
	Proxy string

	// EnableLog 是否启用日志
	EnableLog bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Timeout:      120 * time.Second,
		MaxRetries:   3,
		Temperature:  0.7,
		TopP:         0.9,
		MaxTokens:    2048,
		Stream:       false,
		ExtraHeaders: make(map[string]string),
	}
}

// AIBridge AI聚合器接口
type AIBridge interface {
	// Chat 执行对话（非流式）
	Chat(ctx context.Context, messages []*schema.Message) (*schema.Message, error)

	// ChatStream 执行对话（流式）
	ChatStream(ctx context.Context, messages []*schema.Message) (*schema.StreamReader[*schema.Message], error)

	// Generate 生成文本（简化接口）
	Generate(ctx context.Context, prompt string) (string, error)

	// GenerateStream 生成文本（流式）
	GenerateStream(ctx context.Context, prompt string) (string, error)

	// GetModelInfo 获取当前模型信息
	GetModelInfo() *ModelInfo
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content      string `json:"content"`
	Role         string `json:"role"`
	FinishReason string `json:"finish_reason"`
	Usage        Usage  `json:"usage"`
}

// Usage token使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamChunk 流式响应块
type StreamChunk struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
