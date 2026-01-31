package options

import (
	"time"

	"ai-bridge/pkg/types"

	"github.com/cloudwego/eino/components/tool"
)

// Option 配置选项函数类型
type Option func(*types.Config)

// WithAPIKey 设置API密钥
func WithAPIKey(apiKey string) Option {
	return func(c *types.Config) {
		c.APIKey = apiKey
	}
}

// WithBaseURL 设置基础URL
func WithBaseURL(url string) Option {
	return func(c *types.Config) {
		c.BaseURL = url
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *types.Config) {
		c.Timeout = timeout
	}
}

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(retries int) Option {
	return func(c *types.Config) {
		c.MaxRetries = retries
	}
}

// WithTemperature 设置温度参数
func WithTemperature(temp float32) Option {
	return func(c *types.Config) {
		c.Temperature = temp
	}
}

// WithTopP 设置Top P参数
func WithTopP(topP float32) Option {
	return func(c *types.Config) {
		c.TopP = topP
	}
}

// WithMaxTokens 设置最大token数
func WithMaxTokens(tokens int) Option {
	return func(c *types.Config) {
		c.MaxTokens = tokens
	}
}

// WithTools 设置MCP工具
func WithTools(tools ...tool.BaseTool) Option {
	return func(c *types.Config) {
		c.Tools = tools
	}
}

// WithStream 设置是否使用流式响应
func WithStream(stream bool) Option {
	return func(c *types.Config) {
		c.Stream = stream
	}
}

// WithExtraHeaders 设置额外请求头
func WithExtraHeaders(headers map[string]string) Option {
	return func(c *types.Config) {
		if c.ExtraHeaders == nil {
			c.ExtraHeaders = make(map[string]string)
		}
		for k, v := range headers {
			c.ExtraHeaders[k] = v
		}
	}
}

// WithExtraHeader 设置单个额外请求头
func WithExtraHeader(key, value string) Option {
	return func(c *types.Config) {
		if c.ExtraHeaders == nil {
			c.ExtraHeaders = make(map[string]string)
		}
		c.ExtraHeaders[key] = value
	}
}

// WithOrganization 设置组织ID（OpenAI专用）
func WithOrganization(org string) Option {
	return func(c *types.Config) {
		c.Organization = org
	}
}

// WithProxy 设置代理地址
func WithProxy(proxy string) Option {
	return func(c *types.Config) {
		c.Proxy = proxy
	}
}

// WithEnableLog 设置是否启用日志
func WithEnableLog(enable bool) Option {
	return func(c *types.Config) {
		c.EnableLog = enable
	}
}

// WithSystemPrompt 设置系统提示词模板
// 用于定义AI的角色和行为准则
func WithSystemPrompt(prompt string) Option {
	return func(c *types.Config) {
		c.SystemPrompt = prompt
	}
}

// WithSkills 设置要使用的 Agent Skills
// skills 是 Skill 名称列表，如 ["calculator", "weather"]
func WithSkills(skills ...string) Option {
	return func(c *types.Config) {
		c.Skills = append(c.Skills, skills...)
	}
}

// WithSkillPaths 设置 Skill 文件夹路径
// 可以指定多个包含 skill.md 的文件夹路径
func WithSkillPaths(paths ...string) Option {
	return func(c *types.Config) {
		c.SkillPaths = append(c.SkillPaths, paths...)
	}
}

// WithSkillDir 从目录加载所有 Skills
// 目录下的每个子文件夹应该是一个 Skill
func WithSkillDir(dir string) Option {
	return func(c *types.Config) {
		c.SkillPaths = append(c.SkillPaths, dir)
	}
}

// ApplyOptions 应用配置选项
func ApplyOptions(opts ...Option) *types.Config {
	config := types.DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return config
}
