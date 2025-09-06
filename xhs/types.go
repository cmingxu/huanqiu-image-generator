package main

// HTTP API 响应类型

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
}

// JSON-RPC 相关类型

// JSONRPCRequest JSON-RPC 请求
type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
	ID      any    `json:"id"`
}

// JSONRPCResponse JSON-RPC 响应
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  any           `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      any           `json:"id"`
}

// JSONRPCError JSON-RPC 错误
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// MCP 相关类型

// MCPToolCall MCP 工具调用
type MCPToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// MCPToolResult MCP 工具结果
type MCPToolResult struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent MCP 内容
type MCPContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Browser service types (for cover generation)

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

// Xiaohongshu service types

// PublishRequest 发布请求
type PublishRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

// PublishResponse 发布响应
type PublishResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Images  int    `json:"images"`
	Status  string `json:"status"`
	PostID  string `json:"post_id,omitempty"`
}

// LoginStatusResponse 登录状态响应
type LoginStatusResponse struct {
	IsLoggedIn bool   `json:"is_logged_in"`
	Username   string `json:"username,omitempty"`
}

// Feed 动态信息
type Feed struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Likes    int    `json:"likes"`
	Comments int    `json:"comments"`
	URL      string `json:"url"`
}

// FeedsListResponse Feeds列表响应
type FeedsListResponse struct {
	Feeds []Feed `json:"feeds"`
	Count int    `json:"count"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit,omitempty"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Keyword string `json:"keyword"`
	Results []Feed `json:"results"`
	Total   int    `json:"total"`
}