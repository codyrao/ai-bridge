package mcp

import (
	"context"
	"fmt"
	"testing"
)

func TestNewToolRegistry(t *testing.T) {
	registry := NewToolRegistry()
	if registry == nil {
		t.Fatal("NewToolRegistry() returned nil")
	}
	if registry.tools == nil {
		t.Error("ToolRegistry.tools is nil")
	}
}

func TestToolRegistry_RegisterAndGet(t *testing.T) {
	registry := NewToolRegistry()

	tool := NewTool(
		"test_tool",
		"A test tool",
		CreateParameterSchema(
			map[string]interface{}{
				"param1": CreateStringProperty("Parameter 1"),
			},
			[]string{"param1"},
		),
		func(ctx context.Context, params map[string]interface{}) (string, error) {
			return "test result", nil
		},
	)

	// 注册工具
	registry.Register(tool)

	// 获取工具
	got, ok := registry.Get("test_tool")
	if !ok {
		t.Error("Get() returned false for registered tool")
	}
	if got.Definition.Name != "test_tool" {
		t.Errorf("Expected tool name 'test_tool', got '%s'", got.Definition.Name)
	}

	// 获取不存在的工具
	_, ok = registry.Get("nonexistent")
	if ok {
		t.Error("Get() returned true for nonexistent tool")
	}
}

func TestToolRegistry_GetAll(t *testing.T) {
	registry := NewToolRegistry()

	// 注册多个工具
	for i := 0; i < 3; i++ {
		tool := NewTool(
			fmt.Sprintf("tool_%d", i),
			fmt.Sprintf("Tool %d", i),
			CreateParameterSchema(map[string]interface{}{}, []string{}),
			func(ctx context.Context, params map[string]interface{}) (string, error) {
				return "", nil
			},
		)
		registry.Register(tool)
	}

	all := registry.GetAll()
	if len(all) != 3 {
		t.Errorf("Expected 3 tools, got %d", len(all))
	}
}

func TestMCPTool_Handler(t *testing.T) {
	tool := NewTool(
		"calculator",
		"Calculate sum",
		CreateParameterSchema(
			map[string]interface{}{
				"a": CreateNumberProperty("First number"),
				"b": CreateNumberProperty("Second number"),
			},
			[]string{"a", "b"},
		),
		func(ctx context.Context, params map[string]interface{}) (string, error) {
			a, _ := params["a"].(float64)
			b, _ := params["b"].(float64)
			return fmt.Sprintf("%.0f", a+b), nil
		},
	)

	ctx := context.Background()
	params := map[string]interface{}{
		"a": 10.0,
		"b": 20.0,
	}

	result, err := tool.Handler(ctx, params)
	if err != nil {
		t.Errorf("Handler() returned error: %v", err)
	}
	if result != "30" {
		t.Errorf("Expected '30', got '%s'", result)
	}
}

func TestCreateParameterSchema(t *testing.T) {
	properties := map[string]interface{}{
		"name": CreateStringProperty("Name"),
		"age":  CreateIntegerProperty("Age"),
	}
	required := []string{"name"}

	schema := CreateParameterSchema(properties, required)

	if schema["type"] != "object" {
		t.Errorf("Expected type 'object', got '%v'", schema["type"])
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("properties is not a map")
	}
	if _, ok := props["name"]; !ok {
		t.Error("properties does not contain 'name'")
	}

	req, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required is not a string slice")
	}
	if len(req) != 1 || req[0] != "name" {
		t.Errorf("Expected required ['name'], got %v", req)
	}
}

func TestCreateStringProperty(t *testing.T) {
	prop := CreateStringProperty("A description")
	if prop["type"] != "string" {
		t.Errorf("Expected type 'string', got '%v'", prop["type"])
	}
	if prop["description"] != "A description" {
		t.Errorf("Expected description 'A description', got '%v'", prop["description"])
	}
}

func TestCreateNumberProperty(t *testing.T) {
	prop := CreateNumberProperty("A number")
	if prop["type"] != "number" {
		t.Errorf("Expected type 'number', got '%v'", prop["type"])
	}
}

func TestCreateIntegerProperty(t *testing.T) {
	prop := CreateIntegerProperty("An integer")
	if prop["type"] != "integer" {
		t.Errorf("Expected type 'integer', got '%v'", prop["type"])
	}
}

func TestCreateBooleanProperty(t *testing.T) {
	prop := CreateBooleanProperty("A boolean")
	if prop["type"] != "boolean" {
		t.Errorf("Expected type 'boolean', got '%v'", prop["type"])
	}
}

func TestCreateEnumProperty(t *testing.T) {
	enum := []string{"option1", "option2", "option3"}
	prop := CreateEnumProperty("An enum", enum)
	if prop["type"] != "string" {
		t.Errorf("Expected type 'string', got '%v'", prop["type"])
	}
	propEnum, ok := prop["enum"].([]string)
	if !ok {
		t.Fatal("enum is not a string slice")
	}
	if len(propEnum) != 3 {
		t.Errorf("Expected 3 enum values, got %d", len(propEnum))
	}
}

func TestExampleTools(t *testing.T) {
	tools := ExampleTools()
	if len(tools) == 0 {
		t.Error("ExampleTools() returned empty list")
	}

	// 检查是否包含计算器工具
	foundCalc := false
	for _, tool := range tools {
		if tool.Definition.Name == "calculator" {
			foundCalc = true
			break
		}
	}
	if !foundCalc {
		t.Error("calculator tool not found in example tools")
	}
}
