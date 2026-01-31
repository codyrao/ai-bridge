package main

import (
	"context"
	"fmt"
	"log"

	"ai-bridge/pkg/adapters"
	"ai-bridge/pkg/options"
	"ai-bridge/pkg/types"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("     Claude风格Tool Use示例")
	fmt.Println("========================================")
	fmt.Println()

	// 示例1: 使用代码定义工具
	demoToolsInCode()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println()

	// 示例2: 从JSON文件加载工具
	demoToolsFromFile()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println()

	// 示例3: 工具调用流程演示
	demoToolCallFlow()
}

// demoToolsInCode 在代码中定义工具（Claude风格）
func demoToolsInCode() {
	fmt.Println("【示例1: 代码中定义工具（Claude风格）】")
	fmt.Println()

	// 定义计算器工具
	calculatorTool := types.Tool{
		Name:        "calculator",
		Description: "执行数学计算，支持基本运算（加减乘除）、幂运算和开方。当用户需要进行数学计算时，使用此工具。",
		InputSchema: types.ToolInputSchema{
			Type: "object",
			Properties: map[string]types.ToolProperty{
				"expression": {
					Type:        "string",
					Description: "数学表达式，例如 \"2 + 2\", \"sqrt(16)\", \"2^10\"",
				},
				"precision": {
					Type:        "integer",
					Description: "结果精度（小数位数），默认为2位",
					Default:     2,
				},
			},
			Required: []string{"expression"},
		},
	}

	// 定义天气工具
	weatherTool := types.Tool{
		Name:        "get_weather",
		Description: "获取指定城市的当前天气信息。当用户询问天气时使用此工具。",
		InputSchema: types.ToolInputSchema{
			Type: "object",
			Properties: map[string]types.ToolProperty{
				"city": {
					Type:        "string",
					Description: "城市名称，例如 \"北京\", \"Shanghai\"",
				},
				"unit": {
					Type:        "string",
					Description: "温度单位",
					Enum:        []string{"celsius", "fahrenheit"},
					Default:     "celsius",
				},
			},
			Required: []string{"city"},
		},
	}

	// 创建适配器，传入工具定义
	adapter, err := adapters.GetAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithTools(calculatorTool, weatherTool),
		options.WithToolHandler("calculator", calculatorHandler),
		options.WithToolHandler("get_weather", weatherHandler),
		options.WithTemperature(0.7),
	)
	if err != nil {
		log.Printf("创建适配器失败: %v\n", err)
		return
	}

	fmt.Println("✓ 适配器创建成功")
	fmt.Printf("  配置的工具:\n")
	fmt.Printf("    - calculator: 执行数学计算\n")
	fmt.Printf("    - get_weather: 获取天气信息\n")
	fmt.Println()

	// 测试对话
	ctx := context.Background()
	testPrompts := []string{
		"帮我计算 123 乘以 456 等于多少？",
		"北京今天的天气怎么样？",
	}

	for _, prompt := range testPrompts {
		fmt.Printf("用户: %s\n", prompt)
		resp, err := adapter.Generate(ctx, prompt)
		if err != nil {
			log.Printf("  生成失败: %v\n", err)
			continue
		}
		fmt.Printf("AI: %s\n\n", resp)
	}
}

// demoToolsFromFile 从JSON文件加载工具
demoToolsFromFile() {
	fmt.Println("【示例2: 从JSON文件加载工具】")
	fmt.Println()

	// 从JSON文件加载工具定义
	calculatorTool, err := loadToolFromFile("skools/calculator.json")
	if err != nil {
		log.Printf("加载calculator工具失败: %v\n", err)
		return
	}

	weatherTool, err := loadToolFromFile("skools/weather.json")
	if err != nil {
		log.Printf("加载weather工具失败: %v\n", err)
		return
	}

	codeSearchTool, err := loadToolFromFile("skools/code_search.json")
	if err != nil {
		log.Printf("加载code_search工具失败: %v\n", err)
		return
	}

	// 创建适配器
	adapter, err := adapters.GetAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithTools(calculatorTool, weatherTool, codeSearchTool),
		options.WithToolHandler("calculator", calculatorHandler),
		options.WithToolHandler("get_weather", weatherHandler),
		options.WithToolHandler("code_search", codeSearchHandler),
		options.WithTemperature(0.7),
	)
	if err != nil {
		log.Printf("创建适配器失败: %v\n", err)
		return
	}

	fmt.Println("✓ 从JSON文件加载工具成功")
	fmt.Printf("  加载的工具:\n")
	fmt.Printf("    - %s: %s\n", calculatorTool.Name, calculatorTool.Description)
	fmt.Printf("    - %s: %s\n", weatherTool.Name, weatherTool.Description)
	fmt.Printf("    - %s: %s\n", codeSearchTool.Name, codeSearchTool.Description)
	fmt.Println()

	// 测试对话
	ctx := context.Background()
	prompt := "帮我搜索一下项目中的 main 函数"
	fmt.Printf("用户: %s\n", prompt)
	resp, err := adapter.Generate(ctx, prompt)
	if err != nil {
		log.Printf("  生成失败: %v\n", err)
		return
	}
	fmt.Printf("AI: %s\n\n", resp)
}

