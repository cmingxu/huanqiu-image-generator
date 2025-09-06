package xhs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"xiaohongshu-content/internal/config"
)

// PostRequest represents a request to post content to Xiaohongshu
type PostRequest struct {
	Title   string   `json:"title"`   // Post title
	Content string   `json:"content"` // Post content
	Images  []string `json:"images"`  // Image URLs or paths (at least one required)
}

// LoginStatusRequest represents a request to check login status
type LoginStatusRequest struct{}

// FeedsRequest represents a request to get feeds list
type FeedsRequest struct{}

// SearchRequest represents a request to search content
type SearchRequest struct {
	Keyword string `json:"keyword"` // Search keyword
}

// PostResponse represents the response from posting to Xiaohongshu
type PostResponse struct {
	PostID      string    `json:"post_id"`      // Generated post ID
	URL         string    `json:"url"`          // Post URL
	Status      string    `json:"status"`       // posted, scheduled, failed
	Message     string    `json:"message"`      // Status message
	PostedAt    time.Time `json:"posted_at"`    // When the post was created
	Error       string    `json:"error,omitempty"` // Error message if any
}

// LoginStatusResponse represents the response from checking login status
type LoginStatusResponse struct {
	LoggedIn bool   `json:"logged_in"` // Whether user is logged in
	Message  string `json:"message"`   // Status message
	Error    string `json:"error,omitempty"` // Error message if any
}

// FeedsResponse represents the response from getting feeds list
type FeedsResponse struct {
	Feeds   []Feed `json:"feeds"`   // List of feeds
	Message string `json:"message"` // Status message
	Error   string `json:"error,omitempty"` // Error message if any
}

// SearchResponse represents the response from searching content
type SearchResponse struct {
	Results []SearchResult `json:"results"` // Search results
	Message string         `json:"message"` // Status message
	Error   string         `json:"error,omitempty"` // Error message if any
}

// Feed represents a single feed item
type Feed struct {
	ID      string `json:"id"`      // Feed ID
	Title   string `json:"title"`   // Feed title
	Content string `json:"content"` // Feed content
	Author  string `json:"author"`  // Author name
	URL     string `json:"url"`     // Feed URL
}

// SearchResult represents a single search result
type SearchResult struct {
	ID      string `json:"id"`      // Result ID
	Title   string `json:"title"`   // Result title
	Content string `json:"content"` // Result content
	Author  string `json:"author"`  // Author name
	URL     string `json:"url"`     // Result URL
}

// Client handles communication with Xiaohongshu MCP server
type Client struct {
	cfg        *config.Config
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Xiaohongshu client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:     cfg,
		baseURL: cfg.Xiaohongshu.ServerURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// MCPRequest represents a generic MCP request
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      string      `json:"id"`
}

