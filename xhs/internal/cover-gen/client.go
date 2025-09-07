package covergen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	"xiaohongshu-unified/internal/config"
)

// ImageRequest represents a request to generate an image
type ImageRequest struct {
	Prompt      string            `json:"prompt"`       // Image generation prompt
	Style       string            `json:"style,omitempty"` // Image style (optional)
	Size        string            `json:"size,omitempty"`  // Image size (optional)
	Quality     string            `json:"quality,omitempty"` // Image quality (optional)
	Parameters  map[string]interface{} `json:"parameters,omitempty"` // Additional parameters
}

// ImageResponse represents the response from image generation
type ImageResponse struct {
	ImageURL    string `json:"image_url"`    // Generated image URL
	ImagePath   string `json:"image_path"`   // Local image path
	ImageData   string `json:"image_data"`   // Base64 encoded image data
	Prompt      string `json:"prompt"`       // Used prompt
	GeneratedAt time.Time `json:"generated_at"`
	Error       string `json:"error,omitempty"` // Error message if any
}

// Client handles communication with MCP server
type Client struct {
	cfg        *config.Config
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new MCP client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:     cfg,
		baseURL: cfg.MCP.ServerURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Image generation can take time
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

// MCPResponse represents the response from MCP server
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

// AssetsResponse represents the response from the assets API
type AssetsResponse struct {
	Images []string `json:"images"`
}

// fetchAvailableAssets fetches the list of available assets from the cover service
func (c *Client) fetchAvailableAssets() ([]string, error) {
	url := c.cfg.MCP.BaseURL + "/api/assets"
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch assets: status %d", resp.StatusCode)
	}

	var assetsResp AssetsResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetsResp); err != nil {
		return nil, fmt.Errorf("failed to decode assets response: %w", err)
	}

	return assetsResp.Images, nil
}

// selectRandomAsset selects a random asset from the available list
func (c *Client) selectRandomAsset() (string, error) {
	assets, err := c.fetchAvailableAssets()
	if err != nil {
		return "", err
	}

	if len(assets) == 0 {
		return "", fmt.Errorf("no assets available")
	}

	// Select a random asset
	index := rand.Intn(len(assets))
	return assets[index], nil
}

// generateOutputPath generates a timestamped output path in the configured directory
func (c *Client) generateOutputPath() string {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("cover_%s.jpeg", timestamp)
	return filepath.Join(c.cfg.MCP.OutDir, filename)
}

// GenerateImage generates an image using the MCP server
func (c *Client) GenerateImage(req *ImageRequest) (*ImageResponse, error) {
	// For now, delegate to GenerateXiaohongshuCover with the prompt
	// This maintains backward compatibility while using the new MCP format
	return c.GenerateXiaohongshuCover(req.Prompt, "Generated Image")
}

// GenerateXiaohongshuCover generates a cover image specifically for Xiaohongshu
func (c *Client) GenerateXiaohongshuCover(prompt, coverText string) (*ImageResponse, error) {
	// Use the coverText parameter as the overlay text
	if coverText == "" {
		coverText = "Sample Text"
	}

	// Select a random asset image
	randomAsset, err := c.selectRandomAsset()
	if err != nil {
		return nil, fmt.Errorf("failed to select random asset: %w", err)
	}

	// Generate timestamped output path
	outputPath := c.generateOutputPath()

	// Prepare MCP request using tools/call method with proper parameters
	mcpReq := MCPRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name": "generate_xiaohongshu_cover",
			"arguments": map[string]interface{}{
				"baseUrl":         c.cfg.MCP.BaseURL,
				"selector":        "#exportable",
				"image":           randomAsset,
				"text":            coverText,
				"output_path":     outputPath,
				"headless":        c.cfg.MCP.Headless,
				"fontFamily":      "Comic Sans MS",
				"fontSize":        48,
				"fontWeight":      "bold",
				"color":           "#0e0d0c",
				"backgroundColor": "#f4f750",
				"textShadow":      "2px 2px 4px #000000",
				"border":          "1px solidrgb(187, 23, 23)",
				"borderRadius":    32,
				"borderWidth":     2,
				"borderStyle":     "dashed",
				"padding":         40,
				"scaleX":          1.0,
				"scaleY":          1.0,
				"skewX":           -15,
				"skewY":           0.0,
				"opacity":         0.8,
				"overlayColor":    "#443c3c",
				"x":               50,
				"y":               50,
			},
		},
		ID: fmt.Sprintf("cover_%d", time.Now().UnixNano()),
	}

	// Convert to JSON
	reqBody, err := json.Marshal(mcpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MCP request: %w", err)
	}

	// Make HTTP request to MCP server
	apiURL := c.baseURL
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call MCP server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MCP server returned status %d", resp.StatusCode)
	}

	// Parse MCP response
	var mcpResp MCPResponse
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		return nil, fmt.Errorf("failed to decode MCP response: %w", err)
	}

	if mcpResp.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", mcpResp.Error.Code, mcpResp.Error.Message)
	}

	// Create ImageResponse with the generated image path
	imageResp := &ImageResponse{
		ImagePath:   outputPath,
		ImageURL:    outputPath,
		Prompt:      coverText,
		GeneratedAt: time.Now(),
	}

	return imageResp, nil
}

// TestConnection tests the connection to MCP server
func (c *Client) TestConnection() error {
	// Try to ping the MCP server
	apiURL := fmt.Sprintf("%s/health", c.baseURL)
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to connect to MCP server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MCP server health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetCapabilities gets the capabilities of the MCP server
func (c *Client) GetCapabilities() (map[string]interface{}, error) {
	mcpReq := MCPRequest{
		Method: "get_capabilities",
		Params: nil,
		ID:     fmt.Sprintf("cap_%d", time.Now().UnixNano()),
	}

	reqBody, err := json.Marshal(mcpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capabilities request: %w", err)
	}

	apiURL := fmt.Sprintf("%s/mcp", c.baseURL)
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create capabilities request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get capabilities: %w", err)
	}
	defer resp.Body.Close()

	var mcpResp MCPResponse
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		return nil, fmt.Errorf("failed to decode capabilities response: %w", err)
	}

	if mcpResp.Error != nil {
		return nil, fmt.Errorf("capabilities error %d: %s", mcpResp.Error.Code, mcpResp.Error.Message)
	}

	if capabilities, ok := mcpResp.Result.(map[string]interface{}); ok {
		return capabilities, nil
	}

	return nil, fmt.Errorf("invalid capabilities response format")
}

// GenerateImageWithRetry generates an image with retry logic
func (c *Client) GenerateImageWithRetry(req *ImageRequest, maxRetries int) (*ImageResponse, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		resp, err := c.GenerateImage(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if i < maxRetries-1 {
			// Wait before retry
			time.Sleep(time.Duration(i+1) * 5 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to generate image after %d retries: %w", maxRetries, lastErr)
}