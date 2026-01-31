package bridge

import (
	"testing"

	"ai-bridge/pkg/types"
)

func TestGetProviders(t *testing.T) {
	providers := GetProviders()
	if len(providers) == 0 {
		t.Error("GetProviders() returned empty list")
	}

	// 检查是否包含主要厂商
	expectedProviders := []types.Provider{
		types.ProviderQWen,
		types.ProviderKimi,
		types.ProviderGLM,
		types.ProviderGPT,
		types.ProviderClaude,
		types.ProviderGemini,
		types.ProviderDeepseek,
		types.ProviderOllama,
	}

	providerMap := make(map[types.Provider]bool)
	for _, p := range providers {
		providerMap[p] = true
	}

	for _, expected := range expectedProviders {
		if !providerMap[expected] {
			t.Errorf("Expected provider %s not found", expected)
		}
	}
}

func TestGetModels(t *testing.T) {
	// 测试获取GPT模型
	gptModels := GetModels(types.ProviderGPT)
	if len(gptModels) == 0 {
		t.Error("GetModels(ProviderGPT) returned empty list")
	}

	// 检查是否包含GPT-4
	foundGPT4 := false
	for _, m := range gptModels {
		if m.Name == "gpt-4" {
			foundGPT4 = true
			break
		}
	}
	if !foundGPT4 {
		t.Error("GPT-4 model not found")
	}

	// 测试获取不存在的厂商
	models := GetModels("nonexistent")
	if models != nil {
		t.Error("GetModels(nonexistent) should return nil")
	}
}

func TestIsValidProvider(t *testing.T) {
	tests := []struct {
		provider types.Provider
		want     bool
	}{
		{types.ProviderGPT, true},
		{types.ProviderQWen, true},
		{types.ProviderDeepseek, true},
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
		provider  types.Provider
		modelName string
		want      bool
	}{
		{types.ProviderGPT, "gpt-4", true},
		{types.ProviderGPT, "gpt-3.5-turbo", true},
		{types.ProviderGPT, "invalid-model", false},
		{types.ProviderQWen, "qwen-turbo", true},
		{types.ProviderDeepseek, "deepseek-chat", true},
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

func TestGetModelInfo(t *testing.T) {
	// 测试获取存在的模型信息
	info := GetModelInfo(types.ProviderGPT, "gpt-4")
	if info == nil {
		t.Fatal("GetModelInfo(ProviderGPT, gpt-4) returned nil")
	}
	if info.Name != "gpt-4" {
		t.Errorf("Expected model name 'gpt-4', got '%s'", info.Name)
	}
	if info.Provider != types.ProviderGPT {
		t.Errorf("Expected provider 'gpt', got '%s'", info.Provider)
	}

	// 测试获取不存在的模型
	info = GetModelInfo(types.ProviderGPT, "nonexistent")
	if info != nil {
		t.Error("GetModelInfo should return nil for nonexistent model")
	}
}
