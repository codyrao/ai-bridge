---
name: weather
description: 获取指定城市的当前天气信息，包括温度、天气状况、湿度等。当用户询问天气时使用此技能。
version: 1.0.0
author: ai-bridge
tags: weather, forecast, location
dependencies:
---

# Weather Skill

## 概述

这是一个天气查询技能，帮助 AI 助手获取指定城市的天气信息。

## 使用场景

- 用户询问某个城市的当前天气
- 用户需要了解温度、湿度、天气状况
- 用户需要天气预警信息
- 用户计划出行需要天气参考

## 工作流程

1. **识别城市**：从用户输入中提取城市名称
2. **标准化城市名**：将城市名转换为标准格式
3. **查询天气数据**：调用天气 API 获取实时数据
4. **格式化输出**：将天气信息以友好的方式呈现

## 输入参数

### 必需参数
- **city** (string): 城市名称，支持中文和英文
  - 示例："北京", "Shanghai", "New York"

### 可选参数
- **unit** (string): 温度单位
  - 可选值："celsius" (摄氏度), "fahrenheit" (华氏度)
  - 默认值："celsius"

## 输出格式

```json
{
  "city": "北京",
  "temperature": 25,
  "unit": "celsius",
  "condition": "晴天",
  "humidity": "45%",
  "wind": "3级",
  "update_time": "2024-01-15 14:30"
}
```

## 示例

### 示例 1：查询今天天气
**用户**：北京今天天气怎么样？
**AI**：北京今天的天气是：晴天，温度25°C，湿度45%，风力3级。

### 示例 2：查询其他城市
**用户**：上海现在多少度？
**AI**：上海现在的温度是22°C，天气多云，湿度60%。

### 示例 3：华氏度查询
**用户**：What's the weather in New York?
**AI**：The weather in New York is: Sunny, 68°F, humidity 50%.

## 注意事项

1. 城市名称支持中英文
2. 默认返回摄氏度，用户可要求华氏度
3. 天气数据为实时数据，可能有短暂延迟
4. 部分偏远地区可能数据不完整

## 相关技能

- `forecast`: 获取未来几天的天气预报
- `air_quality`: 查询空气质量指数
