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
func (s *MCPServer) handleGenerateXiaohongshuCover(ctx context.Context, args map[string]interface{}) *MCPToolResult {
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
				Text: fmt.Sprintf("Xiaohongshu cover generation failed: %s", err.Error()),
			}},
			IsError: true,
		}
	}

	// Format result as JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Failed to format result: %s", err.Error()),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: fmt.Sprintf("Xiaohongshu cover generated successfully:\n%s\n\nGenerated URL: %s", string(jsonData), fullURL),
		}},
	}
}

