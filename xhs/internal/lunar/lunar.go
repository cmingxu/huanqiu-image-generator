package lunar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LunarInfo represents Chinese lunar calendar information
type LunarInfo struct {
	Date         time.Time `json:"date"`
	LunarDate    string    `json:"lunar_date"`     // 农历日期，如"腊月初八"
	LunarYear    string    `json:"lunar_year"`     // 农历年份，如"癸卯年"
	Zodiac       string    `json:"zodiac"`         // 生肖，如"兔"
	SolarTerm    string    `json:"solar_term"`     // 节气，如"大寒"
	Festival     string    `json:"festival"`       // 节日，如"腊八节"
	Suit         []string  `json:"suit"`           // 宜做的事情
	Avoid        []string  `json:"avoid"`          // 忌做的事情
	LuckyColor   string    `json:"lucky_color"`    // 幸运颜色
	LuckyNumber  string    `json:"lucky_number"`   // 幸运数字
	Constellation string   `json:"constellation"`  // 星座
	Timestamp    time.Time `json:"timestamp"`
}

// Service handles lunar calendar information fetching
type Service struct {
	client *http.Client
}

// NewService creates a new lunar service
func NewService() *Service {
	return &Service{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// APIResponse represents the response from lunar calendar API
type APIResponse struct {
	Code int `json:"code"`
	Data struct {
		Date       string `json:"date"`
		DateY      string `json:"dateY"`
		DateD      string `json:"dateD"`
		DateMC     string `json:"dateMC"`
		Jieqi      string `json:"jieqi"`
		Shichen    string `json:"shichen"`
		Xingzuoyunshi struct {
			Yiyan          string `json:"yiyan"`
			Yunshi         int    `json:"yunshi"`
			Xingyuncolor   string `json:"xingyuncolor"`
			Xingyunnumber  int    `json:"xingyunnumber"`
			Supeixingzuo   string `json:"supeixingzuo"`
		} `json:"xingzuoyunshi"`
		Lunardate string `json:"lunardate"`
		Week      string `json:"week"`
		Hseb      string `json:"hseb"`
		Text      string `json:"text"`
		FromWho   string `json:"from_who"`
	} `json:"data"`
}

// GetLunarInfo fetches current lunar calendar information
func (s *Service) GetLunarInfo() (*LunarInfo, error) {
	now := time.Now()
	return s.GetLunarInfoForDate(now)
}

// GetLunarInfoForDate fetches lunar calendar information for a specific date
func (s *Service) GetLunarInfoForDate(date time.Time) (*LunarInfo, error) {
	// Use the new lunar calendar API
	apiURL := "https://api.xcvts.cn/api/huangli"
	
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lunar data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lunar API returned status %d", resp.StatusCode)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode lunar response: %w", err)
	}

	if apiResp.Code != 1 {
		return nil, fmt.Errorf("lunar API error: invalid response code %d", apiResp.Code)
	}

	// Convert the API response to our LunarInfo structure
	lunarInfo := &LunarInfo{
		Date:          date,
		LunarDate:     apiResp.Data.Lunardate,
		LunarYear:     apiResp.Data.Hseb, // 乙巳年 甲申月 戊寅日
		Zodiac:        "", // Not provided in this API
		SolarTerm:     apiResp.Data.Jieqi,
		Festival:      "", // Not directly provided
		Suit:          []string{}, // Not provided in this API
		Avoid:         []string{}, // Not provided in this API
		LuckyColor:    apiResp.Data.Xingzuoyunshi.Xingyuncolor,
		LuckyNumber:   fmt.Sprintf("%d", apiResp.Data.Xingzuoyunshi.Xingyunnumber),
		Constellation: apiResp.Data.Xingzuoyunshi.Supeixingzuo,
		Timestamp:     time.Now(),
	}

	return lunarInfo, nil
}

// generateMockLunarInfo generates mock lunar calendar data
func (s *Service) generateMockLunarInfo(date time.Time) *LunarInfo {
	// Simple mock data generation based on date
	day := date.Day()
	month := date.Month()
	year := date.Year()

	// Generate lunar date (simplified)
	lunarMonths := []string{"正月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "冬月", "腊月"}
	lunarDays := []string{"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十",
		"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
		"廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十"}

	lunarMonth := lunarMonths[(int(month)-1+day)%12]
	lunarDay := lunarDays[(day-1)%30]
	lunarDate := fmt.Sprintf("%s%s", lunarMonth, lunarDay)

	// Generate zodiac (simplified)
	zodiacs := []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	zodiac := zodiacs[year%12]

	// Generate lunar year
	lunarYear := fmt.Sprintf("%s年", zodiac)

	// Generate solar terms
	solarTerms := []string{"立春", "雨水", "惊蛰", "春分", "清明", "谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑",
		"立秋", "处暑", "白露", "秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至", "小寒", "大寒"}
	solarTerm := solarTerms[(int(month)*2-2+day/15)%24]

	// Generate suitable and avoid activities
	suitActivities := [][]string{
		{"祈福", "出行", "搬家"},
		{"开业", "签约", "投资"},
		{"结婚", "订婚", "相亲"},
		{"装修", "动土", "破土"},
		{"祭祀", "扫墓", "上香"},
	}
	avoidActivities := [][]string{
		{"诉讼", "争执", "吵架"},
		{"借贷", "放债", "赌博"},
		{"手术", "针灸", "拔牙"},
		{"远行", "出海", "登高"},
		{"搬家", "入宅", "安床"},
	}

	suitIndex := day % len(suitActivities)
	avoidIndex := (day + 1) % len(avoidActivities)

	// Generate lucky elements
	colors := []string{"红色", "金色", "绿色", "蓝色", "紫色", "白色", "黄色"}
	numbers := []string{"1", "3", "6", "8", "9"}
	constellations := []string{"白羊座", "金牛座", "双子座", "巨蟹座", "狮子座", "处女座", "天秤座", "天蝎座", "射手座", "摩羯座", "水瓶座", "双鱼座"}

	return &LunarInfo{
		Date:          date,
		LunarDate:     lunarDate,
		LunarYear:     lunarYear,
		Zodiac:        zodiac,
		SolarTerm:     solarTerm,
		Festival:      "", // Would be populated based on special dates
		Suit:          suitActivities[suitIndex],
		Avoid:         avoidActivities[avoidIndex],
		LuckyColor:    colors[day%len(colors)],
		LuckyNumber:   numbers[day%len(numbers)],
		Constellation: constellations[(int(month)-1)%12],
		Timestamp:     time.Now(),
	}
}

// GetFormattedLunar returns lunar information in a human-readable format
func (l *LunarInfo) GetFormattedLunar() string {
	result := fmt.Sprintf("📅 农历：%s %s（%s%s），节气：%s", l.LunarYear, l.LunarDate, l.Zodiac, "年", l.SolarTerm)
	
	if l.Festival != "" {
		result += fmt.Sprintf("，节日：%s", l.Festival)
	}
	
	if len(l.Suit) > 0 {
		result += fmt.Sprintf("\n✅ 宜：%s", joinStrings(l.Suit, "、"))
	}
	
	if len(l.Avoid) > 0 {
		result += fmt.Sprintf("\n❌ 忌：%s", joinStrings(l.Avoid, "、"))
	}
	
	result += fmt.Sprintf("\n🍀 幸运色：%s，幸运数字：%s", l.LuckyColor, l.LuckyNumber)
	
	return result
}

// joinStrings joins string slice with separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}