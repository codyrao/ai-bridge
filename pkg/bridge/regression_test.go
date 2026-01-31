package bridge

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"

	"github.com/cloudwego/eino/schema"
)

// RegressionTestResult 回归测试结果
type RegressionTestResult struct {
	Provider    string
	Model       string
	Status      string // PASS / FAIL / SKIP
	Error       string
	Latency     time.Duration
	TestsPassed int
	TestsTotal  int
}

// TestAllProvidersRegression 全量回归测试所有平台
func TestAllProvidersRegression(t *testing.T) {
	fmt.Println("========================================")
	fmt.Println("     AI Bridge 全量回归测试")
	fmt.Println("========================================")
	fmt.Println()

	var results []RegressionTestResult

	// 定义所有支持的厂商和模型
	testCases := []struct {
		provider types.Provider
		model    string
		envKey   string
		baseURL  string
	}{
		{types.ProviderGPT, "gpt-3.5-turbo", "GPT_API_KEY", ""},
		{types.ProviderQWen, "qwen-turbo", "QWEN_API_KEY", ""},
		{types.ProviderKimi, "moonshot-v1-8k", "KIMI_API_KEY", ""},
		{types.ProviderGLM, "glm-4", "GLM_API_KEY", ""},
		{types.ProviderMiniMax, "abab6.5s-chat", "MINIMAX_API_KEY", ""},
		{types.ProviderClaude, "claude-3-haiku-20240307", "CLAUDE_API_KEY", ""},
		{types.ProviderGemini, "gemini-1.5-flash", "GEMINI_API_KEY", ""},
		{types.ProviderGrok, "grok-1", "GROK_API_KEY", ""},
		{types.ProviderDeepseek, "deepseek-chat", "DEEPSEEK_API_KEY", ""},
		{types.ProviderOllama, "qwen3-coder:30b", "", "http://localhost:11434"},
	}

	for _, tc := range testCases {
		result := runProviderRegressionTest(t, tc.provider, tc.model, tc.envKey, tc.baseURL)
		results = append(results, result)
	}

	// 打印测试报告
	printRegressionReport(results)
}

// runProviderRegressionTest 运行单个厂商的回归测试
func runProviderRegressionTest(t *testing.T, provider types.Provider, modelName, envKey, baseURL string) RegressionTestResult {
	result := RegressionTestResult{
		Provider:   string(provider),
		Model:      modelName,
		Status:     "SKIP",
		TestsTotal: 4,
	}

	fmt.Printf("【%s - %s】\n", provider, modelName)

	// 检查API Key（Ollama除外）
	apiKey := ""
	if envKey != "" {
		apiKey = os.Getenv(envKey)
		if apiKey == "" {
			result.Error = fmt.Sprintf("环境变量 %s 未设置", envKey)
			fmt.Printf("  ⚠️  跳过: %s\n\n", result.Error)
			return result
		}
	}

	// 构建选项
	opts := []options.Option{
		options.WithAPIKey(apiKey),
		options.WithTemperature(0.7),
		options.WithMaxTokens(256),
		options.WithTimeout(60 * time.Second),
	}
	if baseURL != "" {
		opts = append(opts, options.WithBaseURL(baseURL))
	}

	// 创建客户端
	start := time.Now()
	client, err := NewAIClient(provider, modelName, opts...)
	result.Latency = time.Since(start)

	if err != nil {
		result.Status = "FAIL"
		result.Error = fmt.Sprintf("创建客户端失败: %v", err)
		fmt.Printf("  ❌ 失败: %s\n\n", result.Error)
		return result
	}

	fmt.Printf("  ✓ 客户端创建成功 (%.2fs)\n", result.Latency.Seconds())

	// 测试1: 基础对话
	test1Pass := testBasicChat(t, client, &result)

	// 测试2: 流式对话
	test2Pass := testStreamChat(t, client, &result)

	// 测试3: 多轮对话
	test3Pass := testMultiTurnChat(t, client, &result)

	// 测试4: 模型信息
	test4Pass := testModelInfo(t, client, &result)

	// 统计结果
	if test1Pass && test2Pass && test3Pass && test4Pass {
		result.Status = "PASS"
		result.TestsPassed = 4
		fmt.Printf("  ✅ 全部通过\n\n")
	} else {
		result.Status = "FAIL"
		result.TestsPassed = countPassedTests(test1Pass, test2Pass, test3Pass, test4Pass)
		fmt.Printf("  ❌ 部分失败 (%d/4)\n\n", result.TestsPassed)
	}

	return result
}