// demoToolCallFlow 演示工具调用流程
demoToolCallFlow() {
	fmt.Println("【示例3: 工具调用流程演示】")
	fmt.Println()

	fmt.Println("Claude风格的Tool Use流程：")
	fmt.Println()
	fmt.Println("1. 定义工具（Tool Definition）")
	fmt.Println("   - 工具名称、描述、输入参数Schema")
	fmt.Println()
	fmt.Println("2. 发送工具给模型")
	fmt.Println("   - 模型了解可用工具及其参数")
	fmt.Println()
	fmt.Println("3. 用户提问")
	fmt.Println("   - 用户：\"北京今天天气怎么样？\"")
	fmt.Println()
	fmt.Println("4. 模型决定使用工具")
	fmt.Println("   - 模型生成 ToolCall：")
	fmt.Println(`     {
       "id": "tool_01X",
       "name": "get_weather",
       "input": {
         "city": "北京",
         "unit": "celsius"
       }
     }`)
	fmt.Println()
	fmt.Println("5. 执行工具")
	fmt.Println("   - 系统调用 get_weather 处理函数")
	fmt.Println("   - 获取天气数据")
	fmt.Println()
	fmt.Println("6. 返回工具结果")
	fmt.Println("   - 将结果包装为 ToolResult：")
	fmt.Println(`     {
       "tool_call_id": "tool_01X",
       "content": {
         "temperature": 25,
         "condition": "晴天",
         "humidity": "45%"
       }
     }`)
	fmt.Println()
	fmt.Println("7. 模型生成最终回复")
	fmt.Println("   - AI：\"北京今天天气晴朗，温度25°C，湿度45%。\"")
	fmt.Println()
	fmt.Println("✓ 完整流程演示结束")
}

// 工具处理函数实现

// calculatorHandler 计算器工具处理函数
func calculatorHandler(input map[string]interface{}) (interface{}, error) {
	expression, ok := input["expression"].(string)
	if !ok || expression == "" {
		return nil, fmt.Errorf("缺少 expression 参数")
	}

	// 这里可以使用第三方库如 govaluate 来计算表达式
	// 简化示例，直接返回表达式
	return map[string]interface{}{
		"expression": expression,
		"result":     "[计算结果]",
		"note":       "实际实现需要使用表达式计算库",
	}, nil
}

// weatherHandler 天气工具处理函数
func weatherHandler(input map[string]interface{}) (interface{}, error) {
	city, ok := input["city"].(string)
	if !ok || city == "" {
		return nil, fmt.Errorf("缺少 city 参数")
	}

	unit := "celsius"
	if u, ok := input["unit"].(string); ok {
		unit = u
	}

	// 模拟天气数据
	return map[string]interface{}{
		"city":        city,
		"temperature": 25,
		"unit":        unit,
		"condition":   "晴天",
		"humidity":    "45%",
		"wind":        "3级",
	}, nil
}

// codeSearchHandler 代码搜索工具处理函数
func codeSearchHandler(input map[string]interface{}) (interface{}, error) {
	query, ok := input["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("缺少 query 参数")
	}

	language := ""
	if lang, ok := input["language"].(string); ok {
		language = lang
	}

	// 模拟搜索结果
	return []map[string]interface{}{
		{
			"file":    "main.go",
			"line":    15,
			"content": "func main() { ... }",
		},
		{
			"file":    "cmd/server/main.go",
			"line":    20,
			"content": "func main() { ... }",
		},
	}, nil
}

// loadToolFromFile 从JSON文件加载工具定义
func loadToolFromFile(path string) (types.Tool, error) {
	// 简化实现，实际应该读取和解析JSON文件
	// 这里返回空工具，实际使用时需要实现文件读取
	return types.Tool{}, fmt.Errorf("文件加载功能需要实现")
}
