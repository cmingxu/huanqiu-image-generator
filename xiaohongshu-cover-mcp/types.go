package main

// JSON-RPC related types

// JSONRPCRequest JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
	ID      any    `json:"id"`
}

// JSONRPCResponse JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  any           `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      any           `json:"id"`
}

// JSONRPCError JSON-RPC error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// MCP related types

// MCPToolCall MCP tool call
type MCPToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// MCPToolResult MCP tool result
type MCPToolResult struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent MCP content
type MCPContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Browser service types

// ScreenshotRequest screenshot request parameters
type ScreenshotRequest struct {
	URL        string `json:"url"`
	Selector   string `json:"selector,omitempty"`
	OutputPath string `json:"output_path,omitempty"`
	WaitTime   int    `json:"wait_time,omitempty"`
}

// ScreenshotResult screenshot operation result
type ScreenshotResult struct {
	Success    bool   `json:"success"`
	OutputPath string `json:"output_path,omitempty"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
}