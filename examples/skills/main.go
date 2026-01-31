package main

import (
	"context"
	"fmt"
	"log"

	"ai-bridge/pkg/adapters"
	"ai-bridge/pkg/options"
	"ai-bridge/pkg/skills"
	"ai-bridge/pkg/types"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("     Agent Skills 使用示例")
	fmt.Println("========================================")
	fmt.Println()

	// 示例1: 直接加载 Skill 文件夹
	demoLoadSkill()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println()

	// 示例2: 使用 Skill 进行对话
	demoSkillConversation()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println()

	// 示例3: 批量加载多个 Skills
	demoMultipleSkills()
}

// demoLoadSkill 演示如何加载 Skill
func demoLoadSkill() {
	fmt.Println("【示例1: 加载 Agent Skill】")
	fmt.Println()

	// 从文件夹加载单个 Skill
	skill, err := skills.LoadSkill(".github/skills/calculator")
	if err != nil {
		log.Printf("加载 Skill 失败: %v\n", err)
		return
	}

	fmt.Printf("✓ Skill 加载成功\n")
	fmt.Printf("  名称: %s\n", skill.Name)
	fmt.Printf("  描述: %s\n", skill.Metadata.Description)
	fmt.Printf("  版本: %s\n", skill.Metadata.Version)
	fmt.Printf("  作者: %s\n", skill.Metadata.Author)
	fmt.Printf("  标签: %v\n", skill.Metadata.Tags)
	fmt.Println()

	// 显示 Skill 的系统提示词
	fmt.Println("生成的系统提示词（前200字符）:")
	prompt := skill.GetSystemPrompt()
	if len(prompt) > 200 {
		fmt.Printf("%s...\n", prompt[:200])
	} else {
		fmt.Println(prompt)
	}
	fmt.Println()

	// 显示 Skill 指令（去掉 front matter）
	fmt.Println("Skill 指令内容（前200字符）:")
	instruction := skill.GetInstruction()
	if len(instruction) > 200 {
		fmt.Printf("%s...\n", instruction[:200])
	} else {
		fmt.Println(instruction)
	}
}

// demoSkillConversation 演示如何使用 Skill 进行对话
func demoSkillConversation() {
	fmt.Println("【示例2: 使用 Skill 进行对话】")
	fmt.Println()

	// 加载 Skill
	skill, err := skills.LoadSkill(".github/skills/calculator")
	if err != nil {
		log.Printf("加载 Skill 失败: %v\n", err)
		return
	}

	// 创建适配器，使用 Skill 作为系统提示词
	// 将 Skill 的内容添加到系统提示词中
	systemPrompt := `你是一位智能助手，可以使用以下技能：

` + skill.GetSystemPrompt()

	adapter, err := adapters.GetAdapter(
		types.ProviderOllama,
		"qwen3-coder:30b",
		options.WithBaseURL("http://localhost:11434"),
		options.WithSystemPrompt(systemPrompt),
		options.WithTemperature(0.7),
	)
	if err != nil {
		log.Printf("创建适配器失败: %v\n", err)
		return
	}

	fmt.Println("✓ 适配器创建成功（已加载 calculator skill）")
	fmt.Println()

	// 测试对话
	ctx := context.Background()
	testPrompts := []string{
		"帮我计算 123 乘以 456 等于多少？",
		"2 的 10 次方是多少？",
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

// demoMultipleSkills 演示如何批量加载多个 Skills
func demoMultipleSkills() {
	fmt.Println("【示例3: 批量加载多个 Skills】")
	fmt.Println()

	// 创建 Skill 注册表
	registry := skills.NewRegistry()

	// 从目录加载所有 Skills
	err := registry.LoadFromDir(".github/skills")
	if err != nil {
		log.Printf("加载 Skills 失败: %v\n", err)
		return
	}

	// 显示所有加载的 Skills
	allSkills := registry.GetAll()
	fmt.Printf("✓ 成功加载 %d 个 Skills:\n", len(allSkills))
	for _, skill := range allSkills {
		fmt.Printf("  - %s: %s\n", skill.Name, skill.Metadata.Description)
	}
	fmt.Println()

	// 获取特定 Skill
	if calcSkill, ok := registry.Get("calculator"); ok {
		fmt.Println("获取 calculator skill:")
		fmt.Printf("  名称: %s\n", calcSkill.Metadata.Name)
		fmt.Printf("  版本: %s\n", calcSkill.Metadata.Version)
	}

	if weatherSkill, ok := registry.Get("weather"); ok {
		fmt.Println("\n获取 weather skill:")
		fmt.Printf("  名称: %s\n", weatherSkill.Metadata.Name)
		fmt.Printf("  版本: %s\n", weatherSkill.Metadata.Version)
	}

	// 使用多个 Skills 创建系统提示词
	fmt.Println("\n组合多个 Skills 的系统提示词:")
	var combinedPrompt string
	combinedPrompt += "你是一位智能助手，具备以下技能：\n\n"

	for _, skill := range allSkills {
		combinedPrompt += fmt.Sprintf("## %s\n", skill.Metadata.Name)
		combinedPrompt += fmt.Sprintf("%s\n\n", skill.Metadata.Description)
	}

	fmt.Printf("前300字符: %s...\n", combinedPrompt[:min(300, len(combinedPrompt))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
