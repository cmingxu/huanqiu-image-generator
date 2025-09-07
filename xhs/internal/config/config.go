package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	// Weather API configuration
	WeatherAPI struct {
		APIKey  string `json:"api_key"`
		BaseURL string `json:"base_url"`
		City    string `json:"city"`
	} `json:"weather_api"`

	// Traffic API configuration
	TrafficAPI struct {
		APIKey  string `json:"api_key"`
		BaseURL string `json:"base_url"`
		City    string `json:"city"`
	} `json:"traffic_api"`

	// DeepSeek LLM configuration
	DeepSeekLLM struct {
		APIKey  string `json:"api_key"`
		BaseURL string `json:"base_url"`
		Model   string `json:"model"`
	} `json:"deepseek_llm"`

	// MCP configuration for cover generation
	MCP struct {
		ServerURL string `json:"server_url"`
		Headless  bool   `json:"headless"`
		BaseURL   string `json:"base_url"`
		OutDir    string `json:"out_dir"`
	} `json:"mcp"`

	// Xiaohongshu configuration
	Xiaohongshu struct {
		ServerURL string `json:"server_url"`
		Headless  bool   `json:"headless"`
	} `json:"xiaohongshu"`

	// Weibo configuration
	Weibo struct {
		UID     string `json:"uid"`
		Cookies string `json:"cookies"`
		Token   string `json:"token"`
	} `json:"weibo"`

	// General settings
	Settings struct {
		PostInterval string `json:"post_interval"` // e.g., "1h", "24h"
		LogLevel     string `json:"log_level"`
	} `json:"settings"`
}

// Load loads configuration from config.json file or environment variables
func Load() (*Config, error) {
	cfg := &Config{}

	// Try to load from config.json first
	if data, err := os.ReadFile("config.json"); err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config.json: %w", err)
		}
	}

	// Override with environment variables if present
	loadFromEnv(cfg)

	// Set defaults
	setDefaults(cfg)

	// Validate configuration
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	if key := os.Getenv("WEATHER_API_KEY"); key != "" {
		cfg.WeatherAPI.APIKey = key
	}
	if url := os.Getenv("WEATHER_API_URL"); url != "" {
		cfg.WeatherAPI.BaseURL = url
	}
	if city := os.Getenv("CITY"); city != "" {
		cfg.WeatherAPI.City = city
		cfg.TrafficAPI.City = city
	}

	if key := os.Getenv("TRAFFIC_API_KEY"); key != "" {
		cfg.TrafficAPI.APIKey = key
	}
	if url := os.Getenv("TRAFFIC_API_URL"); url != "" {
		cfg.TrafficAPI.BaseURL = url
	}

	if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
		cfg.DeepSeekLLM.APIKey = key
	}
	if url := os.Getenv("DEEPSEEK_API_URL"); url != "" {
		cfg.DeepSeekLLM.BaseURL = url
	}
	if model := os.Getenv("DEEPSEEK_MODEL"); model != "" {
		cfg.DeepSeekLLM.Model = model
	}

	if url := os.Getenv("MCP_SERVER_URL"); url != "" {
		cfg.MCP.ServerURL = url
	}

	if url := os.Getenv("XHS_SERVER_URL"); url != "" {
		cfg.Xiaohongshu.ServerURL = url
	}

	if uid := os.Getenv("WEIBO_UID"); uid != "" {
		cfg.Weibo.UID = uid
	}
	if cookies := os.Getenv("WEIBO_COOKIES"); cookies != "" {
		cfg.Weibo.Cookies = cookies
	}
	if token := os.Getenv("WEIBO_TOKEN"); token != "" {
		cfg.Weibo.Token = token
	}

	if interval := os.Getenv("POST_INTERVAL"); interval != "" {
		cfg.Settings.PostInterval = interval
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Settings.LogLevel = logLevel
	}
}

// setDefaults sets default values for configuration
func setDefaults(cfg *Config) {
	if cfg.WeatherAPI.BaseURL == "" {
		cfg.WeatherAPI.BaseURL = "https://api.openweathermap.org/data/2.5"
	}
	if cfg.WeatherAPI.City == "" {
		cfg.WeatherAPI.City = "Beijing"
	}

	if cfg.TrafficAPI.City == "" {
		cfg.TrafficAPI.City = "Beijing"
	}

	if cfg.DeepSeekLLM.BaseURL == "" {
		cfg.DeepSeekLLM.BaseURL = "https://api.deepseek.com"
	}
	if cfg.DeepSeekLLM.Model == "" {
		cfg.DeepSeekLLM.Model = "deepseek-chat"
	}

	if cfg.MCP.ServerURL == "" {
		cfg.MCP.ServerURL = "http://localhost:18062"
	}
	if cfg.MCP.BaseURL == "" {
		cfg.MCP.BaseURL = "http://localhost:3000"
	}
	if cfg.MCP.OutDir == "" {
		cfg.MCP.OutDir = "/Users/kx/Desktop"
	}
	// Default headless to false for MCP
	cfg.MCP.Headless = false

	if cfg.Xiaohongshu.ServerURL == "" {
		cfg.Xiaohongshu.ServerURL = "http://localhost:18062"
	}
	// Default headless to false for Xiaohongshu
	cfg.Xiaohongshu.Headless = false

	// Set default weibo configuration
	if cfg.Weibo.UID == "" {
		cfg.Weibo.UID = "3937775216"
	}
	if cfg.Weibo.Cookies == "" {
		cfg.Weibo.Cookies = "SUB=_2AkMfzfZxf8NxqwFRmfscymjibox_zA3EieKpkQeqJRMxHRl-yT9kqnIitRB6NE3Ynp3g3XUjfERDfRvDu2Ob-V0AV-Ht; XSRF-TOKEN=JVS9su9p3gsRZyDzgsijAdx5; WBPSESS=gJ7ElPMf_3q2cdj5JUfmvNSXzQofuuhpbfKWU-JmetuhhFVlp1s7T3D6PJClzn45urDFp34oVajUL4N7sYweJyZs74npFsMnIJ9PUcbSjV9Pwg5IdiwWIEUuHTqDSRsJ3pCe78X7Zm38ENkYYoFzAwkxKSCkNQ3Kb-j9COTqz14="
	}
	if cfg.Weibo.Token == "" {
		cfg.Weibo.Token = "JVS9su9p3gsRZyDzgsijAdx5"
	}

	if cfg.Settings.PostInterval == "" {
		cfg.Settings.PostInterval = "24h"
	}
	if cfg.Settings.LogLevel == "" {
		cfg.Settings.LogLevel = "info"
	}
}

// validate validates the configuration
func validate(cfg *Config) error {
	if cfg.WeatherAPI.APIKey == "" {
		return fmt.Errorf("weather API key is required")
	}
	if cfg.DeepSeekLLM.APIKey == "" {
		return fmt.Errorf("DeepSeek API key is required")
	}
	return nil
}

