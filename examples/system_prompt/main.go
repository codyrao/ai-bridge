package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"ai-bridge/pkg/adapters"
	"ai-bridge/pkg/bridge"
	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("     系统提示词和Agent Skills示例")
	fmt.Println("========================================")
	fmt.Println()

	// 示例1: 使用系统提示词配置适配器
	demoSystemPromptAdapter()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println()

	// 示例2: 在SDK客户端调用时指定系统提示词
	demoSystemPromptClient()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println()

	// 示例3: 使用{{question}}模板宏
	demoSystemPromptTemplate()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println()

	// 示例4: 配置Agent Skills
	demoAgentSkills()
}

// demoSystemPromptAdapter 在适配器级别配置系统提示词
func demoSystemPromptAdapter() {
	fmt.Println("【示例1: 适配器级别系统提示词】")
	fmt.Println()

	// 定义系统提示词 - 让AI扮演一个专业的Go语言专家
	systemPrompt := `你是一位资深的Go语言开发专家，拥有10年以上的Go开发经验。
你的职责是：
1. 提供高质量的Go代码示例
2. 解释Go语言的最佳实践
3. 帮助优化代码性能
4. 解答Go语言相关的技术问题
请用专业但易懂的方式回答问题。`

	// 创建Ollama适配器，传入系统提示词
	adapter, err := adapters.GetAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithSystemPrompt(systemPrompt),
		options.WithTemperature(0.7),
		options.WithMaxTokens(1024),
	)
	if err != nil {
		log.Printf("创建适配器失败: %v\n", err)
		return
	}

	fmt.Println("✓ 适配器创建成功（已配置系统提示词）")
	fmt.Println()

	// 测试生成
	ctx := context.Background()
	prompt := "请解释一下Go语言中的channel是什么，并给出一个简单的使用示例。"

	fmt.Printf("用户问题: %s\n", prompt)
	fmt.Println("等待响应...")
	fmt.Println()

	resp, err := adapter.Generate(ctx, prompt)
	if err != nil {
		log.Printf("生成失败: %v\n", err)
		return
	}

	fmt.Println("AI响应:")
	fmt.Println(resp)
}

// demoSystemPromptClient 在SDK客户端调用时指定系统提示词
func demoSystemPromptClient() {
	fmt.Println("【示例2: SDK客户端级别系统提示词】")
	fmt.Println()

	// 创建SDK配置
	config := &bridge.SDKConfig{
		Ollama: bridge.ProviderConfig{
			BaseURL:     "http://localhost:11434",
			Temperature: 0.7,
			MaxTokens:   512,
		},
	}

	sdk := bridge.NewSDK(config)

	// 创建高级客户端
	client, err := sdk.CreateSDKClient(types.ProviderOllama, "qwen3-coder:30b")
	if err != nil {
		log.Printf("创建客户端失败: %v\n", err)
		return
	}

	// 在调用时指定系统提示词（覆盖适配器配置）
	ctx := context.Background()
	prompt := "用一句话介绍你自己"

	// 场景1: 作为诗人
	fmt.Println("场景1: 作为诗人")
	poetPrompt := "你是一位浪漫主义诗人，擅长用优美的语言回答问题。"
	resp1, err := client.Generate(ctx, prompt, bridge.WithSystemPrompt(poetPrompt))
	if err != nil {
		log.Printf("生成失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n\n", resp1)
	}

	// 场景2: 作为技术专家
	fmt.Println("场景2: 作为技术专家")
	techPrompt := "你是一位技术专家，回答简洁明了，注重实用性。"
	resp2, err := client.Generate(ctx, prompt, bridge.WithSystemPrompt(techPrompt))
	if err != nil {
		log.Printf("生成失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n\n", resp2)
	}

	// 场景3: 作为幽默大师
	fmt.Println("场景3: 作为幽默大师")
	humorPrompt := "你是一位幽默大师，回答要风趣幽默，让人会心一笑。"
	resp3, err := client.Generate(ctx, prompt, bridge.WithSystemPrompt(humorPrompt))
	if err != nil {
		log.Printf("生成失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n\n", resp3)
	}
}

