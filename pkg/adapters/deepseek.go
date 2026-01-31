package adapters

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/deepseek"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// DeepseekAdapter Deepseek适配器
type DeepseekAdapter struct {
	BaseAdapter
}

// NewDeepseekAdapter 创建Deepseek适配器
func NewDeepseekAdapter(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	cfg := options.ApplyOptions(opts...)
	modelInfo := types.GetModelInfo(provider, modelName)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key is required for Deepseek")
	}

	// 创建Deepseek配置
	deepseekCfg := &deepseek.ChatModelConfig{
		APIKey:      cfg.APIKey,
		Model:       modelName,
		MaxTokens:   cfg.MaxTokens,
		Temperature: cfg.Temperature,
		TopP:        cfg.TopP,
	}

	// 设置基础URL（如果使用自定义端点）
	if cfg.BaseURL != "" {
		deepseekCfg.BaseURL = cfg.BaseURL
	}

	// 创建ChatModel
	chatModel, err := deepseek.NewChatModel(context.Background(), deepseekCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create deepseek chat model: %w", err)
	}

	adapter := &DeepseekAdapter{
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
	RegisterAdapter(types.ProviderDeepseek, NewDeepseekAdapter)
}
