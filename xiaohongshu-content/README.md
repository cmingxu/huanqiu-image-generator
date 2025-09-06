# Xiaohongshu Content Generator

一个自动化的小红书内容生成和发布工具，集成天气、农历、交通信息，使用 DeepSeek LLM 生成内容，并通过 MCP 工具生成封面图片。

## 功能特性

- 🌤️ **天气信息获取**: 集成 OpenWeatherMap API 获取实时天气数据
- 📅 **农历信息**: 提供中国农历日期、节气、宜忌等信息
- 🚗 **交通信息**: 获取城市交通状况和出行建议
- 🤖 **AI 内容生成**: 使用 DeepSeek LLM 基于收集的信息生成小红书帖子
- 🎨 **封面图片生成**: 通过 MCP 工具自动生成配套封面图片
- 📱 **自动发布**: 自动发布内容和图片到小红书平台

## 项目结构

```
xiaohongshu-content/
├── main.go                    # 主程序入口
├── go.mod                     # Go 模块文件
├── config.json.example        # 配置文件示例
├── README.md                  # 项目说明文档
└── internal/                  # 内部包
    ├── config/                # 配置管理
    │   └── config.go
    ├── weather/               # 天气服务
    │   └── weather.go
    ├── lunar/                 # 农历服务
    │   └── lunar.go
    ├── traffic/               # 交通服务
    │   └── traffic.go
    ├── llm/                   # LLM 服务
    │   └── deepseek.go
    ├── mcp/                   # MCP 客户端
    │   └── client.go
    ├── xhs/                   # 小红书客户端
    │   └── client.go
    └── orchestrator/          # 工作流编排器
        └── orchestrator.go
```

## 快速开始

### 1. 环境要求

- Go 1.19 或更高版本
- 运行中的 xiaohongshu-cover-mcp 服务 (端口 18061)
- 运行中的 xiaohongshu-mcp 服务 (端口 18060)

### 2. 安装依赖

```bash
cd xiaohongshu-content
go mod tidy
```

### 3. 配置设置

复制配置文件示例并填入你的 API 密钥：

```bash
cp config.json.example config.json
```

编辑 `config.json` 文件，填入以下信息：

- `weather_api.api_key`: OpenWeatherMap API 密钥
- `deepseek_llm.api_key`: DeepSeek API 密钥
- 其他服务的 URL 和配置

### 4. 环境变量（可选）

你也可以通过环境变量设置配置：

```bash
export WEATHER_API_KEY="your_weather_api_key"
export DEEPSEEK_API_KEY="your_deepseek_api_key"
export CITY="Beijing"
export MCP_SERVER_URL="http://localhost:18061"
export XHS_SERVER_URL="http://localhost:18060"
```

### 5. 运行程序

```bash
go run main.go
```

## 工作流程

1. **信息收集阶段**
   - 获取当前天气信息（温度、湿度、风速等）
   - 获取农历信息（农历日期、节气、宜忌等）
   - 获取交通状况信息

2. **内容生成阶段**
   - 将收集的信息发送给 DeepSeek LLM
   - 生成符合小红书风格的标题和正文
   - 生成相关标签和图片提示词

3. **图片生成阶段**
   - 使用 MCP 工具根据提示词生成封面图片
   - 优化图片尺寸和风格适配小红书平台

4. **发布阶段**
   - 将生成的内容和图片发布到小红书
   - 返回发布结果和链接

## API 集成

### 天气 API

使用 OpenWeatherMap API 获取天气信息。需要注册账号获取免费 API 密钥。

### DeepSeek LLM

使用 DeepSeek 的聊天 API 生成内容。需要注册账号并获取 API 密钥。

### MCP 服务

- **xiaohongshu-cover-mcp**: 图片生成服务
- **xiaohongshu-mcp**: 小红书发布服务

确保这两个服务在运行此程序前已经启动。

## 配置说明

### 配置文件结构

```json
{
  "weather_api": {
    "api_key": "天气 API 密钥",
    "base_url": "API 基础 URL",
    "city": "城市名称"
  },
  "deepseek_llm": {
    "api_key": "DeepSeek API 密钥",
    "base_url": "DeepSeek API URL",
    "model": "使用的模型名称"
  },
  "mcp": {
    "server_url": "MCP 服务器 URL"
  },
  "xiaohongshu": {
    "server_url": "小红书 MCP 服务器 URL"
  },
  "settings": {
    "post_interval": "发布间隔（如 24h）",
    "log_level": "日志级别"
  }
}
```

## 故障排除

### 常见问题

1. **连接失败**
   - 检查 MCP 服务是否正在运行
   - 验证服务 URL 和端口配置

2. **API 密钥错误**
   - 确认 API 密钥正确且有效
   - 检查 API 配额是否用完

3. **内容生成失败**
   - 检查 DeepSeek API 连接
   - 验证提示词格式

### 日志输出

程序会输出详细的执行日志，包括：
- 🚀 工作流开始
- 📊 信息收集进度
- 🤖 内容生成状态
- 🎨 图片生成结果
- 📱 发布状态
- ✅ 成功完成或 ❌ 错误信息

## 开发说明

### 添加新的信息源

1. 在 `internal/` 下创建新的服务包
2. 实现信息获取接口
3. 在 `orchestrator.go` 中集成新服务
4. 更新 LLM 提示词以包含新信息

### 自定义内容模板

修改 `internal/llm/deepseek.go` 中的 `buildPrompt` 函数来自定义内容生成模板。

## 许可证

MIT License