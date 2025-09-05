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

// MCPServer MCP server structure
type MCPServer struct {
	browserService *BrowserService
	router         *gin.Engine
	httpServer     *http.Server
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(browserService *BrowserService) *MCPServer {
	return &MCPServer{
		browserService: browserService,
	}
}

// Start starts the MCP server
func (s *MCPServer) Start(port string) error {
	s.router = s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:    port,
		Handler: s.router,
	}

	// Start server goroutine
	go func() {
		logrus.Infof("Starting MCP server on port %s", port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Server startup failed: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.Errorf("Server shutdown failed: %v", err)
		return err
	}

	logrus.Info("Server shutdown complete")
	return nil
}

// setupRoutes sets up the router configuration
func (s *MCPServer) setupRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// Health check
	router.GET("/health", s.healthHandler)

	// MCP endpoint
	mcpHandler := s.createMCPHandler()
	router.Any("/mcp", gin.WrapH(mcpHandler))
	router.Any("/mcp/*path", gin.WrapH(mcpHandler))

	return router
}

// healthHandler handles health check requests
func (s *MCPServer) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "xhs-exporter-mcp",
	})
}

// corsMiddleware adds CORS headers
func (s *MCPServer) corsMiddleware() gin.HandlerFunc {
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

// createMCPHandler creates the main MCP handler
func (s *MCPServer) createMCPHandler() http.Handler {
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

		switch req.Method {
		case "initialize":
			s.handleInitialize(w, req)
		case "tools/list":
			s.handleToolsList(w, req)
		case "tools/call":
			s.handleToolsCall(w, req)
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

// handleInitialize handles the initialize method
func (s *MCPServer) handleInitialize(w http.ResponseWriter, req JSONRPCRequest) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "xiaohongshu-cover-mcp",
			"version": "v1.0.0",
		},
	}

	s.sendJSONRPCResponse(w, req.ID, result)
}

// handleToolsList handles the tools/list method
func (s *MCPServer) handleToolsList(w http.ResponseWriter, req JSONRPCRequest) {
	tools := map[string]interface{}{
		"tools": []map[string]interface{}{
			{
				"name":        "generate_xiaohongshu_cover",
				"description": "Generate a Xiaohongshu cover image with customizable text and styling",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"baseUrl": map[string]interface{}{
							"type":        "string",
							"description": "The URL to generate cover from (default: http://localhost:3000)",
						},
						"selector": map[string]interface{}{
							"type":        "string",
							"description": "CSS selector of element to screenshot (default: #exportable)",
						},
						"image": map[string]interface{}{
							"type":        "string",
							"description": "Path to the background image (default: /assets/sample1.jpg)",
						},
						"text": map[string]interface{}{
							"type":        "string",
							"description": "Text content to overlay (supports HTML, default: 'Sample Text')",
						},
						"output_path": map[string]interface{}{
							"type":        "string",
							"description": "Output file path for the generated image (default: /tmp/xiaohongshu_cover.png)",
						},
						"fontFamily": map[string]interface{}{
							"type":        "string",
							"description": "Font family name (default: 'Arial')",
						},
						"fontSize": map[string]interface{}{
							"type":        "integer",
							"description": "Font size in pixels (default: 48)",
						},
						"fontWeight": map[string]interface{}{
							"type":        "string",
							"description": "Font weight (default: 'bold')",
						},
						"color": map[string]interface{}{
							"type":        "string",
							"description": "Text color hex code (default: '#ffffff')",
						},
						"backgroundColor": map[string]interface{}{
							"type":        "string",
							"description": "Background color hex code (default: '#000000')",
						},
						"textShadow": map[string]interface{}{
							"type":        "string",
							"description": "CSS text shadow (default: '2px 2px 4px #000000')",
						},
						"border": map[string]interface{}{
							"type":        "string",
							"description": "CSS border (default: '1px solid #000000')",
						},
						"borderRadius": map[string]interface{}{
							"type":        "integer",
							"description": "Border radius in pixels (default: 0)",
						},
						"borderWidth": map[string]interface{}{
							"type":        "integer",
							"description": "Border width in pixels (default: 1)",
						},
						"borderStyle": map[string]interface{}{
							"type":        "string",
							"description": "Border style (default: 'solid')",
						},
						"padding": map[string]interface{}{
							"type":        "integer",
							"description": "Padding in pixels (default: 20)",
						},
						"scaleX": map[string]interface{}{
							"type":        "number",
							"description": "Horizontal scale (default: 1.0)",
						},
						"scaleY": map[string]interface{}{
							"type":        "number",
							"description": "Vertical scale (default: 1.0)",
						},
						"skewX": map[string]interface{}{
							"type":        "number",
							"description": "Horizontal skew in degrees (default: 0)",
						},
						"skewY": map[string]interface{}{
							"type":        "number",
							"description": "Vertical skew in degrees (default: 0)",
						},
						"opacity": map[string]interface{}{
							"type":        "number",
							"description": "Overlay opacity (0.0 to 1.0, default: 0.8)",
						},
						"overlayColor": map[string]interface{}{
							"type":        "string",
							"description": "Overlay color hex code (default: '#000000')",
						},
						"x": map[string]interface{}{
							"type":        "integer",
							"description": "Horizontal position in pixels (default: 50)",
						},
						"y": map[string]interface{}{
							"type":        "integer",
							"description": "Vertical position in pixels (default: 50)",
						},
					},
					"required": []interface{}{},
				},
			},
		},
	}

	s.sendJSONRPCResponse(w, req.ID, tools)
}

// handleToolsCall handles the tools/call method
func (s *MCPServer) handleToolsCall(w http.ResponseWriter, req JSONRPCRequest) {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		s.sendJSONRPCError(w, req.ID, -32602, "Invalid params", "Expected object")
		return
	}

	toolName, ok := params["name"].(string)
	if !ok {
		s.sendJSONRPCError(w, req.ID, -32602, "Invalid params", "Missing tool name")
		return
	}

	args, ok := params["arguments"].(map[string]interface{})
	if !ok {
		args = make(map[string]interface{})
	}

	switch toolName {
	case "generate_xiaohongshu_cover":
		result := s.handleGenerateXiaohongshuCover(context.Background(), args)
		s.sendJSONRPCResponse(w, req.ID, result)
	default:
		s.sendJSONRPCError(w, req.ID, -32601, "Method not found", fmt.Sprintf("Unknown tool: %s", toolName))
	}
}

// sendJSONRPCResponse sends a JSON-RPC response
func (s *MCPServer) sendJSONRPCResponse(w http.ResponseWriter, id interface{}, result interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sendJSONRPCError sends a JSON-RPC error response
func (s *MCPServer) sendJSONRPCError(w http.ResponseWriter, id interface{}, code int, message string, data interface{}) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}