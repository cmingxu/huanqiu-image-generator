package orchestrator

import (
	"fmt"
	"log"
	"time"

	"xiaohongshu-unified/internal/config"
	"xiaohongshu-unified/internal/llm"
	"xiaohongshu-unified/internal/lunar"
	"xiaohongshu-unified/internal/cover-gen"
	"xiaohongshu-unified/internal/traffic"
	"xiaohongshu-unified/internal/visitor"
	"xiaohongshu-unified/internal/weather"
	"xiaohongshu-unified/internal/weibo"
	"xiaohongshu-unified/internal/xhs"
)

// Orchestrator coordinates all services to generate and post content
type Orchestrator struct {
	cfg           *config.Config
	weatherSvc    *weather.Service
	lunarSvc      *lunar.Service
	trafficSvc    *traffic.Service
	visitorSvc    *visitor.Service
	weiboSvc      *weibo.Service
	llmSvc        *llm.Service
	coverMCPClient     *covergen.Client
	xhsClient     *xhs.Client
}

// New creates a new orchestrator
func New(cfg *config.Config) *Orchestrator {
	return &Orchestrator{
		cfg:        cfg,
		weatherSvc: weather.NewService(cfg),
		lunarSvc:   lunar.NewService(),
		trafficSvc: traffic.NewService(cfg),
		visitorSvc: visitor.NewService(cfg),
		weiboSvc:   weibo.NewService(cfg.Weibo.UID, cfg.Weibo.Cookies, cfg.Weibo.Token),
		llmSvc:     llm.NewService(cfg),
		coverMCPClient:  covergen.NewClient(cfg),
		xhsClient:  xhs.NewClient(cfg),
	}
}

// WorkflowResult represents the result of the complete workflow
type WorkflowResult struct {
	WeatherInfo     *weather.WeatherInfo `json:"weather_info"`
	LunarInfo       *lunar.LunarInfo     `json:"lunar_info"`
	TrafficInfo     *traffic.TrafficInfo `json:"traffic_info"`
	VisitorInfo     *visitor.VisitorInfo `json:"visitor_info"`
	WeiboContent    string               `json:"weibo_content"`
	GeneratedContent *llm.GeneratedContent `json:"generated_content"`
	ImageResponse   *covergen.ImageResponse   `json:"image_response"`
	PostResponse    *xhs.PostResponse    `json:"post_response"`
	ExecutionTime   time.Duration        `json:"execution_time"`
	Timestamp       time.Time            `json:"timestamp"`
	Success         bool                 `json:"success"`
	Error           string               `json:"error,omitempty"`
}

// Run executes the complete workflow
func (o *Orchestrator) Run() error {
	start := time.Now()
	result := &WorkflowResult{
		Timestamp: start,
	}

	log.Println("🚀 Starting Xiaohongshu content generation workflow...")

	// Step 1: Test connections
	if err := o.testConnections(); err != nil {
		result.Error = fmt.Sprintf("Connection test failed: %v", err)
		log.Printf("❌ %s", result.Error)
		return fmt.Errorf(result.Error)
	}
	log.Println("✅ All service connections tested successfully")

	// Step 2: Gather information
	log.Println("📊 Gathering information...")

	// Get weather information
	weatherInfo, err := o.weatherSvc.GetWeatherInfo()
	if err != nil {
		log.Printf("⚠️ Failed to get weather info: %v", err)
		// Continue without weather info
	} else {
		result.WeatherInfo = weatherInfo
		log.Printf("🌤️ Weather: %s", weatherInfo.GetFormattedWeather())
	}

	// Get lunar information
	lunarInfo, err := o.lunarSvc.GetLunarInfo()
	if err != nil {
		log.Printf("⚠️ Failed to get lunar info: %v", err)
		// Continue without lunar info
	} else {
		result.LunarInfo = lunarInfo
		log.Printf("📅 Lunar: %s", lunarInfo.GetFormattedLunar())
	}

	// Get traffic information
	trafficInfo, err := o.trafficSvc.GetTrafficInfo()
	if err != nil {
		log.Printf("⚠️ Failed to get traffic info: %v", err)
		// Continue without traffic info
	} else {
		result.TrafficInfo = trafficInfo
		log.Printf("🚗 Traffic: %s overall status", trafficInfo.OverallStatus)
	}

	// Get visitor information
	visitorInfo, err := o.visitorSvc.GetVisitorInfo()
	if err != nil {
		log.Printf("⚠️ Failed to get visitor info: %v", err)
		// Continue without visitor info
	} else {
		result.VisitorInfo = visitorInfo
		log.Printf("👥 Visitor: %s", visitorInfo.GetFormattedVisitor())
	}

	// Get weibo content for summary
	weiboSummaryContent, err := o.weiboSvc.GetTop2PostsForSummary()
	if err != nil {
		log.Printf("Warning: Failed to get weibo summary content: %v", err)
		weiboSummaryContent = "" // Continue without weibo content
	}
	result.WeiboContent = weiboSummaryContent
	if weiboSummaryContent != "" {
		log.Printf("📱 Weibo: Got recent content")
	}

	// Step 3: Generate content using LLM
	log.Println("🤖 Generating content with DeepSeek LLM...")
	contentReq := &llm.ContentRequest{
		Weather: weatherInfo,
		Lunar:   lunarInfo,
		// Traffic: trafficInfo, // Omitted per user request
		Visitor: visitorInfo,
		Weibo:   weiboSummaryContent,
		Theme:   "daily life sharing", // You can make this configurable
	}

	generatedContent, err := o.llmSvc.GenerateContent(contentReq)
	if err != nil {
		result.Error = fmt.Sprintf("Content generation failed: %v", err)
		log.Printf("❌ %s", result.Error)
		return fmt.Errorf(result.Error)
	}
	result.GeneratedContent = generatedContent
	log.Printf("✅ Content generated: %s", generatedContent.Title)

	// Step 4: Generate cover image
	log.Println("🎨 Generating cover image...")
	// Use a default image prompt and the cover_text from LLM response
	defaultImagePrompt := "cozy daily life scene, warm lighting, lifestyle photography, Beijing Universal Studios theme park"
	imageResp, err := o.coverMCPClient.GenerateXiaohongshuCover(
		defaultImagePrompt,
		generatedContent.CoverText,
	)
	if err != nil {
		result.Error = fmt.Sprintf("Image generation failed: %v", err)
		log.Printf("❌ %s", result.Error)
		return fmt.Errorf(result.Error)
	}
	result.ImageResponse = imageResp
	log.Printf("✅ Cover image generated: %s", imageResp.ImageURL)

	// Step 5: Post to Xiaohongshu
	log.Println("📱 Posting to Xiaohongshu...")
	postReq := &xhs.PostRequest{
		Title:   generatedContent.Title,
		Content: generatedContent.GetFormattedContent(),
		Images:  []string{imageResp.ImageURL},
	}

	// Validate post request
	if err := o.xhsClient.ValidatePostRequest(postReq); err != nil {
		result.Error = fmt.Sprintf("Post validation failed: %v", err)
		log.Printf("❌ %s", result.Error)
		return fmt.Errorf(result.Error)
	}

	// Post with retry
	postResp, err := o.xhsClient.PostWithRetry(postReq, 3)
	if err != nil {
		result.Error = fmt.Sprintf("Posting failed: %v", err)
		log.Printf("❌ %s", result.Error)
		return fmt.Errorf(result.Error)
	}
	result.PostResponse = postResp
	log.Printf("✅ Posted successfully: %s", postResp.URL)

	// Step 6: Complete workflow
	result.ExecutionTime = time.Since(start)
	result.Success = true

	log.Printf("🎉 Workflow completed successfully in %v", result.ExecutionTime)
	log.Printf("📝 Post ID: %s", postResp.PostID)
	log.Printf("🔗 Post URL: %s", postResp.URL)

	return nil
}

