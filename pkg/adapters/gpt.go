package adapters

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// GPTAdapter OpenAI GPT适配器
type GPTAdapter struct {
	BaseAdapter
}

// NewGPTAdapter 创建GPT适配器
func NewGPTAdapter(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	cfg := options.ApplyOptions(opts...)
	modelInfo := types.GetModelInfo(provider, modelName)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key is required for GPT")
	}

	// 创建OpenAI配置
	openaiCfg := &openai.ChatModelConfig{
		APIKey:      cfg.APIKey,
		Model:       modelName,
		MaxTokens:   &cfg.MaxTokens,
		Temperature: &cfg.Temperature,
		TopP:        &cfg.TopP,
	}

	// 设置基础URL（如果使用自定义端点）
	if cfg.BaseURL != "" {
		openaiCfg.BaseURL = cfg.BaseURL
	}

	// 创建ChatModel
	chatModel, err := openai.NewChatModel(context.Background(), openaiCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create gpt chat model: %w", err)
	}

	adapter := &GPTAdapter{
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
	RegisterAdapter(types.ProviderGPT, NewGPTAdapter)
}