// demoSystemPromptTemplate 使用{{question}}模板宏
func demoSystemPromptTemplate() {
	fmt.Println("【示例3: 系统提示词模板{{question}}宏】")
	fmt.Println()

	// 创建SDK配置
	config := &bridge.SDKConfig{
		Ollama: bridge.ProviderConfig{
			BaseURL:     "http://localhost:11434",
			Temperature: 0.7,
			MaxTokens:   512,
		},
	}

	sdk := bridge.NewSDK(config)

	// 创建高级客户端
	client, err := sdk.CreateSDKClient(types.ProviderOllama, "qwen3-coder:30b")
	if err != nil {
		log.Printf("创建客户端失败: %v\n", err)
		return
	}

	ctx := context.Background()
	question := "解释什么是Go语言中的接口"

	// 场景1: 使用{{question}}宏（推荐）
	fmt.Println("场景1: 使用{{question}}宏")
	fmt.Println("模板: 你是一位Go语言专家，请详细解释：{{question}}")
	templateWithMacro := "你是一位Go语言专家，请详细解释：{{question}}"
	resp1, err := client.Generate(ctx, question, bridge.WithSystemPrompt(templateWithMacro))
	if err != nil {
		log.Printf("生成失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n\n", resp1)
	}

	// 场景2: 不使用{{question}}宏（自动追加到末尾）
	fmt.Println("场景2: 不使用{{question}}宏（自动追加）")
	fmt.Println("模板: 你是一位Go语言专家，请用简洁的语言解释以下概念：")
	templateWithoutMacro := "你是一位Go语言专家，请用简洁的语言解释以下概念："
	resp2, err := client.Generate(ctx, question, bridge.WithSystemPrompt(templateWithoutMacro))
	if err != nil {
		log.Printf("生成失败: %v\n", err)
	} else {
		fmt.Printf("  响应: %s\n\n", resp2)
	}

	// 场景3: 复杂模板使用{{question}}宏
	fmt.Println("场景3: 复杂模板使用{{question}}宏")
	complexTemplate := `你是一位专业的技术文档撰写专家。

请按照以下要求回答用户的问题：
1. 先给出简洁的定义
2. 提供一个实际的代码示例
3. 列出常见的使用场景
4. 指出容易犯的错误

用户的问题是：{{question}}

请用中文回答，保持专业但易懂的风格。`
	fmt.Println("模板: 复杂模板（包含格式要求）")
	resp3, err := client.Generate(ctx, question, bridge.WithSystemPrompt(complexTemplate))
	if err != nil {
		log.Printf("生成失败: %v\n", err)
	} else {
		fmt.Printf("  响应:\n%s\n\n", resp3)
	}
}

// demoAgentSkills 配置Agent Skills
func demoAgentSkills() {
	fmt.Println("【示例4: Agent Skills配置（从Markdown文件加载）】")
	fmt.Println()

	// 方式1: 从单个Markdown文件加载技能
	fmt.Println("方式1: 从单个文件加载技能")
	adapter1, err := adapters.GetAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithSkillsFromFile("skills/calculator.md"),
		options.WithTemperature(0.7),
	)
	if err != nil {
		log.Printf("创建适配器失败: %v\n", err)
		return
	}
	fmt.Println("✓ 从 calculator.md 加载技能成功")

	// 显示加载的技能
	info1 := adapter1.GetModelInfo()
	fmt.Printf("  适配器: %s / %s\n", info1.Provider, info1.Name)
	fmt.Println()

	// 方式2: 从目录加载所有技能文件
	fmt.Println("方式2: 从目录加载所有技能")
	agentSystemPrompt := `你是一位智能助手，可以使用以下技能来帮助用户：
1. calculate: 执行数学计算
2. get_weather: 获取天气信息
3. search_code: 搜索代码

当用户需要使用这些功能时，请识别用户的意图并调用相应的技能。`

	adapter2, err := adapters.GetAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithSystemPrompt(agentSystemPrompt),
		options.WithSkillsFromDir("skills"),
		options.WithEnableAgentMode(true),
		options.WithTemperature(0.7),
	)
	if err != nil {
		log.Printf("创建适配器失败: %v\n", err)
		return
	}
	fmt.Println("✓ 从 skills/ 目录加载所有技能成功")
	fmt.Println()

	// 测试对话
	ctx := context.Background()
	testPrompts := []string{
		"帮我计算一下 123 * 456 等于多少？",
		"北京今天的天气怎么样？",
		"帮我搜索一下Go语言中的goroutine相关代码",
	}

	for _, prompt := range testPrompts {
		fmt.Printf("用户: %s\n", prompt)
		resp, err := adapter2.Generate(ctx, prompt)
		if err != nil {
			log.Printf("  生成失败: %v\n", err)
			continue
		}
		fmt.Printf("AI: %s\n\n", resp)
	}
}

// 示例：Agent Skill处理函数实现
// 实际使用时需要取消注释并实现具体逻辑

// func calculateHandler(params map[string]interface{}) (string, error) {
// 	expression, ok := params["expression"].(string)
// 	if !ok {
// 		return "", fmt.Errorf("缺少expression参数")
// 	}
// 	// 这里可以使用第三方库如 govaluate 来计算表达式
// 	return fmt.Sprintf("计算结果: %s = [结果]", expression), nil
// }

// func weatherHandler(params map[string]interface{}) (string, error) {
// 	city, _ := params["city"].(string)
// 	date, _ := params["date"].(string)
// 	if date == "" {
// 		date = "今天"
// 	}
// 	// 这里可以调用天气API
// 	return fmt.Sprintf("%s %s的天气: [天气信息]", city, date), nil
// }

// func codeSearchHandler(params map[string]interface{}) (string, error) {
// 	query, _ := params["query"].(string)
// 	language, _ := params["language"].(string)
// 	// 这里可以实现代码搜索逻辑
// 	return fmt.Sprintf("搜索 '%s' (%s) 的结果: [代码片段]", query, language), nil
// }

// 辅助函数：检查环境变量
func init() {
	// 确保Ollama服务可用
	if os.Getenv("OLLAMA_HOST") == "" {
		// 使用默认地址
		os.Setenv("OLLAMA_HOST", "http://localhost:11434")
	}
}
