# AI Bridge SDK - Go语言AI大模型聚合SDK

AI Bridge SDK 是一个基于 Go 语言和字节跳动 Eino 框架实现的 AI 大模型接口聚合SDK。它提供了统一的接口来调用市面上主流的 AI 大模型，包括 QWen、Kimi、GLM、MiniMax、Claude、GPT、Gemini、Grok、Deepseek 以及 Ollama 本地模型。

## 功能特性

- **多厂商支持**：支持 10+ 主流 AI 模型厂商
- **统一接口**：提供一致的 SDK 接口，简化多模型切换
- **Config配置模式**：通过Config结构统一配置各模型API和密钥
- **Option配置模式**：丰富的配置选项，灵活控制模型行为
- **MCP 工具支持**：支持 Model Control Plane 工具调用
- **流式响应**：支持流式和非流式两种对话模式
- **RESTful API**：提供 HTTP API 服务，方便集成

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
        // GPT配置
        GPT: bridge.ProviderConfig{
            APIKey:  "your-openai-api-key",
            BaseURL: "", // 可选，使用默认地址
        },
        // 通义千问配置
        QWen: bridge.ProviderConfig{
            APIKey: "your-qwen-api-key",
        },
        // Deepseek配置
        Deepseek: bridge.ProviderConfig{
            APIKey: "your-deepseek-api-key",
        },
        // Ollama本地模型配置
        Ollama: bridge.ProviderConfig{
            BaseURL: "http://localhost:11434",
        },
    }

    // 创建SDK客户端
    sdk := bridge.NewSDK(config)

    // 使用GPT模型
    gptClient, err := sdk.CreateClient(types.ProviderGPT, "gpt-3.5-turbo")
    if err != nil {
        log.Fatal(err)
    }

    // 执行对话
    resp, err := gptClient.Generate(context.Background(), "你好，请介绍一下自己")
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
// SDKConfig SDK全局配置
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
```

### 完整配置示例

```go
config := &bridge.SDKConfig{
    GPT: bridge.ProviderConfig{
        APIKey:      os.Getenv("GPT_API_KEY"),
        BaseURL:     "",                       // 使用默认地址
        Timeout:     120 * time.Second,
        MaxRetries:  3,
        Temperature: 0.7,
        TopP:        0.9,
        MaxTokens:   4096,
        Proxy:       "http://proxy.example.com:8080",
    },
    QWen: bridge.ProviderConfig{
        APIKey:      os.Getenv("QWEN_API_KEY"),
        Temperature: 0.8,
        MaxTokens:   2048,
    },
    Kimi: bridge.ProviderConfig{
        APIKey: os.Getenv("KIMI_API_KEY"),
    },
    GLM: bridge.ProviderConfig{
        APIKey: os.Getenv("GLM_API_KEY"),
    },
    MiniMax: bridge.ProviderConfig{
        APIKey: os.Getenv("MINIMAX_API_KEY"),
    },
    Claude: bridge.ProviderConfig{
        APIKey: os.Getenv("CLAUDE_API_KEY"),
    },
    Gemini: bridge.ProviderConfig{
        APIKey: os.Getenv("GEMINI_API_KEY"),
    },
    Grok: bridge.ProviderConfig{
        APIKey: os.Getenv("GROK_API_KEY"),
    },
    Deepseek: bridge.ProviderConfig{
        APIKey: os.Getenv("DEEPSEEK_API_KEY"),
    },
    Ollama: bridge.ProviderConfig{
        BaseURL: "http://localhost:11434",  // Ollama默认地址
    },
}

sdk := bridge.NewSDK(config)
```

## SDK使用示例

### 使用高级客户端（推荐）

高级客户端支持更多Option配置：

```go
// 创建高级客户端
client, _ := sdk.CreateSDKClient(types.ProviderGPT, "gpt-3.5-turbo")

// 1. 简单生成（默认启用流式，60秒超时）
resp, _ := client.Generate(context.Background(), "你好")
fmt.Println(resp)

