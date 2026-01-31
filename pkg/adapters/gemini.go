package adapters

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// GeminiAdapter Google Gemini适配器
// 注意：Gemini原生需要使用Google genai客户端，这里使用OpenAI兼容接口简化
type GeminiAdapter struct {
	BaseAdapter
}

// NewGeminiAdapter 创建Gemini适配器
func NewGeminiAdapter(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	cfg := options.ApplyOptions(opts...)
	modelInfo := types.GetModelInfo(provider, modelName)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api key is required for Gemini")
	}

	// 使用OpenAI兼容接口
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
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
		return nil, fmt.Errorf("failed to create gemini chat model: %w", err)
	}

	adapter := &GeminiAdapter{
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
	RegisterAdapter(types.ProviderGemini, NewGeminiAdapter)
}
