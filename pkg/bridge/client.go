package bridge

import (
	"context"
	"io"
	"time"

	"github.com/cloudwego/eino/schema"

	"ai-bridge/pkg/types"
)

// ClientOption 客户端调用选项
type ClientOption func(*ClientConfig)

// ClientConfig 客户端调用配置
type ClientConfig struct {
	History      []*schema.Message // 对话历史
	Stream       bool              // 是否启用流式返回（默认true）
	Timeout      time.Duration     // 超时时间（默认60s）
	SystemPrompt string            // 系统提示词（可选，覆盖适配器配置）
}

// DefaultClientConfig 返回默认客户端配置
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		History: nil,
		Stream:  true,             // 默认启用流式
		Timeout: 60 * time.Second, // 默认60秒超时
	}
}

// WithHistory 设置对话历史
func WithHistory(history []*schema.Message) ClientOption {
	return func(c *ClientConfig) {
		c.History = history
	}
}

// WithStream 设置是否启用流式返回
func WithStream(stream bool) ClientOption {
	return func(c *ClientConfig) {
		c.Stream = stream
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *ClientConfig) {
		c.Timeout = timeout
	}
}

// WithSystemPrompt 设置系统提示词（可选，覆盖适配器配置）
func WithSystemPrompt(prompt string) ClientOption {
	return func(c *ClientConfig) {
		c.SystemPrompt = prompt
	}
}

// SDKClient SDK客户端包装器
type SDKClient struct {
	inner types.AIBridge
}

// NewSDKClient 创建SDK客户端
func NewSDKClient(inner types.AIBridge) *SDKClient {
	return &SDKClient{inner: inner}
}

// Generate 生成文本（支持Option）
// 支持选项：
//   - WithHistory(history): 设置对话历史
//   - WithStream(bool): 是否启用流式（默认true）
//   - WithTimeout(duration): 设置超时时间（默认60s）
//   - WithSystemPrompt(prompt): 设置系统提示词（可选，支持{{question}}宏）
func (c *SDKClient) Generate(ctx context.Context, prompt string, opts ...ClientOption) (string, error) {
	cfg := DefaultClientConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// 设置超时
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	// 构建消息列表
	messages := make([]*schema.Message, 0)

	// 添加系统提示词（如果指定），支持模板渲染
	if cfg.SystemPrompt != "" {
		systemPrompt := renderSystemPromptTemplate(cfg.SystemPrompt, prompt)
		messages = append(messages, schema.SystemMessage(systemPrompt))
	}

	if len(cfg.History) > 0 {
		messages = append(messages, cfg.History...)
	}
	messages = append(messages, schema.UserMessage(prompt))

	// 根据Stream配置选择调用方式
	if cfg.Stream {
		stream, err := c.inner.ChatStream(ctx, messages)
		if err != nil {
			return "", err
		}
		defer stream.Close()

		var result string
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return result, err
			}
			result += msg.Content
		}
		return result, nil
	}

	// 非流式调用
	resp, err := c.inner.Chat(ctx, messages)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

// GenerateStream 生成文本（流式，支持Option）
// 支持选项：
//   - WithHistory(history): 设置对话历史
//   - WithTimeout(duration): 设置超时时间（默认60s）
//   - WithSystemPrompt(prompt): 设置系统提示词（可选，支持{{question}}宏）
func (c *SDKClient) GenerateStream(ctx context.Context, prompt string, opts ...ClientOption) (*StreamReader, error) {
	cfg := DefaultClientConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// 设置超时（流式方法不使用defer，由调用者控制）
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		_ = cancel // 流式读取完成后自动取消
	}

	// 构建消息列表
	messages := make([]*schema.Message, 0)

	// 添加系统提示词（如果指定），支持模板渲染
	if cfg.SystemPrompt != "" {
		systemPrompt := renderSystemPromptTemplate(cfg.SystemPrompt, prompt)
		messages = append(messages, schema.SystemMessage(systemPrompt))
	}

	if len(cfg.History) > 0 {
		messages = append(messages, cfg.History...)
	}
	messages = append(messages, schema.UserMessage(prompt))

	stream, err := c.inner.ChatStream(ctx, messages)
	if err != nil {
		return nil, err
	}

	return &StreamReader{inner: stream}, nil
}

