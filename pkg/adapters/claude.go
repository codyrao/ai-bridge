package adapters

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// ClaudeAdapter Claude适配器（使用OpenAI兼容接口）
type ClaudeAdapter struct {
	BaseAdapter
}

// NewClaudeAdapter 创建Claude适配器
func NewClaudeAdapter(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	cfg := options.ApplyOptions(opts...)
	modelInfo := types.GetModelInfo(provider, modelName)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key is required for Claude")
	}

	// Claude可以通过多种方式调用，这里使用OpenAI兼容接口
	// 用户可以通过BaseURL指向Anthropic API或兼容服务
	baseURL := cfg.BaseURL
	if baseURL == "" {
		// 默认使用Anthropic API
		baseURL = "https://api.anthropic.com/v1"
	}

	// 创建OpenAI配置（兼容模式）
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
		return nil, fmt.Errorf("failed to create claude chat model: %w", err)
	}

	adapter := &ClaudeAdapter{
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
	RegisterAdapter(types.ProviderClaude, NewClaudeAdapter)
}
