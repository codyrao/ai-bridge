package bridge

import (
	"os"
	"time"

	"github.com/cloudwego/eino/components/tool"

	"ai-bridge/pkg/adapters"
	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// SDKConfig SDK全局配置
type SDKConfig struct {
	GPT      ProviderConfig // OpenAI GPT配置
	QWen     ProviderConfig // 阿里通义千问配置
	Kimi     ProviderConfig // Moonshot Kimi配置
	GLM      ProviderConfig // 智谱GLM配置
	MiniMax  ProviderConfig // MiniMax配置
	Claude   ProviderConfig // Anthropic Claude配置
	Gemini   ProviderConfig // Google Gemini配置
	Grok     ProviderConfig // xAI Grok配置
	Deepseek ProviderConfig // Deepseek配置
	Ollama   ProviderConfig // Ollama本地模型配置
}

// ProviderConfig 厂商配置
type ProviderConfig struct {
	APIKey      string        // API密钥
	BaseURL     string        // 自定义API地址（可选）
	Timeout     time.Duration // 请求超时（默认120s）
	MaxRetries  int           // 最大重试次数（默认3）
	Temperature float32       // 温度参数（默认0.7）
	TopP        float32       // Top P参数（默认0.9）
	MaxTokens   int           // 最大Token数（默认2048）
	Proxy       string        // 代理地址（可选）
}

// SDK AI Bridge SDK
type SDK struct {
	config *SDKConfig
}

// NewSDK 创建SDK实例
func NewSDK(config *SDKConfig) *SDK {
	return &SDK{config: config}
}

// CreateClient 创建AI客户端
func (s *SDK) CreateClient(provider types.Provider, modelName string) (types.AIBridge, error) {
	opts := s.buildOptions(provider)
	return adapters.GetAdapter(provider, modelName, opts...)
}

// CreateClientWithTools 创建带MCP工具的AI客户端
func (s *SDK) CreateClientWithTools(provider types.Provider, modelName string, tools []tool.BaseTool) (types.AIBridge, error) {
	opts := s.buildOptions(provider)
	opts = append(opts, options.WithTools(tools...))
	return adapters.GetAdapter(provider, modelName, opts...)
}

// CreateSDKClient 创建支持高级Option的SDK客户端
// 支持Generate、Chat等方法的高级配置
func (s *SDK) CreateSDKClient(provider types.Provider, modelName string) (*SDKClient, error) {
	inner, err := s.CreateClient(provider, modelName)
	if err != nil {
		return nil, err
	}
	return NewSDKClient(inner), nil
}

// CreateSDKClientWithTools 创建带工具的高级SDK客户端
func (s *SDK) CreateSDKClientWithTools(provider types.Provider, modelName string, tools []tool.BaseTool) (*SDKClient, error) {
	inner, err := s.CreateClientWithTools(provider, modelName, tools)
	if err != nil {
		return nil, err
	}
	return NewSDKClient(inner), nil
}

// buildOptions 根据厂商配置构建选项
func (s *SDK) buildOptions(provider types.Provider) []options.Option {
	var cfg ProviderConfig

	switch provider {
	case types.ProviderGPT:
		cfg = s.config.GPT
	case types.ProviderQWen:
		cfg = s.config.QWen
	case types.ProviderKimi:
		cfg = s.config.Kimi
	case types.ProviderGLM:
		cfg = s.config.GLM
	case types.ProviderMiniMax:
		cfg = s.config.MiniMax
	case types.ProviderClaude:
		cfg = s.config.Claude
	case types.ProviderGemini:
		cfg = s.config.Gemini
	case types.ProviderGrok:
		cfg = s.config.Grok
	case types.ProviderDeepseek:
		cfg = s.config.Deepseek
	case types.ProviderOllama:
		cfg = s.config.Ollama
	}

	opts := []options.Option{
		options.WithAPIKey(cfg.APIKey),
	}

	if cfg.BaseURL != "" {
		opts = append(opts, options.WithBaseURL(cfg.BaseURL))
	}
	if cfg.Timeout > 0 {
		opts = append(opts, options.WithTimeout(cfg.Timeout))
	}
	if cfg.MaxRetries > 0 {
		opts = append(opts, options.WithMaxRetries(cfg.MaxRetries))
	}
	if cfg.Temperature > 0 {
		opts = append(opts, options.WithTemperature(cfg.Temperature))
	}
	if cfg.TopP > 0 {
		opts = append(opts, options.WithTopP(cfg.TopP))
	}
	if cfg.MaxTokens > 0 {
		opts = append(opts, options.WithMaxTokens(cfg.MaxTokens))
	}
	if cfg.Proxy != "" {
		opts = append(opts, options.WithProxy(cfg.Proxy))
	}

	return opts
}

// ConfigFromEnv 从环境变量加载配置
func ConfigFromEnv() *SDKConfig {
	return &SDKConfig{
		GPT: ProviderConfig{
			APIKey: os.Getenv("GPT_API_KEY"),
		},
		QWen: ProviderConfig{
			APIKey: os.Getenv("QWEN_API_KEY"),
		},
		Kimi: ProviderConfig{
			APIKey: os.Getenv("KIMI_API_KEY"),
		},
		GLM: ProviderConfig{
			APIKey: os.Getenv("GLM_API_KEY"),
		},
		MiniMax: ProviderConfig{
			APIKey: os.Getenv("MINIMAX_API_KEY"),
		},
		Claude: ProviderConfig{
			APIKey: os.Getenv("CLAUDE_API_KEY"),
		},
		Gemini: ProviderConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		},
		Grok: ProviderConfig{
			APIKey: os.Getenv("GROK_API_KEY"),
		},
		Deepseek: ProviderConfig{
			APIKey: os.Getenv("DEEPSEEK_API_KEY"),
		},
		Ollama: ProviderConfig{
			BaseURL: getEnvOrDefault("OLLAMA_BASE_URL", "http://localhost:11434"),
		},
	}
}

// getEnvOrDefault 获取环境变量，如果不存在返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
