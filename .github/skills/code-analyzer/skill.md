---
name: code-analyzer
description: 分析代码质量，识别潜在问题，提供优化建议。适用于代码审查、重构建议、性能优化等场景。
version: 1.0.0
author: ai-bridge
tags: code, analysis, review, quality
dependencies:
---

# Code Analyzer Skill

## 概述

这是一个代码分析技能，帮助 AI 助手对代码进行深度分析，识别问题并提供改进建议。

## 使用场景

- 代码审查和 Code Review
- 识别潜在的 bug 和安全漏洞
- 性能优化建议
- 代码风格和质量检查
- 重构建议

## 工作流程

1. **接收代码**：获取用户提供的代码片段或文件路径
2. **语言识别**：自动识别编程语言
3. **静态分析**：分析代码结构、复杂度、依赖关系
4. **问题识别**：查找潜在问题和改进点
5. **生成报告**：输出详细的分析报告

## 分析维度

### 1. 代码质量
- 代码复杂度（圈复杂度）
- 代码重复率
- 函数/方法长度
- 命名规范

### 2. 潜在问题
- 空指针风险
- 资源泄漏
- 并发问题
- 异常处理

### 3. 安全漏洞
- SQL 注入风险
- XSS 漏洞
- 敏感信息泄露
- 不安全的依赖

### 4. 性能优化
- 算法复杂度
- 内存使用
- I/O 操作优化
- 缓存策略

## 输入参数

### 必需参数
- **code** (string): 代码内容或文件路径

### 可选参数
- **language** (string): 编程语言（自动检测）
  - 示例："go", "python", "java", "javascript"
- **focus** (string): 分析重点
  - 可选值："quality", "security", "performance", "all"
  - 默认值："all"

## 输出格式

```json
{
  "summary": "代码分析摘要",
  "issues": [
    {
      "severity": "high|medium|low",
      "type": "bug|security|performance|style",
      "message": "问题描述",
      "line": 42,
      "suggestion": "改进建议"
    }
  ],
  "metrics": {
    "complexity": 15,
    "lines_of_code": 120,
    "functions": 8
  }
}
```

## 示例

### 示例 1：分析 Go 代码
**用户**：帮我分析一下这段代码
```go
func process(data []string) {
    for i := 0; i < len(data); i++ {
        go func() {
            fmt.Println(data[i])
        }()
    }
}
```

**AI**：发现以下问题：
- **严重**：闭包循环变量捕获问题（第4行）
  - 建议：将循环变量作为参数传递给 goroutine
- **中等**：缺少错误处理
- **低**：函数缺少文档注释

### 示例 2：分析 Python 代码
**用户**：检查这段代码的安全性
```python
def query_db(user_input):
    cursor.execute(f"SELECT * FROM users WHERE name = '{user_input}'")
    return cursor.fetchall()
```

**AI**：**严重安全漏洞**：SQL 注入风险（第2行）
- 建议：使用参数化查询
```python
cursor.execute("SELECT * FROM users WHERE name = ?", (user_input,))
```

## 注意事项

1. 分析结果仅供参考，需要人工复核
2. 复杂逻辑可能需要更多上下文
3. 某些问题可能是误报
4. 建议结合单元测试验证修改

## 最佳实践

1. 定期进行代码分析
2. 将高优先级问题加入待办
3. 建立团队代码规范
4. 结合 CI/CD 自动化检查