// 2. 带历史记录的生成
history := []*schema.Message{
    schema.SystemMessage("你是一个专业助手。"),
    schema.UserMessage("什么是Go语言？"),
    schema.AssistantMessage("Go是一种编程语言。"),
}
resp, _ = client.Generate(context.Background(), 
    "它有什么特点？",
    bridge.WithHistory(history),
)

// 3. 禁用流式，设置超时
resp, _ = client.Generate(context.Background(), 
    "你好",
    bridge.WithStream(false),
    bridge.WithTimeout(30*time.Second),
)

// 4. 流式生成（实时输出）
stream, _ := client.GenerateStream(context.Background(), "你好")
defer stream.Close()

for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    fmt.Print(msg.Content)  // 实时输出
}
```

### 对话模式

```go
client, _ := sdk.CreateSDKClient(types.ProviderQWen, "qwen-turbo")

messages := []*schema.Message{
    schema.SystemMessage("你是一个专业的编程助手。"),
    schema.UserMessage("什么是Go语言？"),
}

// 1. 流式对话（默认）
result, _ := client.Chat(context.Background(), messages)
fmt.Println(result.Content)

// 2. 非流式对话
result, _ := client.Chat(context.Background(), messages, bridge.WithStream(false))
fmt.Println(result.Content)

// 3. 设置超时
result, _ = client.Chat(context.Background(), messages, bridge.WithTimeout(120*time.Second))

// 4. 纯流式对话（手动处理流）
stream, _ := client.ChatStream(context.Background(), messages)
defer stream.Close()

for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    fmt.Print(msg.Content)  // 实时输出
}
```

### 高级客户端Option说明

高级客户端支持以下Option：

```go
// WithHistory 设置对话历史
history := []*schema.Message{
    schema.UserMessage("你好"),
    schema.AssistantMessage("你好！有什么可以帮助你的？"),
}
resp, _ := client.Generate(ctx, "今天天气如何？", bridge.WithHistory(history))

// WithStream 设置是否启用流式（默认true）
resp, _ := client.Generate(ctx, "你好", bridge.WithStream(false))  // 禁用流式
resp, _ := client.Generate(ctx, "你好", bridge.WithStream(true))   // 启用流式（默认）

// WithTimeout 设置超时时间（默认60s）
resp, _ := client.Generate(ctx, "你好", bridge.WithTimeout(30*time.Second))
resp, _ := client.Generate(ctx, "复杂问题", bridge.WithTimeout(5*time.Minute))
```

### 方法对比

| 方法 | 流式支持 | 历史支持 | 适用场景 |
|------|---------|---------|---------|
| `Generate()` | ✅ Option控制 | ✅ WithHistory | 简单文本生成 |
| `GenerateStream()` | ✅ 强制流式 | ✅ WithHistory | 需要实时输出的生成 |
| `Chat()` | ✅ Option控制 | ✅ 参数传入 | 多轮对话 |
| `ChatStream()` | ✅ 强制流式 | ✅ 参数传入 | 需要实时输出的对话 |

### 使用MCP工具

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

// 使用工具进行对话
resp, _ := client.Chat(context.Background(), messages)
```

### 切换模型

```go
// 使用不同厂商的模型
providers := []struct {
    provider types.Provider
    model    string
}{
    {types.ProviderGPT, "gpt-4"},
    {types.ProviderQWen, "qwen-max"},
    {types.ProviderDeepseek, "deepseek-coder"},
}

for _, p := range providers {
    client, err := sdk.CreateClient(p.provider, p.model)
    if err != nil {
        log.Printf("创建 %s 客户端失败: %v", p.provider, err)
        continue
    }
    
    resp, err := client.Generate(context.Background(), "你好")
    if err != nil {
        log.Printf("%s 调用失败: %v", p.provider, err)
        continue
    }
    
    fmt.Printf("%s: %s\n", p.provider, resp)
}
```

## 环境变量配置

支持通过环境变量配置API密钥：

