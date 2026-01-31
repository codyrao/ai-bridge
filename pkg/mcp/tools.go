package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

// ToolDefinition 工具定义
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolHandler 工具处理函数
type ToolHandler func(ctx context.Context, params map[string]interface{}) (string, error)

// MCPTool MCP工具结构
type MCPTool struct {
	Definition ToolDefinition
	Handler    ToolHandler
}

// NewTool 创建新工具
func NewTool(name, description string, params map[string]interface{}, handler ToolHandler) *MCPTool {
	return &MCPTool{
		Definition: ToolDefinition{
			Name:        name,
			Description: description,
			Parameters:  params,
		},
		Handler: handler,
	}
}

// ToEinoTool 转换为Eino工具
func (t *MCPTool) ToEinoTool() tool.BaseTool {
	toolInfo := &schema.ToolInfo{
		Name: t.Definition.Name,
		Desc: t.Definition.Description,
	}

	// 设置参数
	if t.Definition.Parameters != nil {
		toolInfo.ParamsOneOf = schema.NewParamsOneOfByJSONSchema(nil)
	}

	return utils.NewTool(toolInfo, func(ctx context.Context, input *schema.ToolCall) (string, error) {
		params, err := ParseToolArguments(input.Function.Arguments)
		if err != nil {
			return "", err
		}
		return t.Handler(ctx, params)
	})
}

// ParseToolArguments 解析工具参数
func ParseToolArguments(arguments string) (map[string]interface{}, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return nil, fmt.Errorf("failed to parse tool arguments: %w", err)
	}
	return params, nil
}

// ToolRegistry 工具注册表
type ToolRegistry struct {
	tools map[string]*MCPTool
}

// NewToolRegistry 创建工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]*MCPTool),
	}
}

// Register 注册工具
func (r *ToolRegistry) Register(tool *MCPTool) {
	r.tools[tool.Definition.Name] = tool
}

// Get 获取工具
func (r *ToolRegistry) Get(name string) (*MCPTool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

// GetAll 获取所有工具
func (r *ToolRegistry) GetAll() []*MCPTool {
	result := make([]*MCPTool, 0, len(r.tools))
	for _, tool := range r.tools {
		result = append(result, tool)
	}
	return result
}

// ToEinoTools 转换为Eino工具列表
func (r *ToolRegistry) ToEinoTools() []tool.BaseTool {
	result := make([]tool.BaseTool, 0, len(r.tools))
	for _, t := range r.tools {
		result = append(result, t.ToEinoTool())
	}
	return result
}

// CreateParameterSchema 创建参数schema
func CreateParameterSchema(properties map[string]interface{}, required []string) map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": properties,
		"required":   required,
	}
}

// CreateStringProperty 创建字符串属性
func CreateStringProperty(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": description,
	}
}

// CreateNumberProperty 创建数字属性
func CreateNumberProperty(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "number",
		"description": description,
	}
}

// CreateIntegerProperty 创建整数属性
func CreateIntegerProperty(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "integer",
		"description": description,
	}
}

// CreateBooleanProperty 创建布尔属性
func CreateBooleanProperty(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "boolean",
		"description": description,
	}
}

// CreateArrayProperty 创建数组属性
func CreateArrayProperty(description string, items map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":        "array",
		"description": description,
		"items":       items,
	}
}

// CreateEnumProperty 创建枚举属性
func CreateEnumProperty(description string, enum []string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": description,
		"enum":        enum,
	}
}

// ExampleTools 示例工具集合
func ExampleTools() []*MCPTool {
	// 计算器工具
	calcTool := NewTool(
		"calculator",
		"执行数学计算",
		CreateParameterSchema(
			map[string]interface{}{
				"expression": CreateStringProperty("数学表达式，例如: 1 + 2 * 3"),
			},
			[]string{"expression"},
		),
		func(ctx context.Context, params map[string]interface{}) (string, error) {
			expression, ok := params["expression"].(string)
			if !ok {
				return "", fmt.Errorf("expression parameter is required")
			}
			// 这里可以实现实际的计算逻辑
			return fmt.Sprintf("计算结果: %s = [计算逻辑待实现]", expression), nil
		},
	)

	// 天气查询工具
	weatherTool := NewTool(
		"weather",
		"查询指定城市的天气",
		CreateParameterSchema(
			map[string]interface{}{
				"city":    CreateStringProperty("城市名称"),
				"country": CreateStringProperty("国家代码（可选）"),
			},
			[]string{"city"},
		),
		func(ctx context.Context, params map[string]interface{}) (string, error) {
			city, ok := params["city"].(string)
			if !ok {
				return "", fmt.Errorf("city parameter is required")
			}
			country, _ := params["country"].(string)
			if country != "" {
				return fmt.Sprintf("%s, %s 的天气: 晴朗，25°C", city, country), nil
			}
			return fmt.Sprintf("%s 的天气: 晴朗，25°C", city), nil
		},
	)

	// JSON格式化工具
	jsonTool := NewTool(
		"format_json",
		"格式化JSON字符串",
		CreateParameterSchema(
			map[string]interface{}{
				"json": CreateStringProperty("需要格式化的JSON字符串"),
			},
			[]string{"json"},
		),
		func(ctx context.Context, params map[string]interface{}) (string, error) {
			jsonStr, ok := params["json"].(string)
			if !ok {
				return "", fmt.Errorf("json parameter is required")
			}
			var data interface{}
			if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
				return "", fmt.Errorf("invalid json: %w", err)
			}
			formatted, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return "", fmt.Errorf("failed to format json: %w", err)
			}
			return string(formatted), nil
		},
	)

	return []*MCPTool{calcTool, weatherTool, jsonTool}
}
