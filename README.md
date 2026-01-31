# AI Bridge SDK

<div align="center">

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-blue)](https://hub.docker.com/r/codyrao/ai-bridge)

**English** | [中文](#中文文档)

A unified Go SDK for AI large language models based on ByteDance's Eino framework. Supports 10+ mainstream AI providers including OpenAI, QWen, Kimi, GLM, Claude, Gemini, Deepseek, and Ollama.

</div>

---

## Table of Contents

- [Features](#features)
- [Supported Providers](#supported-providers)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Advanced Usage](#advanced-usage)
- [Agent Skills](#agent-skills)
- [Docker Deployment](#docker-deployment)
- [License](#license)

## Features

- **Multi-Provider Support**: 10+ mainstream AI model providers
- **Unified Interface**: Consistent SDK API for easy model switching
- **Config-Based Setup**: Centralized configuration via `Config` struct
- **Option Pattern**: Flexible configuration with extensive options
- **MCP Tool Support**: Model Control Plane tool integration
- **Streaming & Non-Streaming**: Support both response modes
- **RESTful API**: HTTP API service for easy integration
- **Agent Skills**: Anthropic's open standard for AI capabilities

## Supported Providers

| Provider | Identifier | Models |
|----------|------------|--------|
| OpenAI GPT | `gpt` | GPT-3.5/4/4o series |
| Alibaba QWen | `qwen` | QWen Turbo/Plus/Max |
| Moonshot Kimi | `kimi` | Kimi 8K/32K/128K |
| Zhipu GLM | `glm` | GLM-4 series |
| MiniMax | `minimax` | MiniMax 6.5 series |
| Anthropic Claude | `claude` | Claude 3 series |
| Google Gemini | `gemini` | Gemini 1.5 Pro/Flash |
| xAI Grok | `grok` | Grok-1/2 |
| Deepseek | `deepseek` | Deepseek Chat/Coder |
| Ollama | `ollama` | Local models support |

## Installation

```bash
go get github.com/codyrao/ai-bridge
```

## Quick Start

### Method 1: Using Config Struct (Recommended)

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/codyrao/ai-bridge/pkg/bridge"
    "github.com/codyrao/ai-bridge/pkg/types"
)

func main() {
    // Create SDK configuration
    config := &bridge.SDKConfig{
        GPT: bridge.ProviderConfig{
            APIKey: "your-openai-api-key",
        },
        QWen: bridge.ProviderConfig{
            APIKey: "your-qwen-api-key",
        },
        Ollama: bridge.ProviderConfig{
            BaseURL: "http://localhost:11434",
        },
    }

    // Create SDK client
    sdk := bridge.NewSDK(config)

    // Use GPT model
    client, err := sdk.CreateClient(types.ProviderGPT, "gpt-3.5-turbo")
    if err != nil {
        log.Fatal(err)
    }

    // Generate response
    resp, err := client.Generate(context.Background(), "Hello, introduce yourself")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp)
}
```

### Method 2: Using Option Pattern

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/codyrao/ai-bridge/pkg/bridge"
    "github.com/codyrao/ai-bridge/pkg/options"
    "github.com/codyrao/ai-bridge/pkg/types"
)

func main() {
    // Create client with options
    client, err := bridge.NewAIClient(
        types.ProviderGPT,
        "gpt-4",
        options.WithAPIKey("your-api-key"),
        options.WithTemperature(0.7),
        options.WithMaxTokens(2048),
    )
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Generate(context.Background(), "Hello")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp)
}
```

## Configuration

### SDKConfig Structure

```go
type SDKConfig struct {
    GPT      ProviderConfig  // OpenAI GPT config
    QWen     ProviderConfig  // Alibaba QWen config
    Kimi     ProviderConfig  // Moonshot Kimi config
    GLM      ProviderConfig  // Zhipu GLM config
    MiniMax  ProviderConfig  // MiniMax config
    Claude   ProviderConfig  // Anthropic Claude config
    Gemini   ProviderConfig  // Google Gemini config
    Grok     ProviderConfig  // xAI Grok config
    Deepseek ProviderConfig  // Deepseek config
    Ollama   ProviderConfig  // Ollama local model config
}

type ProviderConfig struct {
    APIKey      string        // API key
    BaseURL     string        // Custom API endpoint (optional)
    Timeout     time.Duration // Request timeout (default 120s)
    MaxRetries  int           // Max retry attempts (default 3)
    Temperature float32       // Temperature (default 0.7)
    TopP        float32       // Top P (default 0.9)
    MaxTokens   int           // Max tokens (default 2048)
    Proxy       string        // Proxy address (optional)
}
```

### Environment Variables

```bash
export GPT_API_KEY="your-openai-api-key"
export QWEN_API_KEY="your-qwen-api-key"
export KIMI_API_KEY="your-kimi-api-key"
export DEEPSEEK_API_KEY="your-deepseek-api-key"
```

```go
// Load from environment variables
config := bridge.ConfigFromEnv()
sdk := bridge.NewSDK(config)
```

## Advanced Usage

### System Prompt Templates

Support `{{question}}` macro for precise question placement:

```go
template := `You are a Go expert. Explain: {{question}}

Requirements:
1. Provide code examples
2. Explain use cases
3. Point out common mistakes`

resp, _ := client.Generate(ctx, "What is goroutine", bridge.WithSystemPrompt(template))
```

### Streaming Response

```go
// Stream generation
stream, _ := client.GenerateStream(context.Background(), "Tell me a story")
defer stream.Close()

for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    fmt.Print(msg.Content)  // Real-time output
}
```

### MCP Tools

```go
import (
    "ai-bridge/pkg/mcp"
    "github.com/cloudwego/eino/components/tool"
)

// Create tool
weatherTool := mcp.NewMCPTool(
    "get_weather",
    "Query weather for a city",
    mcp.CreateParameterSchema(
        mcp.CreateStringProperty("city", "City name", true),
    ),
    func(ctx context.Context, params map[string]interface{}) (string, error) {
        city := params["city"].(string)
        return fmt.Sprintf("Weather in %s: Sunny, 25°C", city), nil
    },
)

// Create client with tools
client, _ := sdk.CreateClientWithTools(
    types.ProviderGPT,
    "gpt-4",
    []tool.BaseTool{weatherTool.ToEinoTool()},
)
```

## Agent Skills

Agent Skills is Anthropic's open standard for providing reusable capabilities to AI agents.

### Skill Structure

```
skills/
├── calculator/
│   └── SKILL.md
├── weather/
│   └── SKILL.md
└── code-search/
    └── SKILL.md
```

### SKILL.md Format

```markdown
---
name: calculator
description: Perform mathematical calculations
version: "1.0"
metadata:
  author: ai-bridge
  tags: math, calculation
---

# Calculator Skill

## Overview
This skill helps AI assistants perform mathematical calculations.

## Usage
- Basic operations: addition, subtraction, multiplication, division
- Advanced: power, square root

## Examples
**User**: Calculate 123 + 456
**AI**: 123 + 456 = 579
```

### Loading Skills

```go
import "ai-bridge/pkg/skills"

// Load single skill
skill, _ := skills.LoadSkill("skills/calculator")
systemPrompt := skill.GetSystemPrompt()

// Load all skills from directory
registry := skills.NewRegistry()
registry.LoadFromDir("skills")

// Get specific skill
calcSkill, ok := registry.Get("calculator")
```

## Docker Deployment

### Using Docker

```bash
# Pull image
docker pull codyrao/ai-bridge:latest

# Run container
docker run -d \
  -p 8080:8080 \
  -e GPT_API_KEY="your-openai-api-key" \
  -e QWEN_API_KEY="your-qwen-api-key" \
  --name ai-bridge \
  codyrao/ai-bridge:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  ai-bridge:
    image: codyrao/ai-bridge:latest
    ports:
      - "8080:8080"
    environment:
      - GPT_API_KEY=${GPT_API_KEY}
      - QWEN_API_KEY=${QWEN_API_KEY}
      - DEEPSEEK_API_KEY=${DEEPSEEK_API_KEY}
    restart: unless-stopped
```

## License

MIT License

---

# 中文文档

AI Bridge SDK 是一个基于 Go 语言和字节跳动 Eino 框架实现的 AI 大模型接口聚合SDK。

## 功能特性

- **多厂商支持**：支持 10+ 主流 AI 模型厂商
- **统一接口**：提供一致的 SDK 接口，简化多模型切换
- **Config配置模式**：通过Config结构统一配置各模型API和密钥
- **Option配置模式**：丰富的配置选项，灵活控制模型行为
- **MCP 工具支持**：支持 Model Control Plane 工具调用
- **流式响应**：支持流式和非流式两种对话模式
- **Agent Skills**：支持 Anthropic 的 Agent Skills 开放标准

## 支持的模型厂商

| 厂商 | 标识符 | 说明 |
|------|--------|------|
| OpenAI GPT | `gpt` | GPT-3.5/4/4o 系列 |
| 阿里通义千问 | `qwen` | QWen Turbo/Plus/Max |
| Moonshot Kimi | `kimi` | Kimi 8K/32K/128K |
| 智谱 GLM | `glm` | GLM-4 系列 |
| MiniMax | `minimax` | MiniMax 6.5 系列 |
| Anthropic Claude | `claude` | Claude 3 系列 |
| Google Gemini | `gemini` | Gemini 1.5 Pro/Flash |
| xAI Grok | `grok` | Grok-1/2 |
| Deepseek | `deepseek` | Deepseek Chat/Coder |
| Ollama | `ollama` | 本地模型支持 |

## 安装

```bash
go get github.com/codyrao/ai-bridge
```

## 快速开始

### 方式一：使用Config结构配置（推荐）

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/codyrao/ai-bridge/pkg/bridge"
    "github.com/codyrao/ai-bridge/pkg/types"
)

func main() {
    // 创建SDK配置
    config := &bridge.SDKConfig{
        GPT: bridge.ProviderConfig{
            APIKey: "your-openai-api-key",
        },
        QWen: bridge.ProviderConfig{
            APIKey: "your-qwen-api-key",
        },
        Ollama: bridge.ProviderConfig{
            BaseURL: "http://localhost:11434",
        },
    }

    // 创建SDK客户端
    sdk := bridge.NewSDK(config)

    // 使用GPT模型
    client, err := sdk.CreateClient(types.ProviderGPT, "gpt-3.5-turbo")
    if err != nil {
        log.Fatal(err)
    }

    // 执行对话
    resp, err := client.Generate(context.Background(), "你好，请介绍一下自己")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp)
}
```

### 方式二：使用Option模式配置

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/codyrao/ai-bridge/pkg/bridge"
    "github.com/codyrao/ai-bridge/pkg/options"
    "github.com/codyrao/ai-bridge/pkg/types"
)

func main() {
    // 创建客户端
    client, err := bridge.NewAIClient(
        types.ProviderGPT,
        "gpt-4",
        options.WithAPIKey("your-api-key"),
        options.WithTemperature(0.7),
        options.WithMaxTokens(2048),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 执行对话
    resp, err := client.Generate(context.Background(), "你好")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp)
}
```

## SDK配置详解

### Config结构定义

```go
type SDKConfig struct {
    GPT      ProviderConfig  // OpenAI GPT配置
    QWen     ProviderConfig  // 阿里通义千问配置
    Kimi     ProviderConfig  // Moonshot Kimi配置
    GLM      ProviderConfig  // 智谱GLM配置
    MiniMax  ProviderConfig  // MiniMax配置
    Claude   ProviderConfig  // Anthropic Claude配置
    Gemini   ProviderConfig  // Google Gemini配置
    Grok     ProviderConfig  // xAI Grok配置
    Deepseek ProviderConfig  // Deepseek配置
    Ollama   ProviderConfig  // Ollama本地模型配置
}

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
```

### 环境变量配置

```bash
export GPT_API_KEY="your-openai-api-key"
export QWEN_API_KEY="your-qwen-api-key"
export KIMI_API_KEY="your-kimi-api-key"
export GLM_API_KEY="your-glm-api-key"
export DEEPSEEK_API_KEY="your-deepseek-api-key"
```

```go
// 从环境变量自动加载配置
config := bridge.ConfigFromEnv()
sdk := bridge.NewSDK(config)
```

## 高级用法

### 系统提示词模板

支持 `{{question}}` 宏定义：

```go
template := `你是一位Go语言专家，请详细解释：{{question}}

要求：
1. 提供代码示例
2. 说明使用场景
3. 指出常见错误`

resp, _ := client.Generate(ctx, "什么是goroutine", bridge.WithSystemPrompt(template))
```

### 流式响应

```go
// 流式生成
stream, _ := client.GenerateStream(context.Background(), "讲个故事")
defer stream.Close()

for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    fmt.Print(msg.Content)  // 实时输出
}
```

### MCP工具使用

```go
import (
    "ai-bridge/pkg/mcp"
    "github.com/cloudwego/eino/components/tool"
)

// 创建工具
weatherTool := mcp.NewMCPTool(
    "get_weather",
    "查询指定城市的天气",
    mcp.CreateParameterSchema(
        mcp.CreateStringProperty("city", "城市名称", true),
    ),
    func(ctx context.Context, params map[string]interface{}) (string, error) {
        city := params["city"].(string)
        return fmt.Sprintf("%s的天气：晴天，25°C", city), nil
    },
)

// 创建带工具的客户端
client, _ := sdk.CreateClientWithTools(
    types.ProviderGPT,
    "gpt-4",
    []tool.BaseTool{weatherTool.ToEinoTool()},
)
```

## Agent Skills

Agent Skills 是 Anthropic 提出的开放标准，用于给 AI Agent 提供可复用的能力。

### Skill 文件夹结构

```
skills/
├── calculator/
│   └── SKILL.md
├── weather/
│   └── SKILL.md
└── code-search/
    └── SKILL.md
```

### SKILL.md 格式

```markdown
---
name: calculator
description: 执行数学计算，支持基本运算（加减乘除）、幂运算和开方
version: "1.0"
metadata:
  author: ai-bridge
  tags: math, calculation
---

# Calculator Skill

## 概述
这是一个数学计算技能，帮助 AI 助手执行各种数学运算。

## 使用场景
- 基本运算：加减乘除
- 高级运算：幂运算、开方

## 示例
**用户**：帮我计算 123 + 456
**AI**：123 + 456 = 579
```

### 加载和使用 Skills

```go
import "ai-bridge/pkg/skills"

// 加载单个 Skill
skill, _ := skills.LoadSkill("skills/calculator")
systemPrompt := skill.GetSystemPrompt()

// 从目录加载所有 Skills
registry := skills.NewRegistry()
registry.LoadFromDir("skills")

// 获取特定 Skill
calcSkill, ok := registry.Get("calculator")
```

## Docker 运行

### 使用 Docker

```bash
# 拉取镜像
docker pull codyrao/ai-bridge:latest

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e GPT_API_KEY="your-openai-api-key" \
  -e QWEN_API_KEY="your-qwen-api-key" \
  --name ai-bridge \
  codyrao/ai-bridge:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  ai-bridge:
    image: codyrao/ai-bridge:latest
    ports:
      - "8080:8080"
    environment:
      - GPT_API_KEY=${GPT_API_KEY}
      - QWEN_API_KEY=${QWEN_API_KEY}
      - DEEPSEEK_API_KEY=${DEEPSEEK_API_KEY}
    restart: unless-stopped
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
