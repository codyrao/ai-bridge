---
name: code-search
description: 在代码库中搜索相关代码片段，支持按语言、文件类型过滤。当用户需要查找代码时使用此技能。
version: "1.0"
metadata:
  author: ai-bridge
  tags: code, search, analysis
---

# Code Search Skill

## 概述

这是一个代码搜索技能，帮助 AI 助手在代码库中快速定位相关代码。

## 使用场景

- 用户想查找某个函数的定义
- 用户需要搜索特定的代码模式
- 用户想查找 TODO 或 FIXME 注释
- 用户需要分析代码结构

## 工作流程

1. **理解搜索意图**：分析用户想查找什么
2. **构建搜索查询**：提取关键词和过滤条件
3. **执行搜索**：在代码库中搜索匹配项
4. **格式化结果**：展示文件路径、行号和代码片段

## 搜索能力

### 基本搜索
- 按关键词搜索：函数名、变量名、字符串
- 按代码模式搜索：正则表达式
- 按注释搜索：TODO, FIXME, BUG 等

### 过滤条件
- **language**: 编程语言（go, python, java, javascript 等）
- **file_pattern**: 文件模式（*.go, *test*.py 等）
- **path**: 搜索路径
- **max_results**: 最大返回结果数

## 示例

### 示例 1：搜索函数定义
**用户**：帮我找一下项目中所有的 main 函数
**AI**：我来为您搜索项目中的 main 函数。

找到以下 main 函数：
1. **cmd/server/main.go** (第 15 行)
   ```go
   func main() {
       server := NewServer()
       server.Start()
   }
   ```

2. **cmd/client/main.go** (第 20 行)
   ```go
   func main() {
       client := NewClient()
       client.Run()
   }
   ```

### 示例 2：搜索 TODO 注释
**用户**：项目中有哪些待办事项？
**AI**：我找到了以下 TODO 注释：

1. **server.go** (第 42 行)
   ```go
   // TODO: add error handling
   ```

2. **database.go** (第 88 行)
   ```go
   // TODO: implement connection pooling
   ```

### 示例 3：搜索特定模式
**用户**：查找所有 User 结构体的定义
**AI**：找到以下 User 结构体定义：

1. **models/user.go** (第 10 行)
   ```go
   type User struct {
       ID       int
       Name     string
       Email    string
   }
   ```

## 注意事项

1. 搜索关键词应该具体明确
2. 合理使用过滤条件缩小搜索范围
3. 对于大型代码库，限制返回结果数量
4. 搜索结果包含上下文信息便于理解

## 最佳实践

1. 先进行宽泛搜索，再逐步细化
2. 结合文件路径和语言过滤提高效率
3. 使用正则表达式进行复杂模式匹配
4. 保存常用搜索模式便于复用
