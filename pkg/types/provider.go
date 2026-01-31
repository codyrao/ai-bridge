package types

// Provider 定义AI模型厂商类型
type Provider string

const (
	ProviderQWen     Provider = "qwen"     // 阿里通义千问
	ProviderKimi     Provider = "kimi"     // Moonshot Kimi
	ProviderGLM      Provider = "glm"      // 智谱GLM
	ProviderMiniMax  Provider = "minimax"  // MiniMax
	ProviderClaude   Provider = "claude"   // Anthropic Claude
	ProviderGPT      Provider = "gpt"      // OpenAI GPT
	ProviderGemini   Provider = "gemini"   // Google Gemini
	ProviderGrok     Provider = "grok"     // xAI Grok
	ProviderDeepseek Provider = "deepseek" // Deepseek
	ProviderOllama   Provider = "ollama"   // Ollama本地模型
)

// ModelInfo 模型信息
type ModelInfo struct {
	Name        string
	Provider    Provider
	MaxTokens   int
	Description string
}

// ModelRegistry 模型注册表
var ModelRegistry = map[Provider][]ModelInfo{
	ProviderQWen: {
		{Name: "qwen-turbo", Provider: ProviderQWen, MaxTokens: 8192, Description: "通义千问Turbo"},
		{Name: "qwen-plus", Provider: ProviderQWen, MaxTokens: 32768, Description: "通义千问Plus"},
		{Name: "qwen-max", Provider: ProviderQWen, MaxTokens: 32768, Description: "通义千问Max"},
		{Name: "qwen-coder-plus", Provider: ProviderQWen, MaxTokens: 32768, Description: "通义千问代码模型"},
	},
	ProviderKimi: {
		{Name: "moonshot-v1-8k", Provider: ProviderKimi, MaxTokens: 8192, Description: "Kimi 8K"},
		{Name: "moonshot-v1-32k", Provider: ProviderKimi, MaxTokens: 32768, Description: "Kimi 32K"},
		{Name: "moonshot-v1-128k", Provider: ProviderKimi, MaxTokens: 131072, Description: "Kimi 128K"},
	},
	ProviderGLM: {
		{Name: "glm-4", Provider: ProviderGLM, MaxTokens: 8192, Description: "GLM-4"},
		{Name: "glm-4-plus", Provider: ProviderGLM, MaxTokens: 8192, Description: "GLM-4 Plus"},
		{Name: "glm-4-flash", Provider: ProviderGLM, MaxTokens: 8192, Description: "GLM-4 Flash"},
		{Name: "glm-4v", Provider: ProviderGLM, MaxTokens: 2048, Description: "GLM-4V多模态"},
	},
	ProviderMiniMax: {
		{Name: "abab6.5s-chat", Provider: ProviderMiniMax, MaxTokens: 8192, Description: "MiniMax 6.5s"},
		{Name: "abab6.5-chat", Provider: ProviderMiniMax, MaxTokens: 8192, Description: "MiniMax 6.5"},
		{Name: "abab6-chat", Provider: ProviderMiniMax, MaxTokens: 8192, Description: "MiniMax 6"},
	},
	ProviderClaude: {
		{Name: "claude-3-opus-20240229", Provider: ProviderClaude, MaxTokens: 200000, Description: "Claude 3 Opus"},
		{Name: "claude-3-sonnet-20240229", Provider: ProviderClaude, MaxTokens: 200000, Description: "Claude 3 Sonnet"},
		{Name: "claude-3-haiku-20240307", Provider: ProviderClaude, MaxTokens: 200000, Description: "Claude 3 Haiku"},
		{Name: "claude-3-5-sonnet-20240620", Provider: ProviderClaude, MaxTokens: 200000, Description: "Claude 3.5 Sonnet"},
	},
	ProviderGPT: {
		{Name: "gpt-3.5-turbo", Provider: ProviderGPT, MaxTokens: 16385, Description: "GPT-3.5 Turbo"},
		{Name: "gpt-4", Provider: ProviderGPT, MaxTokens: 8192, Description: "GPT-4"},
		{Name: "gpt-4-turbo", Provider: ProviderGPT, MaxTokens: 128000, Description: "GPT-4 Turbo"},
		{Name: "gpt-4o", Provider: ProviderGPT, MaxTokens: 128000, Description: "GPT-4o"},
		{Name: "gpt-4o-mini", Provider: ProviderGPT, MaxTokens: 128000, Description: "GPT-4o Mini"},
	},
	ProviderGemini: {
		{Name: "gemini-1.5-pro", Provider: ProviderGemini, MaxTokens: 2097152, Description: "Gemini 1.5 Pro"},
		{Name: "gemini-1.5-flash", Provider: ProviderGemini, MaxTokens: 1048576, Description: "Gemini 1.5 Flash"},
		{Name: "gemini-1.0-pro", Provider: ProviderGemini, MaxTokens: 32768, Description: "Gemini 1.0 Pro"},
	},
	ProviderGrok: {
		{Name: "grok-1", Provider: ProviderGrok, MaxTokens: 131072, Description: "Grok-1"},
		{Name: "grok-2", Provider: ProviderGrok, MaxTokens: 131072, Description: "Grok-2"},
	},
	ProviderDeepseek: {
		{Name: "deepseek-chat", Provider: ProviderDeepseek, MaxTokens: 32768, Description: "Deepseek Chat"},
		{Name: "deepseek-coder", Provider: ProviderDeepseek, MaxTokens: 32768, Description: "Deepseek Coder"},
		{Name: "deepseek-reasoner", Provider: ProviderDeepseek, MaxTokens: 32768, Description: "Deepseek Reasoner"},
	},
}

// GetModelsByProvider 获取指定厂商的所有模型
func GetModelsByProvider(provider Provider) []ModelInfo {
	if models, ok := ModelRegistry[provider]; ok {
		return models
	}
	return nil
}

// GetModelInfo 获取指定模型信息
func GetModelInfo(provider Provider, modelName string) *ModelInfo {
	if models, ok := ModelRegistry[provider]; ok {
		for _, m := range models {
			if m.Name == modelName {
				return &m
			}
		}
	}
	return nil
}

// IsValidProvider 检查厂商是否有效
func IsValidProvider(provider Provider) bool {
	// Ollama是特殊厂商，允许任意模型名
	if provider == ProviderOllama {
		return true
	}
	_, ok := ModelRegistry[provider]
	return ok
}

// IsValidModel 检查模型是否有效
func IsValidModel(provider Provider, modelName string) bool {
	return GetModelInfo(provider, modelName) != nil
}
