# XHS Scheduler Usage

这个调度器可以在每天北京时间晚上8点自动执行内容生成和发布工作流。

## 使用方式

### 1. 每日调度模式（推荐）
在每天北京时间20:00自动运行生成和发布流程：
```bash
go run main.go mcp_server.go mcp_handlers.go browser_service.go xiaohongshu_service.go types.go -scheduler
```

### 2. 立即执行一次（测试用）
立即执行一次完整的工作流程：
```bash
go run main.go mcp_server.go mcp_handlers.go browser_service.go xiaohongshu_service.go types.go -run-once
```

### 3. 正常API服务器模式（默认）
启动API服务器，可通过HTTP接口调用：
```bash
go run main.go mcp_server.go mcp_handlers.go browser_service.go xiaohongshu_service.go types.go
```

## 调度器特性

- ✅ **精确时间控制**: 每天北京时间20:00:00执行
- ✅ **时区自动处理**: 自动处理北京时间（Asia/Shanghai）
- ✅ **优雅关闭**: 支持SIGINT/SIGTERM信号优雅停止
- ✅ **错误恢复**: 即使某次执行失败，也会继续下一天的调度
- ✅ **执行窗口**: 1分钟执行窗口（20:00:00-20:00:59）确保可靠性
- ✅ **独立运行**: 不依赖API服务器，纯调度逻辑

## 工作流程

调度器会按以下顺序执行完整的内容生成和发布流程：

1. **信息收集**: 获取天气、农历、交通、访客等信息
2. **内容生成**: 使用DeepSeek LLM生成小红书内容
3. **封面制作**: 自动生成配套的封面图片
4. **自动发布**: 发布到小红书平台

## 日志示例

```
INFO[2025-09-07T09:43:27+08:00] Starting Xiaohongshu Unified Server...
INFO[2025-09-07T09:43:27+08:00] Scheduler mode: true
2025/09/07 09:43:27 🕐 Starting daily scheduler for 8pm Beijing time...
2025/09/07 09:43:27 📅 Next scheduled run: 2025-09-07 20:00:00 CST
```

## 注意事项

1. **配置文件**: 确保`config.json`配置正确
2. **依赖服务**: 需要MCP服务器和小红书MCP服务器运行
3. **网络连接**: 需要稳定的网络连接访问各种API
4. **浏览器**: 需要Chrome/Chromium用于小红书登录

## 停止调度器

使用Ctrl+C或发送SIGTERM信号即可优雅停止调度器。