package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	// 解析命令行参数
	headless := flag.Bool("headless", true, "Run browser in headless mode")
	port := flag.String("port", ":18060", "Server port")
	flag.Parse()

	logrus.Infof("Starting Xiaohongshu Unified MCP Server...")
	logrus.Infof("Headless mode: %v", *headless)
	logrus.Infof("Port: %s", *port)

	// 初始化浏览器服务
	browserService := NewBrowserService(*headless)
	if browserService == nil {
		logrus.Error("Failed to initialize browser service")
		os.Exit(1)
	}

	// 初始化小红书服务
	xiaohongshuService := NewXiaohongshuService(*headless)
	if xiaohongshuService == nil {
		logrus.Error("Failed to initialize xiaohongshu service")
		os.Exit(1)
	}

	// 创建统一MCP服务器
	mcpServer := NewUnifiedMCPServer(browserService, xiaohongshuService)

	// 启动服务器
	logrus.Info("Available tools:")
	logrus.Info("  - generate_xiaohongshu_cover: Generate Xiaohongshu cover images")
	logrus.Info("  - check_login_status: Check Xiaohongshu login status")
	logrus.Info("  - publish_content: Publish content to Xiaohongshu")
	logrus.Info("  - list_feeds: List Xiaohongshu feeds")
	logrus.Info("  - search_feeds: Search Xiaohongshu feeds")

	if err := mcpServer.Start(*port); err != nil {
		logrus.Errorf("Server failed to start: %v", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped gracefully")
}