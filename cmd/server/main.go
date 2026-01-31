package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/eino/schema"

	"ai-bridge/pkg/bridge"
	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

// ChatRequest 聊天请求
type ChatRequest struct {
	Provider  string            `json:"provider"`
	Model     string            `json:"model"`
	Messages  []Message         `json:"messages"`
	APIKey    string            `json:"api_key,omitempty"`
	Stream    bool              `json:"stream,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content      string    `json:"content,omitempty"`
	Error        string    `json:"error,omitempty"`
	Provider     string    `json:"provider"`
	Model        string    `json:"model"`
	Stream       bool      `json:"stream"`
}

// ProvidersResponse 厂商列表响应
type ProvidersResponse struct {
	Providers []ProviderInfo `json:"providers"`
}

// ProviderInfo 厂商信息
type ProviderInfo struct {
	Name   string       `json:"name"`
	Models []ModelInfo  `json:"models"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxTokens   int    `json:"max_tokens"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/providers", providersHandler)
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/chat/stream", chatStreamHandler)

	log.Printf("AI Bridge Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func providersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	providers := bridge.GetProviders()
	var result ProvidersResponse

	for _, provider := range providers {
		models := bridge.GetModels(provider)
		var modelInfos []ModelInfo

		for _, m := range models {
			modelInfos = append(modelInfos, ModelInfo{
				Name:        m.Name,
				Description: m.Description,
				MaxTokens:   m.MaxTokens,
			})
		}

		result.Providers = append(result.Providers, ProviderInfo{
			Name:   string(provider),
			Models: modelInfos,
		})
	}

	json.NewEncoder(w).Encode(result)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证必需参数
	if req.Provider == "" || req.Model == "" || len(req.Messages) == 0 {
		http.Error(w, "Missing required parameters: provider, model, messages", http.StatusBadRequest)
		return
	}

	// 获取API Key（优先从请求体，其次从环境变量）
	apiKey := req.APIKey
	if apiKey == "" {
		apiKey = os.Getenv(fmt.Sprintf("%s_API_KEY", req.Provider))
	}

	if apiKey == "" {
		http.Error(w, "API key not provided", http.StatusUnauthorized)
		return
	}

	// 构建选项
	opts := []options.Option{
		options.WithAPIKey(apiKey),
	}

	// 应用额外选项
	if temp, ok := req.Options["temperature"].(float64); ok {
		opts = append(opts, options.WithTemperature(float32(temp)))
	}
	if maxTokens, ok := req.Options["max_tokens"].(float64); ok {
		opts = append(opts, options.WithMaxTokens(int(maxTokens)))
	}
	if topP, ok := req.Options["top_p"].(float64); ok {
		opts = append(opts, options.WithTopP(float32(topP)))
	}

	// 创建客户端
	client, err := bridge.NewAIClient(types.Provider(req.Provider), req.Model, opts...)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Error:    err.Error(),
			Provider: req.Provider,
			Model:    req.Model,
		})
		return
	}

	// 转换消息格式
	var messages []*schema.Message
	for _, msg := range req.Messages {
		var role schema.RoleType
		switch msg.Role {
		case "user":
			role = schema.User
		case "assistant":
			role = schema.Assistant
		case "system":
			role = schema.System
		default:
			role = schema.User
		}
		messages = append(messages, &schema.Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 执行对话
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.Chat(ctx, messages)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Error:    err.Error(),
			Provider: req.Provider,
			Model:    req.Model,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{
		Content:  resp.Content,
		Provider: req.Provider,
		Model:    req.Model,
		Stream:   false,
	})
}

func chatStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 验证必需参数
	if req.Provider == "" || req.Model == "" || len(req.Messages) == 0 {
		http.Error(w, "Missing required parameters: provider, model, messages", http.StatusBadRequest)
		return
	}

	// 获取API Key
	apiKey := req.APIKey
	if apiKey == "" {
		apiKey = os.Getenv(fmt.Sprintf("%s_API_KEY", req.Provider))
	}

	if apiKey == "" {
		http.Error(w, "API key not provided", http.StatusUnauthorized)
		return
	}

	// 构建选项
	opts := []options.Option{
		options.WithAPIKey(apiKey),
	}

	// 创建客户端
	client, err := bridge.NewAIClient(types.Provider(req.Provider), req.Model, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 转换消息格式
	var messages []*schema.Message
	for _, msg := range req.Messages {
		var role schema.RoleType
		switch msg.Role {
		case "user":
			role = schema.User
		case "assistant":
			role = schema.Assistant
		case "system":
			role = schema.System
		default:
			role = schema.User
		}
		messages = append(messages, &schema.Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 设置SSE头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 执行流式对话
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stream, err := client.ChatStream(ctx, messages)
	if err != nil {
		fmt.Fprintf(w, "data: %s\n\n", `{"error": "`+err.Error()+`"}`)
		return
	}
	defer stream.Close()

	// 发送流式响应
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Fprintf(w, "data: %s\n\n", `{"done": true}`)
				w.(http.Flusher).Flush()
			}
			break
		}

		data, _ := json.Marshal(map[string]string{
			"content": msg.Content,
		})
		fmt.Fprintf(w, "data: %s\n\n", string(data))
		w.(http.Flusher).Flush()
	}
}
