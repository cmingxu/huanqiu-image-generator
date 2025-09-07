package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"xiaohongshu-unified/internal/config"
	"xiaohongshu-unified/internal/orchestrator"
	"xiaohongshu-unified/internal/scheduler"
)

func main() {
	// 设置日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	// 解析命令行参数
	headless := flag.Bool("headless", true, "Run browser in headless mode")
	port := flag.String("port", ":18062", "Server port")
	coverBaseURL := flag.String("cover-base-url", "http://localhost:3000", "Base URL for cover generation service")
	coverDir := flag.String("cover-dir", "/Users/kx/Desktop", "Output directory for cover images")
	schedulerMode := flag.Bool("scheduler", false, "Run in scheduler mode (daily 8pm Beijing time)")
	runOnce := flag.Bool("run-once", false, "Run workflow once and exit")
	flag.Parse()

	logrus.Infof("Starting Xiaohongshu Unified Server...")
	logrus.Infof("Headless mode: %v", *headless)
	logrus.Infof("Scheduler mode: %v", *schedulerMode)
	logrus.Infof("Port: %s", *port)

	// Load configuration for content generation
	cfg, err := config.Load()
	if err != nil {
		logrus.Warnf("Failed to load content generation config: %v", err)
		// Continue without content generation features
	}

	// Override config with command-line flags if provided
	if cfg != nil {
		cfg.MCP.BaseURL = *coverBaseURL
		cfg.MCP.OutDir = *coverDir
		logrus.Infof("Cover Base URL: %s", cfg.MCP.BaseURL)
		logrus.Infof("Cover Output Directory: %s", cfg.MCP.OutDir)
	}

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

	// 创建内容生成编排器
	var orch *orchestrator.Orchestrator
	if cfg != nil {
		orch = orchestrator.New(cfg)
		logrus.Info("Content generation orchestrator initialized")
	} else {
		logrus.Warn("Content generation orchestrator not available (config not loaded)")
	}

	// Handle scheduler mode or run-once mode
	if *schedulerMode || *runOnce {
		if orch == nil {
			logrus.Fatal("Cannot run scheduler mode without valid configuration")
		}

		schedulerSvc := scheduler.New(orch)

		if *runOnce {
			logrus.Info("Running workflow once...")
			if err := schedulerSvc.RunOnce(); err != nil {
				logrus.Fatalf("Workflow execution failed: %v", err)
			}
			logrus.Info("Workflow completed successfully")
			return
		}

		// Setup graceful shutdown for scheduler mode
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Start scheduler in a goroutine
		go func() {
			if err := schedulerSvc.Start(); err != nil {
				logrus.Errorf("Scheduler error: %v", err)
			}
		}()

		// Wait for shutdown signal
		<-sigChan
		logrus.Info("Shutdown signal received")
		schedulerSvc.Stop()
		return
	}

	// Regular API server mode
	// 创建统一服务器，包含MCP和内容生成API
	unifiedServer := NewUnifiedServer(mcpServer, orch)

	// 启动服务器
	logrus.Info("Available MCP tools:")
	logrus.Info("  - generate_xiaohongshu_cover: Generate Xiaohongshu cover images")
	logrus.Info("  - check_login_status: Check Xiaohongshu login status")
	logrus.Info("  - publish_content: Publish content to Xiaohongshu")
	logrus.Info("  - list_feeds: List Xiaohongshu feeds")
	logrus.Info("  - search_feeds: Search Xiaohongshu feeds")

	if orch != nil {
		logrus.Info("Available Content Generation APIs:")
		logrus.Info("  - POST /api/generate-and-publish: Auto-generate content, create image, and publish")
		logrus.Info("  - GET /api/status: Get service status")
	}

	logrus.Infof("Starting Unified server on port %s", *port)
	if err := unifiedServer.Start(*port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

// UnifiedServer combines MCP server and content generation API
type UnifiedServer struct {
	mcpServer *UnifiedMCPServer
	orch     *orchestrator.Orchestrator
	router    *gin.Engine
}

// NewUnifiedServer creates a new unified server
func NewUnifiedServer(mcpServer *UnifiedMCPServer, orch *orchestrator.Orchestrator) *UnifiedServer {
	return &UnifiedServer{
		mcpServer: mcpServer,
		orch:     orch,
	}
}

// Start starts the unified server
func (s *UnifiedServer) Start(port string) error {
	s.setupRoutes()
	return s.router.Run(port)
}

// setupRoutes sets up all routes for both MCP and content generation
func (s *UnifiedServer) setupRoutes() {
	s.router = gin.Default()

	// Add CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// MCP routes
	s.router.POST("/", gin.WrapH(s.mcpServer.createMCPHandler()))

	// Content generation API routes
	if s.orch != nil {
		api := s.router.Group("/api")
		{
			api.POST("/generate-and-publish", s.handleGenerateAndPublish)
			api.GET("/status", s.handleStatus)
		}
	}

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// handleGenerateAndPublish handles the auto content generation and publishing
func (s *UnifiedServer) handleGenerateAndPublish(c *gin.Context) {
	if s.orch == nil {
		c.JSON(500, gin.H{"error": "Content generation not configured"})
		return
	}

	logrus.Info("Starting auto content generation and publishing workflow...")

	// Run the orchestrator workflow
	err := s.orch.Run()
	if err != nil {
		logrus.Errorf("Workflow failed: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Workflow failed: %v", err)})
		return
	}

	c.JSON(200, gin.H{
		"message": "Content generation and publishing completed successfully",
		"status":  "success",
	})
}

// handleStatus returns the status of all services
func (s *UnifiedServer) handleStatus(c *gin.Context) {
	status := map[string]interface{}{
		"mcp_server": "running",
	}

	if s.orch != nil {
		status["content_generation"] = "available"
		status["services"] = s.orch.GetServiceStatus()
	} else {
		status["content_generation"] = "not_configured"
	}

	c.JSON(200, status)
}