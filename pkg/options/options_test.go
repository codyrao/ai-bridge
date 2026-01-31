package options

import (
	"testing"
	"time"
)

func TestApplyOptions(t *testing.T) {
	// 测试默认配置
	config := ApplyOptions()
	if config.Timeout != 120*time.Second {
		t.Errorf("Expected default timeout 120s, got %v", config.Timeout)
	}
	if config.Temperature != 0.7 {
		t.Errorf("Expected default temperature 0.7, got %f", config.Temperature)
	}
	if config.MaxTokens != 2048 {
		t.Errorf("Expected default max_tokens 2048, got %d", config.MaxTokens)
	}
}

func TestWithAPIKey(t *testing.T) {
	config := ApplyOptions(WithAPIKey("test-api-key"))
	if config.APIKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", config.APIKey)
	}
}

func TestWithBaseURL(t *testing.T) {
	config := ApplyOptions(WithBaseURL("https://api.example.com"))
	if config.BaseURL != "https://api.example.com" {
		t.Errorf("Expected base URL 'https://api.example.com', got '%s'", config.BaseURL)
	}
}

func TestWithTimeout(t *testing.T) {
	config := ApplyOptions(WithTimeout(30 * time.Second))
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", config.Timeout)
	}
}

func TestWithMaxRetries(t *testing.T) {
	config := ApplyOptions(WithMaxRetries(5))
	if config.MaxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", config.MaxRetries)
	}
}

func TestWithTemperature(t *testing.T) {
	config := ApplyOptions(WithTemperature(0.5))
	if config.Temperature != 0.5 {
		t.Errorf("Expected temperature 0.5, got %f", config.Temperature)
	}
}

func TestWithTopP(t *testing.T) {
	config := ApplyOptions(WithTopP(0.9))
	if config.TopP != 0.9 {
		t.Errorf("Expected top_p 0.9, got %f", config.TopP)
	}
}

func TestWithMaxTokens(t *testing.T) {
	config := ApplyOptions(WithMaxTokens(4096))
	if config.MaxTokens != 4096 {
		t.Errorf("Expected max tokens 4096, got %d", config.MaxTokens)
	}
}

func TestWithStream(t *testing.T) {
	config := ApplyOptions(WithStream(true))
	if !config.Stream {
		t.Error("Expected stream to be true")
	}
}

func TestWithExtraHeaders(t *testing.T) {
	headers := map[string]string{
		"X-Custom-Header": "value",
		"X-Another":       "another-value",
	}
	config := ApplyOptions(WithExtraHeaders(headers))
	if config.ExtraHeaders["X-Custom-Header"] != "value" {
		t.Error("Extra header not set correctly")
	}
	if config.ExtraHeaders["X-Another"] != "another-value" {
		t.Error("Extra header not set correctly")
	}
}

func TestWithExtraHeader(t *testing.T) {
	config := ApplyOptions(WithExtraHeader("X-Single", "single-value"))
	if config.ExtraHeaders["X-Single"] != "single-value" {
		t.Errorf("Expected extra header 'X-Single' to be 'single-value', got '%s'", config.ExtraHeaders["X-Single"])
	}
}

func TestWithOrganization(t *testing.T) {
	config := ApplyOptions(WithOrganization("org-test"))
	if config.Organization != "org-test" {
		t.Errorf("Expected organization 'org-test', got '%s'", config.Organization)
	}
}

func TestWithProxy(t *testing.T) {
	config := ApplyOptions(WithProxy("http://proxy.example.com:8080"))
	if config.Proxy != "http://proxy.example.com:8080" {
		t.Errorf("Expected proxy 'http://proxy.example.com:8080', got '%s'", config.Proxy)
	}
}

func TestWithEnableLog(t *testing.T) {
	config := ApplyOptions(WithEnableLog(true))
	if !config.EnableLog {
		t.Error("Expected enable_log to be true")
	}
}

func TestMultipleOptions(t *testing.T) {
	config := ApplyOptions(
		WithAPIKey("my-key"),
		WithTemperature(0.9),
		WithMaxTokens(1024),
		WithTimeout(60*time.Second),
		WithStream(true),
	)

	if config.APIKey != "my-key" {
		t.Errorf("Expected API key 'my-key', got '%s'", config.APIKey)
	}
	if config.Temperature != 0.9 {
		t.Errorf("Expected temperature 0.9, got %f", config.Temperature)
	}
	if config.MaxTokens != 1024 {
		t.Errorf("Expected max tokens 1024, got %d", config.MaxTokens)
	}
	if config.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", config.Timeout)
	}
	if !config.Stream {
		t.Error("Expected stream to be true")
	}
}
