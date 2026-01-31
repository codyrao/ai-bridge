package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ai-bridge/pkg/types"

	"gopkg.in/yaml.v3"
)

// LoadFromFile 从文件加载技能定义
// 支持 .md, .markdown, .yaml, .yml, .json 格式
func LoadFromFile(path string) ([]types.AgentSkill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".md", ".markdown":
		return parseMarkdown(data)
	case ".yaml", ".yml":
		return parseYAML(data)
	case ".json":
		return parseJSON(data)
	default:
		// 尝试自动检测格式
		if skills, err := parseYAML(data); err == nil {
			return skills, nil
		}
		if skills, err := parseJSON(data); err == nil {
			return skills, nil
		}
		return parseMarkdown(data)
	}
}

// LoadFromDir 从目录加载所有技能文件
func LoadFromDir(dir string) ([]types.AgentSkill, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	var allSkills []types.AgentSkill
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))

		// 只处理支持的文件类型
		if ext != ".md" && ext != ".markdown" &&
			ext != ".yaml" && ext != ".yml" &&
			ext != ".json" {
			continue
		}

		path := filepath.Join(dir, name)
		skills, err := LoadFromFile(path)
		if err != nil {
			// 记录错误但继续加载其他文件
			fmt.Printf("Warning: failed to load skills from %s: %v\n", path, err)
			continue
		}

		allSkills = append(allSkills, skills...)
	}

	return allSkills, nil
}

// parseYAML 解析YAML格式的技能文件
func parseYAML(data []byte) ([]types.AgentSkill, error) {
	var skillFile types.SkillFile
	if err := yaml.Unmarshal(data, &skillFile); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return skillFile.Skills, nil
}

// parseJSON 解析JSON格式的技能文件
func parseJSON(data []byte) ([]types.AgentSkill, error) {
	var skillFile types.SkillFile
	if err := json.Unmarshal(data, &skillFile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return skillFile.Skills, nil
}

// parseMarkdown 解析Markdown格式的技能文件
// Markdown格式支持在代码块中嵌入YAML或JSON
func parseMarkdown(data []byte) ([]types.AgentSkill, error) {
	content := string(data)

	// 尝试提取YAML代码块
	if skills := extractCodeBlock(content, "yaml", "yml"); skills != nil {
		return parseYAML(skills)
	}

	// 尝试提取JSON代码块
	if skills := extractCodeBlock(content, "json"); skills != nil {
		return parseJSON(skills)
	}

	// 尝试直接解析整个内容为YAML
	if skills, err := parseYAML(data); err == nil && len(skills) > 0 {
		return skills, nil
	}

	return nil, fmt.Errorf("no valid skill definition found in markdown")
}

// extractCodeBlock 从Markdown中提取指定语言的代码块
func extractCodeBlock(content string, languages ...string) []byte {
	for _, lang := range languages {
		// 查找 ```lang 开头的代码块
		startTag := "```" + lang
		startIdx := strings.Index(content, startTag)
		if startIdx == -1 {
			continue
		}

		// 找到代码块内容开始位置
		contentStart := startIdx + len(startTag)
		// 跳过可能的换行符
		if contentStart < len(content) && content[contentStart] == '\n' {
			contentStart++
		}

		// 找到代码块结束标记
		endTag := "```"
		endIdx := strings.Index(content[contentStart:], endTag)
		if endIdx == -1 {
			continue
		}

		return []byte(content[contentStart : contentStart+endIdx])
	}
	return nil
}

// ValidateSkill 验证技能定义是否有效
func ValidateSkill(skill types.AgentSkill) error {
	if skill.Name == "" {
		return fmt.Errorf("skill name is required")
	}
	if skill.Description == "" {
		return fmt.Errorf("skill description is required for %s", skill.Name)
	}
	return nil
}

// MergeSkills 合并多个技能列表，后加载的技能会覆盖同名技能
func MergeSkills(skillLists ...[]types.AgentSkill) []types.AgentSkill {
	skillMap := make(map[string]types.AgentSkill)

	for _, skills := range skillLists {
		for _, skill := range skills {
			if err := ValidateSkill(skill); err != nil {
				fmt.Printf("Warning: invalid skill: %v\n", err)
				continue
			}
			skillMap[skill.Name] = skill
		}
	}

	// 转换回切片
	result := make([]types.AgentSkill, 0, len(skillMap))
	for _, skill := range skillMap {
		result = append(result, skill)
	}
	return result
}
