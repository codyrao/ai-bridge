package adapters

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// GLMAdapter 智谱GLM适配器（使用OpenAI兼容接口）
type GLMAdapter struct {
	BaseAdapter
}

// NewGLMAdapter 创建GLM适配器
func NewGLMAdapter(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	cfg := options.ApplyOptions(opts...)
	modelInfo := types.GetModelInfo(provider, modelName)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key is required for GLM")
	}

	// GLM使用OpenAI兼容接口
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://open.bigmodel.cn/api/paas/v4"
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
		return nil, fmt.Errorf("failed to create glm chat model: %w", err)
	}

	adapter := &GLMAdapter{
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
	RegisterAdapter(types.ProviderGLM, NewGLMAdapter)
}
