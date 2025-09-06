package visitor

import (
	"fmt"
	"math/rand"
	"time"

	"xiaohongshu-unified/internal/config"
)

// VisitorInfo represents visitor account information
type VisitorInfo struct {
	Date         time.Time `json:"date"`
	VisitorCount int       `json:"visitor_count"`
	DayType      string    `json:"day_type"`
	Description  string    `json:"description"`
	Timestamp    time.Time `json:"timestamp"`
}

// Service handles visitor information
type Service struct {
	config *config.Config
	rand   *rand.Rand
}

// NewService creates a new visitor service
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetVisitorInfo returns visitor information for today
func (s *Service) GetVisitorInfo() (*VisitorInfo, error) {
	return s.GetVisitorInfoForDate(time.Now())
}

// GetVisitorInfoForDate returns visitor information for a specific date
func (s *Service) GetVisitorInfoForDate(date time.Time) (*VisitorInfo, error) {
	dayType, minVisitors, maxVisitors := s.getDayTypeAndRange(date)
	
	// Generate random visitor count within the range
	visitorCount := minVisitors + s.rand.Intn(maxVisitors-minVisitors+1)
	
	visitorInfo := &VisitorInfo{
		Date:         date,
		VisitorCount: visitorCount,
		DayType:      dayType,
		Description:  s.generateDescription(dayType, visitorCount),
		Timestamp:    time.Now(),
	}
	
	return visitorInfo, nil
}

// getDayTypeAndRange determines the day type and visitor count range
func (s *Service) getDayTypeAndRange(date time.Time) (string, int, int) {
	month := int(date.Month())
	day := date.Day()
	weekday := date.Weekday()
	
	// Check for National Day Holiday (Oct 1-7)
	if month == 10 && day >= 1 && day <= 7 {
		return "国庆节假期", 26000, 35000
	}
	
	// Check for May Holiday (May 1-7)
	if month == 5 && day >= 1 && day <= 7 {
		return "五一假期", 26000, 35000
	}
	
	// Check for Summer Holiday (July and August)
	if month == 7 || month == 8 {
		return "暑假", 25000, 30000
	}
	
	// Check for Winter Holiday (January 15-31, February 1-15)
	if (month == 1 && day >= 15) || (month == 2 && day <= 15) {
		return "寒假", 20000, 25000
	}
	
	// Weekend (Saturday and Sunday)
	if weekday == time.Saturday || weekday == time.Sunday {
		return "周末", 15000, 21000
	}
	
	// Regular weekday
	return "工作日", 12000, 17000
}

// generateDescription creates a human-readable description
func (s *Service) generateDescription(dayType string, visitorCount int) string {
	var level string
	switch {
	case visitorCount >= 30000:
		level = "极高"
	case visitorCount >= 25000:
		level = "很高"
	case visitorCount >= 20000:
		level = "高"
	case visitorCount >= 15000:
		level = "较高"
	default:
		level = "正常"
	}
	
	return fmt.Sprintf("%s，预计游客量%d人，人流量%s", dayType, visitorCount, level)
}

// GetFormattedVisitor returns a formatted string representation
func (v *VisitorInfo) GetFormattedVisitor() string {
	return fmt.Sprintf("👥 游客量：%s", v.Description)
}