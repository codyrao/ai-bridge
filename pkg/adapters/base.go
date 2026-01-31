package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// BaseAdapter 基础适配器
type BaseAdapter struct {
	Provider  types.Provider
	ModelName string
	Config    *types.Config
	ModelInfo *types.ModelInfo
	ChatModel model.ChatModel
}

// Chat 执行对话（非流式）
func (b *BaseAdapter) Chat(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	if b.ChatModel == nil {
		return nil, fmt.Errorf("chat model not initialized")
	}

	// 如果有系统提示词，添加到消息列表开头
	messages = b.prependSystemMessage(messages)

	return b.ChatModel.Generate(ctx, messages)
}

// ChatStream 执行对话（流式）
func (b *BaseAdapter) ChatStream(ctx context.Context, messages []*schema.Message) (*schema.StreamReader[*schema.Message], error) {
	if b.ChatModel == nil {
		return nil, fmt.Errorf("chat model not initialized")
	}

	// 如果有系统提示词，添加到消息列表开头
	messages = b.prependSystemMessage(messages)

	return b.ChatModel.Stream(ctx, messages)
}

// Generate 生成文本（简化接口）
func (b *BaseAdapter) Generate(ctx context.Context, prompt string) (string, error) {
	// 渲染系统提示词模板
	systemPrompt := b.renderSystemPromptTemplate(prompt)

	messages := []*schema.Message{}

	// 如果有系统提示词，添加到消息列表开头
	if systemPrompt != "" {
		messages = append(messages, schema.SystemMessage(systemPrompt))
	}

	messages = append(messages, schema.UserMessage(prompt))

	resp, err := b.ChatModel.Generate(ctx, messages)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// GenerateStream 生成文本（流式）
func (b *BaseAdapter) GenerateStream(ctx context.Context, prompt string) (string, error) {
	// 渲染系统提示词模板
	systemPrompt := b.renderSystemPromptTemplate(prompt)

	messages := []*schema.Message{}

	// 如果有系统提示词，添加到消息列表开头
	if systemPrompt != "" {
		messages = append(messages, schema.SystemMessage(systemPrompt))
	}

	messages = append(messages, schema.UserMessage(prompt))

	stream, err := b.ChatModel.Stream(ctx, messages)
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

// GetModelInfo 获取当前模型信息
func (b *BaseAdapter) GetModelInfo() *types.ModelInfo {
	return b.ModelInfo
}

// prependSystemMessage 如果有系统提示词，添加到消息列表开头
func (b *BaseAdapter) prependSystemMessage(messages []*schema.Message) []*schema.Message {
	if b.Config == nil || b.Config.SystemPrompt == "" {
		return messages
	}

	// 检查消息列表是否已经有系统消息
	for _, msg := range messages {
		if msg.Role == schema.System {
			return messages
		}
	}

	// 在消息列表开头添加系统消息
	systemMsg := schema.SystemMessage(b.Config.SystemPrompt)
	newMessages := make([]*schema.Message, 0, len(messages)+1)
	newMessages = append(newMessages, systemMsg)
	newMessages = append(newMessages, messages...)
	return newMessages
}

// renderSystemPromptTemplate 渲染系统提示词模板
// 支持 {{question}} 宏，如果没有该宏，则将问题追加到提示词最后
func (b *BaseAdapter) renderSystemPromptTemplate(question string) string {
	if b.Config == nil || b.Config.SystemPrompt == "" {
		return ""
	}

	template := b.Config.SystemPrompt

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

// AdapterFactory 适配器工厂函数类型
type AdapterFactory func(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error)

// adapterRegistry 适配器注册表
var adapterRegistry = make(map[types.Provider]AdapterFactory)

// RegisterAdapter 注册适配器工厂
func RegisterAdapter(provider types.Provider, factory AdapterFactory) {
	adapterRegistry[provider] = factory
}

// GetAdapter 获取适配器
func GetAdapter(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	// 验证厂商
	if !types.IsValidProvider(provider) {
		return nil, fmt.Errorf("invalid provider: %s", provider)
	}

	// 验证模型
	if !types.IsValidModel(provider, modelName) {
		// 如果是Ollama，允许任意模型名
		if provider != types.ProviderOllama {
			return nil, fmt.Errorf("invalid model %s for provider %s", modelName, provider)
		}
	}

	// 获取适配器工厂
	factory, ok := adapterRegistry[provider]
	if !ok {
		return nil, fmt.Errorf("adapter not found for provider: %s", provider)
	}

	return factory(provider, modelName, opts...)
}

// ParseToolArguments 解析工具参数
func ParseToolArguments(arguments string) (map[string]interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return nil, fmt.Errorf("failed to parse tool arguments: %w", err)
	}
	return params, nil
}
