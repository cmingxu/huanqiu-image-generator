package traffic

import (
	"fmt"
	"net/http"
	"time"

	"xiaohongshu-unified/internal/config"
)

// TrafficInfo represents traffic information
type TrafficInfo struct {
	City           string                 `json:"city"`
	OverallStatus  string                 `json:"overall_status"`  // 整体路况：畅通、缓行、拥堵
	CongestionLevel int                   `json:"congestion_level"` // 拥堵等级 1-10
	MainRoads      []RoadInfo            `json:"main_roads"`      // 主要道路信息
	Incidents      []TrafficIncident     `json:"incidents"`       // 交通事件
	Recommendation string                 `json:"recommendation"`  // 出行建议
	Timestamp      time.Time             `json:"timestamp"`
}

// RoadInfo represents information about a specific road
type RoadInfo struct {
	Name        string `json:"name"`        // 道路名称
	Status      string `json:"status"`      // 路况状态
	Speed       int    `json:"speed"`       // 平均速度 km/h
	TravelTime  string `json:"travel_time"` // 通行时间
	Description string `json:"description"` // 详细描述
}

// TrafficIncident represents a traffic incident
type TrafficIncident struct {
	Type        string `json:"type"`        // 事件类型：事故、施工、管制等
	Location    string `json:"location"`    // 事件位置
	Description string `json:"description"` // 事件描述
	Severity    string `json:"severity"`    // 严重程度：轻微、一般、严重
	StartTime   string `json:"start_time"`  // 开始时间
	EstimatedEnd string `json:"estimated_end"` // 预计结束时间
}

// Service handles traffic information fetching
type Service struct {
	cfg    *config.Config
	client *http.Client
}

// NewService creates a new traffic service
func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetTrafficInfo fetches current traffic information
func (s *Service) GetTrafficInfo() (*TrafficInfo, error) {
	// For demonstration, we'll generate mock traffic data
	// In production, you would integrate with real traffic APIs like:
	// - 高德地图 API
	// - 百度地图 API
	// - Google Maps API
	// - 腾讯地图 API

	return s.generateMockTrafficInfo(), nil

	// Uncomment below for real API integration
	/*
	apiURL := fmt.Sprintf("%s/traffic", s.cfg.TrafficAPI.BaseURL)
	params := url.Values{}
	params.Add("city", s.cfg.TrafficAPI.City)
	params.Add("key", s.cfg.TrafficAPI.APIKey)

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	resp, err := s.client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch traffic data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("traffic API returned status %d", resp.StatusCode)
	}

	// Parse response based on your chosen API
	var trafficInfo TrafficInfo
	if err := json.NewDecoder(resp.Body).Decode(&trafficInfo); err != nil {
		return nil, fmt.Errorf("failed to decode traffic response: %w", err)
	}

	return &trafficInfo, nil
	*/
}

