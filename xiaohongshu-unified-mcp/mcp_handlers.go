package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

const defaultText = `
8 月 3 日入园人数: <span style="color: #ff0000; font-weight: bold;">19999</span><br/>天气晴朗适合游玩
`

// handleGenerateXiaohongshuCover handles the generate_xiaohongshu_cover tool call
func (s *UnifiedMCPServer) handleGenerateXiaohongshuCover(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: Generating Xiaohongshu cover")

	// Set default values
	defaults := map[string]interface{}{
		"baseUrl":         "http://localhost:3000",
		"selector":        "#exportable",
		"image":           "/assets/6.jpg",
		"text":            defaultText,
		"output_path":     "/tmp/xiaohongshu_cover.png",
		"fontFamily":      "Arial",
		"fontSize":        48,
		"fontWeight":      "bold",
		"color":           "#0e0d0c",
		"backgroundColor": "#f4f750",
		"textShadow":      "2px 2px 4px #000000",
		"border":          "1px solid #000000",
		"borderRadius":    32,
		"borderWidth":     3,
		"borderStyle":     "dashed",
		"padding":         20,
		"scaleX":          1.0,
		"scaleY":          1.0,
		"skewX":           0.0,
		"skewY":           0.0,
		"opacity":         0.8,
		"overlayColor":    "#443c3c",
		"x":               50,
		"y":               50,
	}

	// Apply defaults for missing parameters
	for key, defaultValue := range defaults {
		if _, exists := args[key]; !exists {
			args[key] = defaultValue
		}
	}

	// Get output path
	outputPath, _ := args["output_path"].(string)

	// Build URL with parameters
	baseURL, _ := args["baseUrl"].(string)
	urlParams := url.Values{}

	// Add all parameters to URL
	for key, value := range args {
		if key == "output_path" {
			continue // Skip output_path as it's not a URL parameter
		}

		switch v := value.(type) {
		case string:
			urlParams.Set(key, v)
		case int:
			urlParams.Set(key, strconv.Itoa(v))
		case float64:
			urlParams.Set(key, strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			urlParams.Set(key, strconv.FormatBool(v))
		}
	}

	// Construct full URL
	fullURL := baseURL + "?" + urlParams.Encode()

	logrus.Infof("MCP: Generated URL: %s", fullURL)

	// Build screenshot request
	req := &ScreenshotRequest{
		URL:        fullURL,
		Selector:   args["selector"].(string),
		OutputPath: outputPath,
		WaitTime:   5, // Wait longer for the page to render
	}

	logrus.Infof("MCP: Taking screenshot for Xiaohongshu cover generation")

	// Execute screenshot
	result, err := s.browserService.TakeScreenshot(ctx, req)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Failed to generate cover: %v", err),
			}},
			IsError: true,
		}
	}

	if !result.Success {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Screenshot failed: %s", result.Error),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: fmt.Sprintf("Xiaohongshu cover generated successfully: %s", result.OutputPath),
		}},
	}
}

// handleCheckLoginStatus 处理检查登录状态
func (s *UnifiedMCPServer) handleCheckLoginStatus(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 检查登录状态")

	status, err := s.xiaohongshuService.CheckLoginStatus(ctx)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "检查登录状态失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	resultText := fmt.Sprintf("登录状态检查成功: %+v", status)
	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: resultText,
		}},
	}
}

// handlePublishContent 处理发布内容
func (s *UnifiedMCPServer) handlePublishContent(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 发布内容")

	// 解析参数
	title, _ := args["title"].(string)
	content, _ := args["content"].(string)
	imagePathsInterface, _ := args["images"].([]interface{})

	var imagePaths []string
	for _, path := range imagePathsInterface {
		if pathStr, ok := path.(string); ok {
			imagePaths = append(imagePaths, pathStr)
		}
	}

	logrus.Infof("MCP: 发布内容 - 标题: %s, 图片数量: %d", title, len(imagePaths))

	// 构建发布请求
	req := &PublishRequest{
		Title:   title,
		Content: content,
		Images:  imagePaths,
	}

	// 执行发布
	result, err := s.xiaohongshuService.PublishContent(ctx, req)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "发布失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	resultText := fmt.Sprintf("内容发布成功: %+v", result)
	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: resultText,
		}},
	}
}

// handleListFeeds 处理获取Feeds列表
func (s *UnifiedMCPServer) handleListFeeds(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 获取Feeds列表")

	result, err := s.xiaohongshuService.ListFeeds(ctx)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "获取Feeds列表失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 格式化输出，转换为JSON字符串
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "序列化结果失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: string(jsonData),
		}},
	}
}

// handleSearchFeeds 处理搜索Feeds
func (s *UnifiedMCPServer) handleSearchFeeds(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 搜索Feeds")

	// 解析参数
	keyword, _ := args["keyword"].(string)
	limit, _ := args["limit"].(float64) // JSON numbers are float64

	if keyword == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "搜索关键词不能为空",
			}},
			IsError: true,
		}
	}

	// 构建搜索请求
	req := &SearchRequest{
		Keyword: keyword,
		Limit:   int(limit),
	}

	logrus.Infof("MCP: 搜索关键词: %s, 限制数量: %d", keyword, req.Limit)

	// 执行搜索
	result, err := s.xiaohongshuService.SearchFeeds(ctx, req.Keyword)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "搜索失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	// 格式化输出，转换为JSON字符串
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "序列化结果失败: " + err.Error(),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: string(jsonData),
		}},
	}
}