```bash
export GPT_API_KEY="your-openai-api-key"
export QWEN_API_KEY="your-qwen-api-key"
export KIMI_API_KEY="your-kimi-api-key"
export GLM_API_KEY="your-glm-api-key"
export MINIMAX_API_KEY="your-minimax-api-key"
export CLAUDE_API_KEY="your-claude-api-key"
export GEMINI_API_KEY="your-gemini-api-key"
export GROK_API_KEY="your-grok-api-key"
export DEEPSEEK_API_KEY="your-deepseek-api-key"
```

```go
// 从环境变量自动加载配置
config := bridge.ConfigFromEnv()
sdk := bridge.NewSDK(config)
```

## RESTful API服务

SDK也提供HTTP API服务：

```go
package main

import (
    "ai-bridge/cmd/server"
)

func main() {
    // 启动HTTP服务
    server.Start(":8080")
}
```

### API端点

- `GET /health` - 健康检查
- `GET /providers` - 获取支持的厂商和模型列表
- `POST /chat` - 非流式对话
- `POST /chat/stream` - 流式对话

## Docker运行

### 使用 Docker 运行

```bash
# 拉取镜像
docker pull codyrao/ai-bridge:latest

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e GPT_API_KEY="your-openai-api-key" \
  -e QWEN_API_KEY="your-qwen-api-key" \
  -e DEEPSEEK_API_KEY="your-deepseek-api-key" \
  --name ai-bridge \
  codyrao/ai-bridge:latest
```

### Docker Compose 部署

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
      - KIMI_API_KEY=${KIMI_API_KEY}
      - GLM_API_KEY=${GLM_API_KEY}
      - CLAUDE_API_KEY=${CLAUDE_API_KEY}
      - GEMINI_API_KEY=${GEMINI_API_KEY}
      - GROK_API_KEY=${GROK_API_KEY}
      - MINIMAX_API_KEY=${MINIMAX_API_KEY}
    restart: unless-stopped
```

## 完整示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/codyrao/ai-bridge/pkg/bridge"
    "github.com/codyrao/ai-bridge/pkg/types"
)

func main() {
    // 方式1：使用Config结构
    config := &bridge.SDKConfig{
        GPT: bridge.ProviderConfig{
            APIKey: os.Getenv("GPT_API_KEY"),
        },
        QWen: bridge.ProviderConfig{
            APIKey: os.Getenv("QWEN_API_KEY"),
        },
        Ollama: bridge.ProviderConfig{
            BaseURL: "http://localhost:11434",
        },
    }

    sdk := bridge.NewSDK(config)

    // 测试不同模型
    testModels := []struct {
        name     string
        provider types.Provider
        model    string
    }{
        {"GPT-3.5", types.ProviderGPT, "gpt-3.5-turbo"},
        {"通义千问", types.ProviderQWen, "qwen-turbo"},
        {"Ollama本地", types.ProviderOllama, "qwen3-coder:30b"},
    }

    for _, tm := range testModels {
        client, err := sdk.CreateClient(tm.provider, tm.model)
        if err != nil {
            log.Printf("[%s] 创建失败: %v", tm.name, err)
            continue
        }

        resp, err := client.Generate(context.Background(), "你好，请用一句话介绍自己")
        if err != nil {
            log.Printf("[%s] 调用失败: %v", tm.name, err)
            continue
        }

        fmt.Printf("[%s] %s\n", tm.name, resp)
    }
}
```

## 项目结构

```
ai-bridge/
├── pkg/
│   ├── bridge/          # SDK入口和客户端创建
│   ├── adapters/        # 模型适配器实现
│   ├── mcp/             # MCP工具支持
│   ├── options/         # Option配置模式
│   └── types/           # 类型定义
├── cmd/
│   └── server/          # HTTP服务
├── examples/            # 使用示例
├── docker-compose.yml   # Docker部署配置
├── Dockerfile           # Docker镜像构建
└── README.md            # 项目说明
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 联系方式

如有问题，请通过 GitHub Issues 联系。