// Chat 对话（支持Option）
// 支持选项：
//   - WithStream(bool): 是否启用流式（默认true）
//   - WithTimeout(duration): 设置超时时间（默认60s）
//   - WithSystemPrompt(prompt): 设置系统提示词（可选）
func (c *SDKClient) Chat(ctx context.Context, messages []*schema.Message, opts ...ClientOption) (*ChatResult, error) {
	cfg := DefaultClientConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// 设置超时
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	// 添加系统提示词（如果指定）
	if cfg.SystemPrompt != "" {
		// 检查是否已有系统消息
		hasSystem := false
		for _, msg := range messages {
			if msg.Role == schema.System {
				hasSystem = true
				break
			}
		}
		if !hasSystem {
			messages = append([]*schema.Message{schema.SystemMessage(cfg.SystemPrompt)}, messages...)
		}
	}

	// 根据Stream配置选择调用方式
	if cfg.Stream {
		stream, err := c.inner.ChatStream(ctx, messages)
		if err != nil {
			return nil, err
		}

		// 收集流式响应
		var fullContent string
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				stream.Close()
				return nil, err
			}
			fullContent += msg.Content
		}
		stream.Close()

		return &ChatResult{
			Content: fullContent,
			Stream:  true,
		}, nil
	}

	// 非流式调用
	resp, err := c.inner.Chat(ctx, messages)
	if err != nil {
		return nil, err
	}

	return &ChatResult{
		Content: resp.Content,
		Stream:  false,
	}, nil
}

// ChatStream 对话流式（支持Option）
// 支持选项：
//   - WithTimeout(duration): 设置超时时间（默认60s）
//   - WithSystemPrompt(prompt): 设置系统提示词（可选）
func (c *SDKClient) ChatStream(ctx context.Context, messages []*schema.Message, opts ...ClientOption) (*StreamReader, error) {
	cfg := DefaultClientConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// 设置超时（流式方法不使用defer，由调用者控制）
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		_ = cancel // 流式读取完成后自动取消
	}

	// 添加系统提示词（如果指定）
	if cfg.SystemPrompt != "" {
		// 检查是否已有系统消息
		hasSystem := false
		for _, msg := range messages {
			if msg.Role == schema.System {
				hasSystem = true
				break
			}
		}
		if !hasSystem {
			messages = append([]*schema.Message{schema.SystemMessage(cfg.SystemPrompt)}, messages...)
		}
	}

	stream, err := c.inner.ChatStream(ctx, messages)
	if err != nil {
		return nil, err
	}

	return &StreamReader{inner: stream}, nil
}

// GetModelInfo 获取模型信息
func (c *SDKClient) GetModelInfo() *types.ModelInfo {
	return c.inner.GetModelInfo()
}

// renderSystemPromptTemplate 渲染系统提示词模板
// 支持 {{question}} 宏，如果没有该宏，则将问题追加到提示词最后
func renderSystemPromptTemplate(template, question string) string {
	if template == "" {
		return ""
	}

	// 检查模板中是否包含 {{question}} 宏
	if containsQuestionMacro(template) {
		// 替换 {{question}} 宏
		return replaceQuestionMacro(template, question)
	}

	// 如果没有 {{question}} 宏，将问题追加到提示词最后
	return template + "\n\n用户问题：" + question
}

// containsQuestionMacro 检查模板中是否包含 {{question}} 宏
func containsQuestionMacro(template string) bool {
	return len(template) > 0 && (findSubstr(template, "{{question}}") >= 0 ||
		findSubstr(template, "{{ question }}") >= 0)
}

// replaceQuestionMacro 替换 {{question}} 宏
func replaceQuestionMacro(template, question string) string {
	result := template
	result = replaceAll(result, "{{question}}", question)
	result = replaceAll(result, "{{ question }}", question)
	return result
}

// findSubstr 查找子字符串位置，不存在返回 -1
func findSubstr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// replaceAll 替换所有匹配的子字符串
func replaceAll(s, old, new string) string {
	if old == "" {
		return s
	}
	result := s
	for {
		idx := findSubstr(result, old)
		if idx < 0 {
			break
		}
		result = result[:idx] + new + result[idx+len(old):]
	}
	return result
}

// ChatResult 对话结果
type ChatResult struct {
	Content string // 完整响应内容
	Stream  bool   // 是否来自流式响应
}

// StreamReader 流式读取器包装器
type StreamReader struct {
	inner *schema.StreamReader[*schema.Message]
}

// Recv 接收流式数据
func (s *StreamReader) Recv() (*schema.Message, error) {
	return s.inner.Recv()
}

// Close 关闭流
func (s *StreamReader) Close() {
	s.inner.Close()
}
