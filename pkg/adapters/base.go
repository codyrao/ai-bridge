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

	return b.ChatModel.Generate(ctx, messages)
}

// ChatStream 执行对话（流式）
func (b *BaseAdapter) ChatStream(ctx context.Context, messages []*schema.Message) (*schema.StreamReader[*schema.Message], error) {
	if b.ChatModel == nil {
		return nil, fmt.Errorf("chat model not initialized")
	}

	return b.ChatModel.Stream(ctx, messages)
}

// Generate 生成文本（简化接口）
func (b *BaseAdapter) Generate(ctx context.Context, prompt string) (string, error) {
	messages := []*schema.Message{
		schema.UserMessage(prompt),
	}

	resp, err := b.Chat(ctx, messages)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// GenerateStream 生成文本（流式）
func (b *BaseAdapter) GenerateStream(ctx context.Context, prompt string) (string, error) {
	messages := []*schema.Message{
		schema.UserMessage(prompt),
	}

	stream, err := b.ChatStream(ctx, messages)
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
