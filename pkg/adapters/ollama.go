package adapters

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ollama"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// OllamaAdapter Ollama本地模型适配器
type OllamaAdapter struct {
	BaseAdapter
}

// NewOllamaAdapter 创建Ollama适配器
func NewOllamaAdapter(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	cfg := options.ApplyOptions(opts...)

	// Ollama允许任意模型名，不需要在注册表中验证
	modelInfo := &types.ModelInfo{
		Name:        modelName,
		Provider:    provider,
		MaxTokens:   cfg.MaxTokens,
		Description: "Ollama本地模型: " + modelName,
	}

	// 创建Ollama配置
	ollamaCfg := &ollama.ChatModelConfig{
		Model: modelName,
	}

	// 设置基础URL
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	ollamaCfg.BaseURL = baseURL

	// 设置超时
	if cfg.Timeout > 0 {
		ollamaCfg.Timeout = cfg.Timeout
	}

	// 创建ChatModel
	chatModel, err := ollama.NewChatModel(context.Background(), ollamaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama chat model: %w", err)
	}

	adapter := &OllamaAdapter{
		BaseAdapter: BaseAdapter{
			Provider:  provider,
			ModelName: modelName,
			Config:    cfg,
			ModelInfo: modelInfo,
			ChatModel: chatModel,
		},
	}

	return adapter, nil
}

func init() {
	RegisterAdapter(types.ProviderOllama, NewOllamaAdapter)
}
