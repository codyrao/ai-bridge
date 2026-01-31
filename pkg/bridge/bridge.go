package bridge

import (
	"ai-bridge/pkg/adapters"
	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// NewAIClient 创建AI客户端
// provider: AI厂商类型
// modelName: 模型名称
// opts: 配置选项
func NewAIClient(provider types.Provider, modelName string, opts ...options.Option) (types.AIBridge, error) {
	return adapters.GetAdapter(provider, modelName, opts...)
}

// MustNewAIClient 创建AI客户端（出错时panic）
func MustNewAIClient(provider types.Provider, modelName string, opts ...options.Option) types.AIBridge {
	client, err := NewAIClient(provider, modelName, opts...)
	if err != nil {
		panic(err)
	}
	return client
}

// GetProviders 获取所有支持的厂商列表
func GetProviders() []types.Provider {
	providers := make([]types.Provider, 0, len(types.ModelRegistry)+1)
	for provider := range types.ModelRegistry {
		providers = append(providers, provider)
	}
	// 添加Ollama（不在ModelRegistry中，因为它支持任意模型）
	providers = append(providers, types.ProviderOllama)
	return providers
}

// GetModels 获取指定厂商的所有模型
func GetModels(provider types.Provider) []types.ModelInfo {
	return types.GetModelsByProvider(provider)
}

// IsValidProvider 检查厂商是否有效
func IsValidProvider(provider types.Provider) bool {
	return types.IsValidProvider(provider)
}

// IsValidModel 检查模型是否有效
func IsValidModel(provider types.Provider, modelName string) bool {
	return types.IsValidModel(provider, modelName)
}

// GetModelInfo 获取模型信息
func GetModelInfo(provider types.Provider, modelName string) *types.ModelInfo {
	return types.GetModelInfo(provider, modelName)
}
