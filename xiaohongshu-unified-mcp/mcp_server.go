package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UnifiedMCPServer 统一MCP服务器结构体
type UnifiedMCPServer struct {
	browserService     *BrowserService
	xiaohongshuService *XiaohongshuService
	router             *gin.Engine
	httpServer         *http.Server
}

// NewUnifiedMCPServer 创建新的统一MCP服务器实例
func NewUnifiedMCPServer(browserService *BrowserService, xiaohongshuService *XiaohongshuService) *UnifiedMCPServer {
	return &UnifiedMCPServer{
		browserService:     browserService,
		xiaohongshuService: xiaohongshuService,
	}
}

// Start 启动MCP服务器
func (s *UnifiedMCPServer) Start(port string) error {
	s.router = s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:    port,
		Handler: s.router,
	}

	// 启动服务器goroutine
	go func() {
		logrus.Infof("Starting Unified MCP server on port %s", port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Server startup failed: %v", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.Errorf("Server shutdown failed: %v", err)
		return err
	}

	logrus.Info("Server shutdown complete")
	return nil
}

// setupRoutes 设置路由配置
func (s *UnifiedMCPServer) setupRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// 健康检查
	router.GET("/health", s.healthHandler)

	// MCP端点
	mcpHandler := s.createMCPHandler()
	router.Any("/mcp", gin.WrapH(mcpHandler))
	router.Any("/mcp/*path", gin.WrapH(mcpHandler))

	return router
}

// healthHandler 处理健康检查请求
func (s *UnifiedMCPServer) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "xiaohongshu-unified-mcp",
		"tools": []string{
			"generate_xiaohongshu_cover",
			"check_login_status",
			"publish_content",
			"list_feeds",
			"search_feeds",
		},
	})
}

// corsMiddleware 添加CORS头
func (s *UnifiedMCPServer) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// createMCPHandler 创建MCP处理器
func (s *UnifiedMCPServer) createMCPHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.sendJSONRPCError(w, req.ID, -32700, "Parse error", err.Error())
			return
		}

		logrus.Infof("Received MCP request: method=%s, id=%v", req.Method, req.ID)

		switch req.Method {
		case "initialize":
			s.handleInitialize(w, req)
		case "tools/list":
			s.handleToolsList(w, req)
		case "tools/call":
			s.handleToolsCall(w, req, r.Context())
		case "notifications/initialized":
			// Client notification that initialization is complete, no response needed
			logrus.Info("MCP: Client initialization complete")
			return
		case "notifications/cancelled":
			// Client notification of cancelled request, just log it
			logrus.Info("MCP: Received cancellation notification")
			return
		default:
			s.sendJSONRPCError(w, req.ID, -32601, "Method not found", fmt.Sprintf("Unknown method: %s", req.Method))
		}
	})
}

// handleInitialize 处理初始化请求
func (s *UnifiedMCPServer) handleInitialize(w http.ResponseWriter, req JSONRPCRequest) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "xiaohongshu-unified-mcp",
			"version": "v1.0.0",
		},
	}

	s.sendJSONRPCResponse(w, req.ID, result)
}

// handleToolsList 处理工具列表请求
func (s *UnifiedMCPServer) handleToolsList(w http.ResponseWriter, req JSONRPCRequest) {
	tools := []map[string]interface{}{
		{
			"name":        "generate_xiaohongshu_cover",
			"description": "Generate Xiaohongshu cover image with customizable parameters",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"baseUrl":         map[string]string{"type": "string", "description": "Base URL for the cover generator"},
					"selector":        map[string]string{"type": "string", "description": "CSS selector for screenshot"},
					"image":           map[string]string{"type": "string", "description": "Background image path"},
					"text":            map[string]string{"type": "string", "description": "Text content to display"},
					"output_path":     map[string]string{"type": "string", "description": "Output file path"},
					"fontFamily":      map[string]string{"type": "string", "description": "Font family"},
					"fontSize":        map[string]string{"type": "number", "description": "Font size"},
					"color":           map[string]string{"type": "string", "description": "Text color"},
					"backgroundColor": map[string]string{"type": "string", "description": "Background color"},
				},
			},
		},
		{
			"name":        "check_login_status",
			"description": "Check Xiaohongshu login status",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "publish_content",
			"description": "Publish content to Xiaohongshu",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":   map[string]string{"type": "string", "description": "Content title"},
					"content": map[string]string{"type": "string", "description": "Content body"},
					"images":  map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}, "description": "Image paths"},
				},
				"required": []string{"title", "content", "images"},
			},
		},
		{
			"name":        "list_feeds",
			"description": "List Xiaohongshu feeds",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "search_feeds",
			"description": "Search Xiaohongshu feeds by keyword",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"keyword": map[string]string{"type": "string", "description": "Search keyword"},
					"limit":   map[string]string{"type": "number", "description": "Maximum number of results"},
				},
				"required": []string{"keyword"},
			},
		},
	}

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  map[string]interface{}{"tools": tools},
		ID:      req.ID,
	}

	json.NewEncoder(w).Encode(response)
}

// handleToolsCall 处理工具调用请求
func (s *UnifiedMCPServer) handleToolsCall(w http.ResponseWriter, req JSONRPCRequest, ctx context.Context) {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		s.sendJSONRPCError(w, req.ID, -32602, "Invalid params", "Params must be an object")
		return
	}

	toolCall := MCPToolCall{
		Name:      params["name"].(string),
		Arguments: params["arguments"].(map[string]interface{}),
	}

	logrus.Infof("Executing tool: %s", toolCall.Name)

	var result *MCPToolResult

	switch toolCall.Name {
	case "generate_xiaohongshu_cover":
		result = s.handleGenerateXiaohongshuCover(ctx, toolCall.Arguments)
	case "check_login_status":
		result = s.handleCheckLoginStatus(ctx, toolCall.Arguments)
	case "publish_content":
		result = s.handlePublishContent(ctx, toolCall.Arguments)
	case "list_feeds":
		result = s.handleListFeeds(ctx, toolCall.Arguments)
	case "search_feeds":
		result = s.handleSearchFeeds(ctx, toolCall.Arguments)
	default:
		s.sendJSONRPCError(w, req.ID, -32601, "Method not found", fmt.Sprintf("Unknown tool: %s", toolCall.Name))
		return
	}

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      req.ID,
	}

	json.NewEncoder(w).Encode(response)
}

// sendJSONRPCResponse 发送JSON-RPC响应
func (s *UnifiedMCPServer) sendJSONRPCResponse(w http.ResponseWriter, id interface{}, result interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sendJSONRPCError 发送JSON-RPC错误响应
func (s *UnifiedMCPServer) sendJSONRPCError(w http.ResponseWriter, id interface{}, code int, message, data string) {
	errorResp := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errorResp)
}