// testBasicChat 测试基础对话
func testBasicChat(t *testing.T, client types.AIBridge, result *RegressionTestResult) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messages := []*schema.Message{
		schema.UserMessage("你好，请用一句话介绍自己"),
	}

	resp, err := client.Chat(ctx, messages)
	if err != nil {
		fmt.Printf("  ✗ 基础对话失败: %v\n", err)
		return false
	}

	if resp.Content == "" {
		fmt.Printf("  ✗ 基础对话返回空内容\n")
		return false
	}

	fmt.Printf("  ✓ 基础对话通过 (%d字符)\n", len(resp.Content))
	return true
}

// testStreamChat 测试流式对话
func testStreamChat(t *testing.T, client types.AIBridge, result *RegressionTestResult) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messages := []*schema.Message{
		schema.UserMessage("Hi"),
	}

	stream, err := client.ChatStream(ctx, messages)
	if err != nil {
		fmt.Printf("  ✗ 流式对话失败: %v\n", err)
		return false
	}
	defer stream.Close()

	chunkCount := 0
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
		chunkCount++
	}

	if chunkCount == 0 {
		fmt.Printf("  ✗ 流式对话未返回数据\n")
		return false
	}

	fmt.Printf("  ✓ 流式对话通过 (%d个数据块)\n", chunkCount)
	return true
}

// testMultiTurnChat 测试多轮对话
func testMultiTurnChat(t *testing.T, client types.AIBridge, result *RegressionTestResult) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messages := []*schema.Message{
		schema.UserMessage("什么是Go语言？"),
	}

	resp1, err := client.Chat(ctx, messages)
	if err != nil {
		fmt.Printf("  ✗ 多轮对话第一轮失败: %v\n", err)
		return false
	}

	messages = append(messages, resp1)
	messages = append(messages, schema.UserMessage("它有什么特点？"))

	_, err = client.Chat(ctx, messages)
	if err != nil {
		fmt.Printf("  ✗ 多轮对话第二轮失败: %v\n", err)
		return false
	}

	fmt.Printf("  ✓ 多轮对话通过\n")
	return true
}

// testModelInfo 测试模型信息
func testModelInfo(t *testing.T, client types.AIBridge, result *RegressionTestResult) bool {
	info := client.GetModelInfo()
	if info == nil {
		fmt.Printf("  ✗ 模型信息为空\n")
		return false
	}

	if info.Name == "" || info.Provider == "" {
		fmt.Printf("  ✗ 模型信息不完整\n")
		return false
	}

	fmt.Printf("  ✓ 模型信息通过 (%s)\n", info.Name)
	return true
}

// countPassedTests 统计通过的测试数
func countPassedTests(tests ...bool) int {
	count := 0
	for _, t := range tests {
		if t {
			count++
		}
	}
	return count
}

// printRegressionReport 打印回归测试报告
func printRegressionReport(results []RegressionTestResult) {
	fmt.Println("========================================")
	fmt.Println("         回归测试报告")
	fmt.Println("========================================")

	passCount := 0
	failCount := 0
	skipCount := 0

	for _, r := range results {
		switch r.Status {
		case "PASS":
			passCount++
		case "FAIL":
			failCount++
		case "SKIP":
			skipCount++
		}
	}

	fmt.Printf("\n总计: %d 个厂商\n", len(results))
	fmt.Printf("  ✅ 通过: %d\n", passCount)
	fmt.Printf("  ❌ 失败: %d\n", failCount)
	fmt.Printf("  ⚠️  跳过: %d\n\n", skipCount)

	if failCount > 0 {
		fmt.Println("失败详情:")
		for _, r := range results {
			if r.Status == "FAIL" {
				fmt.Printf("  - %s (%s): %s\n", r.Provider, r.Model, r.Error)
			}
		}
		fmt.Println()
	}

	if skipCount > 0 {
		fmt.Println("跳过详情:")
		for _, r := range results {
			if r.Status == "SKIP" {
				fmt.Printf("  - %s (%s): %s\n", r.Provider, r.Model, r.Error)
			}
		}
		fmt.Println()
	}

	fmt.Println("========================================")
}

