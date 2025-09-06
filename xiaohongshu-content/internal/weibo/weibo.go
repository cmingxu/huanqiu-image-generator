package weibo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// WeiboPost represents a single weibo post
type WeiboPost struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Created  string `json:"created_at"`
	UserName string `json:"user_name"`
}

// WeiboResponse represents the API response structure
type WeiboResponse struct {
	Ok   int `json:"ok"`
	Data struct {
		List []struct {
			Mblogid   string `json:"mblogid"`
			CreatedAt string `json:"created_at"`
			TextRaw   string `json:"text_raw"`
			Text      string `json:"text"`
			User      struct {
				ScreenName string `json:"screen_name"`
			} `json:"user"`
		} `json:"list"`
	} `json:"data"`
}

// Service handles weibo content fetching
type Service struct {
	UID     string
	Cookies string
	Token   string
	client  *http.Client
}

// NewService creates a new weibo service
func NewService(uid, cookies, token string) *Service {
	return &Service{
		UID:     uid,
		Cookies: cookies,
		Token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetLatestPosts fetches the latest weibo posts
func (s *Service) GetLatestPosts(page int) ([]WeiboPost, error) {
	url := fmt.Sprintf("https://weibo.com/ajax/statuses/mymblog?uid=%s&page=%d&feature=0", s.UID, page)

	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set all headers to simulate browser request
	s.setHeaders(req)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}



	var weiboResp WeiboResponse
	if err := json.Unmarshal(body, &weiboResp); err != nil {
		log.Printf("[DEBUG] JSON Parse Error: %v", err)
		log.Printf("[DEBUG] Response body that failed to parse: %s", string(body))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}



	if weiboResp.Ok != 1 {
		return nil, fmt.Errorf("weibo API returned error: %s", string(body))
	}

	var posts []WeiboPost
	for _, item := range weiboResp.Data.List {
		cleanText := s.cleanText(item.TextRaw)
		if cleanText == "" {
			cleanText = s.cleanText(item.Text)
		}
		
		posts = append(posts, WeiboPost{
			ID:       item.Mblogid,
			Text:     cleanText,
			Created:  item.CreatedAt,
			UserName: item.User.ScreenName,
		})
	}

	return posts, nil
}

// setHeaders sets all the headers to simulate browser request
func (s *Service) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,zh-TW;q=0.6")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Client-Version", "v2.47.106")
	req.Header.Set("Cookie", s.Cookies)
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", fmt.Sprintf("https://weibo.com/u/%s?is_all=1", s.UID))
	req.Header.Set("Sec-Ch-Ua", `"Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Server-Version", "v2025.09.05.1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-Xsrf-Token", s.Token)
}

// cleanText removes HTML tags and cleans up the text content
func (s *Service) cleanText(text string) string {
	if text == "" {
		return ""
	}

	// Remove HTML tags
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	cleanedText := htmlTagRegex.ReplaceAllString(text, "")

	// Remove extra whitespace and newlines
	cleanedText = regexp.MustCompile(`\s+`).ReplaceAllString(cleanedText, " ")
	cleanedText = strings.TrimSpace(cleanedText)

	// Remove common weibo artifacts
	cleanedText = regexp.MustCompile(`全文$`).ReplaceAllString(cleanedText, "")
	cleanedText = regexp.MustCompile(`展开$`).ReplaceAllString(cleanedText, "")
	cleanedText = strings.TrimSpace(cleanedText)

	return cleanedText
}

// GetRecentContent gets recent weibo content for content generation
func (s *Service) GetRecentContent() (string, error) {
	posts, err := s.GetLatestPosts(1)
	if err != nil {
		return "", err
	}
	if len(posts) == 0 {
		return "", fmt.Errorf("no posts found")
	}

	// Get the most recent post
	latestPost := posts[0]
	return fmt.Sprintf("最新微博内容：%s (发布时间：%s)", latestPost.Text, latestPost.Created), nil
}

// GetTop4PostsForSummary gets top 4 recent posts for LLM summarization
func (s *Service) GetTop2PostsForSummary() (string, error) {
	posts, err := s.GetLatestPosts(1)
	if err != nil {
		return "", err
	}
	if len(posts) == 0 {
		return "", fmt.Errorf("no posts found")
	}

	// Get top 4 posts (or all if less than 4)
	maxPosts := 2
	if len(posts) < maxPosts {
		maxPosts = len(posts)
	}

	var contentBuilder strings.Builder
	contentBuilder.WriteString("以下是北京环球度假区官方微博最新动态，请总结其中的新闻和活动信息：\n\n")

	for i := 0; i < maxPosts; i++ {
		post := posts[i]
		// Clean HTML tags from text
		cleanText := s.cleanHTMLTags(post.Text)
		contentBuilder.WriteString(fmt.Sprintf("%d. 发布时间：%s\n内容：%s\n\n", i+1, post.Created, cleanText))
	}

	return contentBuilder.String(), nil
}

// cleanHTMLTags removes HTML tags from text
func (s *Service) cleanHTMLTags(text string) string {
	// Remove <br /> tags and replace with spaces
	text = strings.ReplaceAll(text, "<br />", " ")
	text = strings.ReplaceAll(text, "<br/>", " ")
	text = strings.ReplaceAll(text, "<br>", " ")
	
	// Remove other HTML tags using regex
	re := regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, "")
	
	// Clean up multiple spaces
	spaceRe := regexp.MustCompile(`\s+`)
	text = spaceRe.ReplaceAllString(text, " ")
	
	// Trim whitespace
	return strings.TrimSpace(text)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}