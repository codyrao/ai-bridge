package types

import (
	"testing"
)

func TestGetModelsByProvider(t *testing.T) {
	// 测试获取GPT模型
	models := GetModelsByProvider(ProviderGPT)
	if len(models) == 0 {
		t.Error("GetModelsByProvider(ProviderGPT) returned empty list")
	}

	// 测试获取不存在的厂商
	models = GetModelsByProvider("nonexistent")
	if models != nil {
		t.Error("GetModelsByProvider(nonexistent) should return nil")
	}
}

func TestGetModelInfo(t *testing.T) {
	// 测试获取存在的模型
	info := GetModelInfo(ProviderGPT, "gpt-4")
	if info == nil {
		t.Fatal("GetModelInfo(ProviderGPT, gpt-4) returned nil")
	}
	if info.Name != "gpt-4" {
		t.Errorf("Expected name 'gpt-4', got '%s'", info.Name)
	}
	if info.Provider != ProviderGPT {
		t.Errorf("Expected provider 'gpt', got '%s'", info.Provider)
	}

	// 测试获取不存在的模型
	info = GetModelInfo(ProviderGPT, "nonexistent")
	if info != nil {
		t.Error("GetModelInfo should return nil for nonexistent model")
	}

	// 测试获取不存在的厂商
	info = GetModelInfo("nonexistent", "model")
	if info != nil {
		t.Error("GetModelInfo should return nil for nonexistent provider")
	}
}

func TestIsValidProvider(t *testing.T) {
	tests := []struct {
		provider Provider
		want     bool
	}{
		{ProviderGPT, true},
		{ProviderQWen, true},
		{ProviderKimi, true},
		{ProviderGLM, true},
		{ProviderMiniMax, true},
		{ProviderClaude, true},
		{ProviderGemini, true},
		{ProviderGrok, true},
		{ProviderDeepseek, true},
		{ProviderOllama, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.provider), func(t *testing.T) {
			got := IsValidProvider(tt.provider)
			if got != tt.want {
				t.Errorf("IsValidProvider(%s) = %v, want %v", tt.provider, got, tt.want)
			}
		})
	}
}

func TestIsValidModel(t *testing.T) {
	tests := []struct {
		provider  Provider
		modelName string
		want      bool
	}{
		{ProviderGPT, "gpt-4", true},
		{ProviderGPT, "gpt-3.5-turbo", true},
		{ProviderGPT, "gpt-4o", true},
		{ProviderGPT, "invalid-model", false},
		{ProviderQWen, "qwen-turbo", true},
		{ProviderQWen, "qwen-plus", true},
		{ProviderDeepseek, "deepseek-chat", true},
		{ProviderDeepseek, "deepseek-coder", true},
		{ProviderKimi, "moonshot-v1-8k", true},
		{ProviderGLM, "glm-4", true},
		{"invalid", "model", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.provider)+"/"+tt.modelName, func(t *testing.T) {
			got := IsValidModel(tt.provider, tt.modelName)
			if got != tt.want {
				t.Errorf("IsValidModel(%s, %s) = %v, want %v", tt.provider, tt.modelName, got, tt.want)
			}
		})
	}
}

func TestModelRegistry(t *testing.T) {
	// 测试所有厂商都有模型
	for provider := range ModelRegistry {
		models := ModelRegistry[provider]
		if len(models) == 0 {
			t.Errorf("Provider %s has no models", provider)
		}

		// 测试每个模型都有正确的厂商
		for _, model := range models {
			if model.Provider != provider {
				t.Errorf("Model %s has provider %s, expected %s", model.Name, model.Provider, provider)
			}
		}
	}
}
