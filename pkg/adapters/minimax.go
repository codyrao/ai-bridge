package adapters

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// MiniMaxAdapter MiniMax适配器（使用OpenAI兼容接口）
type MiniMaxAdapter struct {
	BaseAdapter
}

// NewMiniMaxAdapter 创建MiniMax适配器
func NewMiniMaxAdapter(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	cfg := options.ApplyOptions(opts...)
	modelInfo := types.GetModelInfo(provider, modelName)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key is required for MiniMax")
	}

	// MiniMax使用OpenAI兼容接口
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.minimax.chat/v1"
	}

	// 创建OpenAI配置
	openaiCfg := &openai.ChatModelConfig{
		APIKey:      cfg.APIKey,
		BaseURL:     baseURL,
		Model:       modelName,
		MaxTokens:   &cfg.MaxTokens,
		Temperature: &cfg.Temperature,
		TopP:        &cfg.TopP,
	}

	// 创建ChatModel
	chatModel, err := openai.NewChatModel(context.Background(), openaiCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create minimax chat model: %w", err)
	}

	adapter := &MiniMaxAdapter{
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
	RegisterAdapter(types.ProviderMiniMax, NewMiniMaxAdapter)
}
