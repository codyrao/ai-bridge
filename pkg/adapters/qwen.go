package adapters

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/qwen"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// QWenAdapter 通义千问适配器
type QWenAdapter struct {
	BaseAdapter
}

// NewQWenAdapter 创建通义千问适配器
func NewQWenAdapter(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	cfg := options.ApplyOptions(opts...)
	modelInfo := types.GetModelInfo(provider, modelName)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key is required for QWen")
	}

	// 创建QWen ChatModel配置
	qwenCfg := &qwen.ChatModelConfig{
		APIKey:      cfg.APIKey,
		Model:       modelName,
		MaxTokens:   &cfg.MaxTokens,
		Temperature: &cfg.Temperature,
		TopP:        &cfg.TopP,
	}

	// 设置基础URL（如果使用自定义端点）
	if cfg.BaseURL != "" {
		qwenCfg.BaseURL = cfg.BaseURL
	} else {
		qwenCfg.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	// 创建ChatModel
	chatModel, err := qwen.NewChatModel(context.Background(), qwenCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create qwen chat model: %w", err)
	}

	adapter := &QWenAdapter{
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
	RegisterAdapter(types.ProviderQWen, NewQWenAdapter)
}