// TestProviderRegistry 测试厂商注册表
func TestProviderRegistry(t *testing.T) {
	fmt.Println("\n【厂商注册表测试】")

	providers := GetProviders()
	fmt.Printf("  注册厂商数量: %d\n", len(providers))

	expectedProviders := []types.Provider{
		types.ProviderGPT,
		types.ProviderQWen,
		types.ProviderKimi,
		types.ProviderGLM,
		types.ProviderMiniMax,
		types.ProviderClaude,
		types.ProviderGemini,
		types.ProviderGrok,
		types.ProviderDeepseek,
		types.ProviderOllama,
	}

	providerMap := make(map[types.Provider]bool)
	for _, p := range providers {
		providerMap[p] = true
	}

	allFound := true
	for _, expected := range expectedProviders {
		if !providerMap[expected] {
			fmt.Printf("  ✗ 缺少厂商: %s\n", expected)
			allFound = false
		}
	}

	if allFound {
		fmt.Printf("  ✓ 所有厂商已注册\n")
	}

	// 测试每个厂商的模型列表
	for _, provider := range expectedProviders {
		models := GetModels(provider)
		if len(models) == 0 && provider != types.ProviderOllama {
			fmt.Printf("  ✗ %s 没有注册模型\n", provider)
		} else {
			fmt.Printf("  ✓ %s: %d 个模型\n", provider, len(models))
		}
	}
}

// TestAdapterFactory 测试适配器工厂
func TestAdapterFactory(t *testing.T) {
	fmt.Println("\n【适配器工厂测试】")

	// 测试无效的厂商
	_, err := NewAIClient("invalid", "model", options.WithAPIKey("test"))
	if err == nil {
		t.Error("应该拒绝无效厂商")
	} else {
		fmt.Printf("  ✓ 正确拒绝无效厂商\n")
	}

	// 测试无效的模型
	_, err = NewAIClient(types.ProviderGPT, "invalid-model", options.WithAPIKey("test"))
	if err == nil {
		t.Error("应该拒绝无效模型")
	} else {
		fmt.Printf("  ✓ 正确拒绝无效模型\n")
	}

	// 测试缺少API Key
	_, err = NewAIClient(types.ProviderGPT, "gpt-3.5-turbo")
	if err == nil {
		t.Error("应该拒绝缺少API Key")
	} else {
		fmt.Printf("  ✓ 正确拒绝缺少API Key\n")
	}
}

// TestOptionPropagation 测试配置选项传递
func TestOptionPropagation(t *testing.T) {
	fmt.Println("\n【配置选项传递测试】")

	// 测试所有配置选项
	opts := []options.Option{
		options.WithAPIKey("test-key"),
		options.WithBaseURL("https://api.example.com"),
		options.WithTimeout(30 * time.Second),
		options.WithMaxRetries(5),
		options.WithTemperature(0.5),
		options.WithTopP(0.8),
		options.WithMaxTokens(1024),
		options.WithStream(true),
		options.WithExtraHeader("X-Custom", "value"),
		options.WithOrganization("org-test"),
		options.WithProxy("http://proxy.example.com"),
		options.WithEnableLog(true),
	}

	config := options.ApplyOptions(opts...)

	checks := []struct {
		name  string
		check bool
	}{
		{"APIKey", config.APIKey == "test-key"},
		{"BaseURL", config.BaseURL == "https://api.example.com"},
		{"Timeout", config.Timeout == 30*time.Second},
		{"MaxRetries", config.MaxRetries == 5},
		{"Temperature", config.Temperature == 0.5},
		{"TopP", config.TopP == 0.8},
		{"MaxTokens", config.MaxTokens == 1024},
		{"Stream", config.Stream == true},
		{"Organization", config.Organization == "org-test"},
		{"Proxy", config.Proxy == "http://proxy.example.com"},
		{"EnableLog", config.EnableLog == true},
		{"ExtraHeaders", config.ExtraHeaders["X-Custom"] == "value"},
	}

	allPass := true
	for _, c := range checks {
		if c.check {
			fmt.Printf("  ✓ %s 正确传递\n", c.name)
		} else {
			fmt.Printf("  ✗ %s 传递失败\n", c.name)
			allPass = false
		}
	}

	if !allPass {
		t.Error("部分配置选项未正确传递")
	}
}