// testConnections tests all external service connections
func (o *Orchestrator) testConnections() error {
	log.Println("🔍 Testing service connections...")

	// Test MCP server connection
	if err := o.coverMCPClient.TestConnection(); err != nil {
		return fmt.Errorf("MCP server connection failed: %w", err)
	}
	log.Println("✅ MCP server connection OK")

	// Test Xiaohongshu MCP server connection
	if err := o.xhsClient.TestConnection(); err != nil {
		return fmt.Errorf("Xiaohongshu MCP server connection failed: %w", err)
	}
	log.Println("✅ Xiaohongshu MCP server connection OK")

	return nil
}

// RunScheduled runs the workflow on a schedule
func (o *Orchestrator) RunScheduled() error {
	interval, err := time.ParseDuration(o.cfg.Settings.PostInterval)
	if err != nil {
		return fmt.Errorf("invalid post interval: %w", err)
	}

	log.Printf("📅 Starting scheduled workflow with interval: %v", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run once immediately
	if err := o.Run(); err != nil {
		log.Printf("❌ Initial run failed: %v", err)
	}

	// Then run on schedule
	for {
		select {
		case <-ticker.C:
			log.Println("⏰ Scheduled run starting...")
			if err := o.Run(); err != nil {
				log.Printf("❌ Scheduled run failed: %v", err)
				// Continue with next scheduled run
			}
		}
	}
}

// RunOnce runs the workflow once and exits
func (o *Orchestrator) RunOnce() error {
	return o.Run()
}

// GetServiceStatus returns the status of all services
func (o *Orchestrator) GetServiceStatus() map[string]string {
	status := make(map[string]string)

	// Test MCP server
	if err := o.coverMCPClient.TestConnection(); err != nil {
		status["mcp_server"] = fmt.Sprintf("❌ Error: %v", err)
	} else {
		status["mcp_server"] = "✅ OK"
	}

	// Test Xiaohongshu MCP server
	if err := o.xhsClient.TestConnection(); err != nil {
		status["xiaohongshu_server"] = fmt.Sprintf("❌ Error: %v", err)
	} else {
		status["xiaohongshu_server"] = "✅ OK"
	}

	// Test weather service (try to get weather info)
	if _, err := o.weatherSvc.GetWeatherInfo(); err != nil {
		status["weather_service"] = fmt.Sprintf("❌ Error: %v", err)
	} else {
		status["weather_service"] = "✅ OK"
	}

	// Lunar service (always available as it uses mock data)
	status["lunar_service"] = "✅ OK"

	// Traffic service (always available as it uses mock data)
	status["traffic_service"] = "✅ OK"

	// Visitor service (always available as it uses mock data)
	status["visitor_service"] = "✅ OK"

	return status
}