// MCPResponse represents a generic MCP response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *MCPError   `json:"error,omitempty"`
	ID      string      `json:"id"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PostContent posts content to Xiaohongshu
func (c *Client) PostContent(req *PostRequest) (*PostResponse, error) {
	// Create arguments map with headless parameter
	arguments := map[string]interface{}{
		"title":    req.Title,
		"content":  req.Content,
		"images":   req.Images,
		"headless": c.cfg.Xiaohongshu.Headless,
	}

	// Debug: Log the arguments being sent
	argumentsJSON, _ := json.MarshalIndent(arguments, "", "  ")
	log.Printf("[DEBUG] PostContent arguments: %s", string(argumentsJSON))

	mcpReq := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "publish_content",
			"arguments": arguments,
		},
		ID: fmt.Sprintf("post_%d", time.Now().UnixNano()),
	}

	// Debug: Log the full MCP request
	mcpReqJSON, _ := json.MarshalIndent(mcpReq, "", "  ")
	log.Printf("[DEBUG] MCP Request: %s", string(mcpReqJSON))

	mcpResp, err := c.callMCP(mcpReq)
	if err != nil {
		return nil, err
	}

	// Convert result to PostResponse
	resultBytes, err := json.Marshal(mcpResp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MCP result: %w", err)
	}

	var result PostResponse
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal post response: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("posting error: %s", result.Error)
	}

	result.PostedAt = time.Now()
	return &result, nil
}

// SchedulePost schedules a post for later publishing
// Note: Scheduling is handled by the MCP server, this method posts immediately
func (c *Client) SchedulePost(req *PostRequest, scheduleTime time.Time) (*PostResponse, error) {
	// For now, just post immediately as the MCP server doesn't support scheduling
	return c.PostContent(req)
}

// GetPostStatus gets the status of a posted content
func (c *Client) GetPostStatus(postID string) (*PostResponse, error) {
	mcpReq := MCPRequest{
		Method: "get_post_status",
		Params: map[string]string{"post_id": postID},
		ID:     fmt.Sprintf("status_%d", time.Now().UnixNano()),
	}

	reqBody, err := json.Marshal(mcpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal status request: %w", err)
	}

	apiURL := fmt.Sprintf("%s/mcp", c.baseURL)
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create status request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get post status: %w", err)
	}
	defer resp.Body.Close()

	var mcpResp MCPResponse
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	if mcpResp.Error != nil {
		return nil, fmt.Errorf("status error %d: %s", mcpResp.Error.Code, mcpResp.Error.Message)
	}

	resultBytes, err := json.Marshal(mcpResp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal status result: %w", err)
	}

	var postResp PostResponse
	if err := json.Unmarshal(resultBytes, &postResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status response: %w", err)
	}

	return &postResp, nil
}

// CheckLoginStatus checks the login status
func (c *Client) CheckLoginStatus() (*LoginStatusResponse, error) {
	// Create arguments map with headless parameter
	arguments := map[string]interface{}{
		"headless": c.cfg.Xiaohongshu.Headless,
	}

	mcpReq := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "check_login_status",
			"arguments": arguments,
		},
		ID: fmt.Sprintf("login_%d", time.Now().UnixNano()),
	}

	mcpResp, err := c.callMCP(mcpReq)
	if err != nil {
		return nil, err
	}

	// Convert result to LoginStatusResponse
	resultBytes, err := json.Marshal(mcpResp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MCP result: %w", err)
	}

	var result LoginStatusResponse
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal login status response: %w", err)
	}

	return &result, nil
}

// ListFeeds gets the feeds list
func (c *Client) ListFeeds() (*FeedsResponse, error) {
	// Create arguments map with headless parameter
	arguments := map[string]interface{}{
		"headless": c.cfg.Xiaohongshu.Headless,
	}

	mcpReq := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "list_feeds",
			"arguments": arguments,
		},
		ID: fmt.Sprintf("feeds_%d", time.Now().UnixNano()),
	}

	mcpResp, err := c.callMCP(mcpReq)
	if err != nil {
		return nil, err
	}

	// Convert result to FeedsResponse
	resultBytes, err := json.Marshal(mcpResp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MCP result: %w", err)
	}

	var result FeedsResponse
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feeds response: %w", err)
	}

	return &result, nil
}

// SearchFeeds searches for content
func (c *Client) SearchFeeds(keyword string) (*SearchResponse, error) {
	// Create arguments map with headless parameter
	arguments := map[string]interface{}{
		"keyword":  keyword,
		"headless": c.cfg.Xiaohongshu.Headless,
	}

	mcpReq := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "search_feeds",
			"arguments": arguments,
		},
		ID: fmt.Sprintf("search_%d", time.Now().UnixNano()),
	}

	mcpResp, err := c.callMCP(mcpReq)
	if err != nil {
		return nil, err
	}

	// Convert result to SearchResponse
	resultBytes, err := json.Marshal(mcpResp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MCP result: %w", err)
	}

	var result SearchResponse
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %w", err)
	}

	return &result, nil
}

// callMCP is a helper method for making MCP calls
func (c *Client) callMCP(mcpReq MCPRequest) (*MCPResponse, error) {
	// Convert to JSON
	reqBody, err := json.Marshal(mcpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MCP request: %w", err)
	}

	// Debug: Log the raw request body
	log.Printf("[DEBUG] Raw MCP request body: %s", string(reqBody))

	// Make HTTP request to MCP server
	apiURL := fmt.Sprintf("%s/mcp", c.baseURL)
	log.Printf("[DEBUG] Sending request to: %s", apiURL)

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Xiaohongshu MCP server: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[DEBUG] MCP server response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Xiaohongshu MCP server returned status %d", resp.StatusCode)
	}

	// Parse MCP response
	var mcpResp MCPResponse
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		return nil, fmt.Errorf("failed to decode MCP response: %w", err)
	}

	// Debug: Log the MCP response
	mcpRespJSON, _ := json.MarshalIndent(mcpResp, "", "  ")
	log.Printf("[DEBUG] MCP Response: %s", string(mcpRespJSON))

	if mcpResp.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", mcpResp.Error.Code, mcpResp.Error.Message)
	}

	return &mcpResp, nil
}

// TestConnection tests the connection to Xiaohongshu MCP server
func (c *Client) TestConnection() error {
	// Test connection by checking login status
	_, err := c.CheckLoginStatus()
	return err
}

// GetAccountInfo gets account information
func (c *Client) GetAccountInfo() (map[string]interface{}, error) {
	mcpReq := MCPRequest{
		Method: "get_account_info",
		Params: nil,
		ID:     fmt.Sprintf("account_%d", time.Now().UnixNano()),
	}

	reqBody, err := json.Marshal(mcpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal account request: %w", err)
	}

	apiURL := fmt.Sprintf("%s/mcp", c.baseURL)
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create account request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}
	defer resp.Body.Close()

	var mcpResp MCPResponse
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		return nil, fmt.Errorf("failed to decode account response: %w", err)
	}

	if mcpResp.Error != nil {
		return nil, fmt.Errorf("account error %d: %s", mcpResp.Error.Code, mcpResp.Error.Message)
	}

	if accountInfo, ok := mcpResp.Result.(map[string]interface{}); ok {
		return accountInfo, nil
	}

	return nil, fmt.Errorf("invalid account info response format")
}

// PostWithRetry posts content with retry logic
func (c *Client) PostWithRetry(req *PostRequest, maxRetries int) (*PostResponse, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		resp, err := c.PostContent(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if i < maxRetries-1 {
			// Wait before retry
			time.Sleep(time.Duration(i+1) * 10 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to post content after %d retries: %w", maxRetries, lastErr)
}

// ValidatePostRequest validates a post request
func (c *Client) ValidatePostRequest(req *PostRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}
	if len(req.Images) == 0 {
		return fmt.Errorf("at least one image is required")
	}
	if len(req.Title) > 100 {
		return fmt.Errorf("title too long (max 100 characters)")
	}
	if len(req.Content) > 1000 {
		return fmt.Errorf("content too long (max 1000 characters)")
	}

	return nil
}