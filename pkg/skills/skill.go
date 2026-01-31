package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Skill Agent Skills 标准定义
// 参考: https://agentskills.io/specification
// Skill 是一个文件夹，包含 skill.md 文件和可选资源
type Skill struct {
	// Name 技能名称（从文件夹名称获取）
	Name string

	// Path 技能文件夹路径
	Path string

	// Metadata 技能元数据（从 skill.md 解析）
	Metadata SkillMetadata

	// Content 技能内容（从 skill.md 解析）
	Content string
}

// SkillMetadata 技能元数据
type SkillMetadata struct {
	// Name 技能名称
	Name string

	// Description 技能描述
	Description string

	// Version 版本号
	Version string

	// Author 作者
	Author string

	// Tags 标签列表
	Tags []string

	// Dependencies 依赖的其他技能
	Dependencies []string
}

// LoadSkill 从文件夹加载 Skill
func LoadSkill(skillPath string) (*Skill, error) {
	// 检查路径是否存在
	info, err := os.Stat(skillPath)
	if err != nil {
		return nil, fmt.Errorf("skill path not found: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("skill path must be a directory: %s", skillPath)
	}

	// 读取 skill.md 文件
	skillMdPath := filepath.Join(skillPath, "skill.md")
	content, err := os.ReadFile(skillMdPath)
	if err != nil {
		return nil, fmt.Errorf("skill.md not found in %s: %w", skillPath, err)
	}

	// 解析 skill.md
	metadata := parseSkillMetadata(string(content))

	skill := &Skill{
		Name:     filepath.Base(skillPath),
		Path:     skillPath,
		Metadata: metadata,
		Content:  string(content),
	}

	return skill, nil
}

// LoadSkillsFromDir 从目录加载所有 Skill
func LoadSkillsFromDir(skillsDir string) ([]*Skill, error) {
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	var skills []*Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillPath := filepath.Join(skillsDir, entry.Name())
		skill, err := LoadSkill(skillPath)
		if err != nil {
			// 记录错误但继续加载其他技能
			fmt.Printf("Warning: failed to load skill from %s: %v\n", skillPath, err)
			continue
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

// parseSkillMetadata 解析 skill.md 的元数据
func parseSkillMetadata(content string) SkillMetadata {
	var metadata SkillMetadata

	lines := strings.Split(content, "\n")
	inFrontMatter := false
	frontMatter := []string{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 检测 front matter 开始/结束
		if trimmed == "---" {
			if !inFrontMatter {
				inFrontMatter = true
				continue
			} else {
				break
			}
		}

		if inFrontMatter {
			frontMatter = append(frontMatter, line)
		}
	}

	// 解析 front matter
	for _, line := range frontMatter {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "name":
					metadata.Name = value
				case "description":
					metadata.Description = value
				case "version":
					metadata.Version = value
				case "author":
					metadata.Author = value
				case "tags":
					// 解析逗号分隔的标签
					tags := strings.Split(value, ",")
					for _, tag := range tags {
						trimmed := strings.TrimSpace(tag)
						if trimmed != "" {
							metadata.Tags = append(metadata.Tags, trimmed)
						}
					}
				case "dependencies":
					// 解析逗号分隔的依赖
					deps := strings.Split(value, ",")
					for _, dep := range deps {
						trimmed := strings.TrimSpace(dep)
						if trimmed != "" {
							metadata.Dependencies = append(metadata.Dependencies, trimmed)
						}
					}
				}
			}
		}
	}

	return metadata
}

// GetSystemPrompt 生成技能的系统提示词
// 将 skill.md 内容转换为系统提示词
func (s *Skill) GetSystemPrompt() string {
	var prompt strings.Builder

	// 添加技能描述
	if s.Metadata.Description != "" {
		prompt.WriteString(fmt.Sprintf("# %s\n\n", s.Metadata.Name))
		prompt.WriteString(fmt.Sprintf("%s\n\n", s.Metadata.Description))
	}

	// 添加完整的 skill.md 内容作为上下文
	prompt.WriteString("## 技能详细说明\n\n")
	prompt.WriteString(s.Content)

	return prompt.String()
}

// GetInstruction 获取技能的指令部分
// 提取 skill.md 中实际的指令内容（去掉 front matter）
func (s *Skill) GetInstruction() string {
	lines := strings.Split(s.Content, "\n")
	var result []string
	inFrontMatter := false
	skipFrontMatter := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "---" {
			if !inFrontMatter {
				inFrontMatter = true
				skipFrontMatter = true
				continue
			} else {
				inFrontMatter = false
				continue
			}
		}

		if inFrontMatter {
			continue
		}

		result = append(result, line)
	}

	return strings.TrimSpace(strings.Join(result, "\n"))
}

// Validate 验证 Skill 是否有效
func (s *Skill) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("skill name is required")
	}

	if s.Content == "" {
		return fmt.Errorf("skill content is empty")
	}

	return nil
}

// Registry Skill 注册表
type Registry struct {
	skills map[string]*Skill
}

// NewRegistry 创建新的 Skill 注册表
func NewRegistry() *Registry {
	return &Registry{
		skills: make(map[string]*Skill),
	}
}

// Register 注册 Skill
func (r *Registry) Register(skill *Skill) error {
	if err := skill.Validate(); err != nil {
		return err
	}

	r.skills[skill.Name] = skill
	return nil
}

// Get 获取 Skill
func (r *Registry) Get(name string) (*Skill, bool) {
	skill, ok := r.skills[name]
	return skill, ok
}

// GetAll 获取所有 Skills
func (r *Registry) GetAll() []*Skill {
	result := make([]*Skill, 0, len(r.skills))
	for _, skill := range r.skills {
		result = append(result, skill)
	}
	return result
}

// LoadFromDir 从目录加载所有 Skills 到注册表
func (r *Registry) LoadFromDir(skillsDir string) error {
	skills, err := LoadSkillsFromDir(skillsDir)
	if err != nil {
		return err
	}

	for _, skill := range skills {
		if err := r.Register(skill); err != nil {
			fmt.Printf("Warning: failed to register skill %s: %v\n", skill.Name, err)
			continue
		}
	}

	return nil
}
