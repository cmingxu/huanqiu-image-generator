package weather

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"xiaohongshu-unified/internal/config"
)

// WeatherInfo represents weather information
type WeatherInfo struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
	Humidity    int     `json:"humidity"`
	Description string  `json:"description"`
	WindSpeed   float64 `json:"wind_speed"`
	Visibility  int     `json:"visibility"`
	UVIndex     float64 `json:"uv_index"`
	Timestamp   time.Time `json:"timestamp"`
}

// Service handles weather information fetching
type Service struct {
	cfg    *config.Config
	client *http.Client
}

// NewService creates a new weather service
func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}



// GetWeatherInfo fetches current weather information by scraping Chinese weather website
func (s *Service) GetWeatherInfo() (*WeatherInfo, error) {
	// Use the Chinese weather website URL
	weatherURL := "https://e.weather.com.cn/mweather/101010100.shtml"

	resp, err := s.client.Get(weatherURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather website returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	htmlContent := string(body)

	// Extract weather information using regex patterns
	weatherInfo, err := s.parseWeatherFromHTML(htmlContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}

	return weatherInfo, nil
}

// parseWeatherFromHTML extracts weather information from HTML content
func (s *Service) parseWeatherFromHTML(htmlContent string) (*WeatherInfo, error) {
	// Extract temperature using regex
	tempRegex := regexp.MustCompile(`<span[^>]*class="[^"]*temp[^"]*"[^>]*>([+-]?\d+)°?</span>`)
	tempMatches := tempRegex.FindStringSubmatch(htmlContent)
	var temperature float64 = 20.0 // default
	if len(tempMatches) > 1 {
		if temp, err := strconv.ParseFloat(tempMatches[1], 64); err == nil {
			temperature = temp
		}
	}

	// Extract weather description
	descRegex := regexp.MustCompile(`<span[^>]*class="[^"]*weather[^"]*"[^>]*>([^<]+)</span>`)
	descMatches := descRegex.FindStringSubmatch(htmlContent)
	description := "晴"
	if len(descMatches) > 1 {
		description = strings.TrimSpace(descMatches[1])
	}

	// Extract humidity
	humidityRegex := regexp.MustCompile(`湿度[：:]?\s*(\d+)%`)
	humidityMatches := humidityRegex.FindStringSubmatch(htmlContent)
	var humidity int = 60 // default
	if len(humidityMatches) > 1 {
		if h, err := strconv.Atoi(humidityMatches[1]); err == nil {
			humidity = h
		}
	}

	// Extract wind information
	windRegex := regexp.MustCompile(`风[力速][：:]?\s*(\d+)[级m/s]`)
	windMatches := windRegex.FindStringSubmatch(htmlContent)
	var windSpeed float64 = 3.0 // default
	if len(windMatches) > 1 {
		if w, err := strconv.ParseFloat(windMatches[1], 64); err == nil {
			windSpeed = w
		}
	}

	weatherInfo := &WeatherInfo{
		City:        "北京",
		Temperature: temperature,
		FeelsLike:   temperature + 1.0, // approximate feels like
		Humidity:    humidity,
		Description: description,
		WindSpeed:   windSpeed,
		Visibility:  10000, // default 10km
		UVIndex:     5.0,   // default moderate
		Timestamp:   time.Now(),
	}

	return weatherInfo, nil
}

// GetFormattedWeather returns weather information in a human-readable format
func (w *WeatherInfo) GetFormattedWeather() string {
	return fmt.Sprintf(
		"🌤️ %s天气：%s，气温%.1f°C（体感%.1f°C），湿度%d%%，风速%.1fm/s，能见度%dm，紫外线指数%.1f",
		w.City,
		w.Description,
		w.Temperature,
		w.FeelsLike,
		w.Humidity,
		w.WindSpeed,
		w.Visibility,
		w.UVIndex,
	)
}