// generateMockTrafficInfo generates mock traffic data
func (s *Service) generateMockTrafficInfo() *TrafficInfo {
	now := time.Now()
	hour := now.Hour()

	// Determine overall status based on time of day
	var overallStatus string
	var congestionLevel int
	var recommendation string

	switch {
	case hour >= 7 && hour <= 9: // Morning rush hour
		overallStatus = "拥堵"
		congestionLevel = 8
		recommendation = "早高峰时段，建议错峰出行或选择公共交通"
	case hour >= 17 && hour <= 19: // Evening rush hour
		overallStatus = "拥堵"
		congestionLevel = 9
		recommendation = "晚高峰时段，道路拥堵严重，建议延后出行"
	case hour >= 10 && hour <= 16: // Daytime
		overallStatus = "缓行"
		congestionLevel = 4
		recommendation = "白天时段，整体路况良好，适合出行"
	case hour >= 20 || hour <= 6: // Night time
		overallStatus = "畅通"
		congestionLevel = 2
		recommendation = "夜间时段，道路畅通，出行便利"
	default:
		overallStatus = "缓行"
		congestionLevel = 5
		recommendation = "路况一般，注意安全驾驶"
	}

	// Generate main roads info
	mainRoads := []RoadInfo{
		{
			Name:        "三环路",
			Status:      getStatusByLevel(congestionLevel),
			Speed:       getSpeedByLevel(congestionLevel),
			TravelTime:  "45-60分钟",
			Description: "主要环路，车流量较大",
		},
		{
			Name:        "长安街",
			Status:      getStatusByLevel(congestionLevel - 1),
			Speed:       getSpeedByLevel(congestionLevel - 1),
			TravelTime:  "30-40分钟",
			Description: "东西主干道，通行状况良好",
		},
		{
			Name:        "京藏高速",
			Status:      getStatusByLevel(congestionLevel + 1),
			Speed:       getSpeedByLevel(congestionLevel + 1),
			TravelTime:  "60-90分钟",
			Description: "进出京主要通道，易发生拥堵",
		},
		{
			Name:        "中关村大街",
			Status:      getStatusByLevel(congestionLevel),
			Speed:       getSpeedByLevel(congestionLevel),
			TravelTime:  "25-35分钟",
			Description: "科技园区主干道，上下班时段较拥堵",
		},
	}

	// Generate incidents based on congestion level
	var incidents []TrafficIncident
	if congestionLevel > 6 {
		incidents = []TrafficIncident{
			{
				Type:        "交通事故",
				Location:    "三环路国贸桥附近",
				Description: "两车追尾，占用一条车道",
				Severity:    "一般",
				StartTime:   now.Add(-30 * time.Minute).Format("15:04"),
				EstimatedEnd: now.Add(20 * time.Minute).Format("15:04"),
			},
			{
				Type:        "道路施工",
				Location:    "京藏高速清河收费站",
				Description: "路面维修，限制通行",
				Severity:    "轻微",
				StartTime:   "09:00",
				EstimatedEnd: "17:00",
			},
		}
	} else if congestionLevel > 3 {
		incidents = []TrafficIncident{
			{
				Type:        "交通管制",
				Location:    "天安门广场周边",
				Description: "临时交通管制，请绕行",
				Severity:    "轻微",
				StartTime:   "08:00",
				EstimatedEnd: "18:00",
			},
		}
	}

	return &TrafficInfo{
		City:           s.cfg.TrafficAPI.City,
		OverallStatus:  overallStatus,
		CongestionLevel: congestionLevel,
		MainRoads:      mainRoads,
		Incidents:      incidents,
		Recommendation: recommendation,
		Timestamp:      now,
	}
}

// getStatusByLevel returns status string based on congestion level
func getStatusByLevel(level int) string {
	if level < 0 {
		level = 0
	}
	if level > 10 {
		level = 10
	}

	switch {
	case level <= 3:
		return "畅通"
	case level <= 6:
		return "缓行"
	default:
		return "拥堵"
	}
}

// getSpeedByLevel returns average speed based on congestion level
func getSpeedByLevel(level int) int {
	if level < 0 {
		level = 0
	}
	if level > 10 {
		level = 10
	}

	// Speed decreases as congestion level increases
	// Max speed: 60 km/h, Min speed: 10 km/h
	return 60 - (level * 5)
}

// GetFormattedTraffic returns traffic information in a human-readable format
func (t *TrafficInfo) GetFormattedTraffic() string {
	result := fmt.Sprintf("🚗 %s交通：整体%s（拥堵等级%d/10）\n", t.City, t.OverallStatus, t.CongestionLevel)
	result += fmt.Sprintf("💡 出行建议：%s\n", t.Recommendation)

	if len(t.MainRoads) > 0 {
		result += "\n🛣️ 主要道路：\n"
		for _, road := range t.MainRoads {
			result += fmt.Sprintf("• %s：%s（平均%dkm/h，预计%s）\n", 
				road.Name, road.Status, road.Speed, road.TravelTime)
		}
	}

	if len(t.Incidents) > 0 {
		result += "\n⚠️ 交通事件：\n"
		for _, incident := range t.Incidents {
			result += fmt.Sprintf("• %s：%s（%s，%s开始）\n", 
				incident.Type, incident.Location, incident.Description, incident.StartTime)
		}
	}

	return result
}