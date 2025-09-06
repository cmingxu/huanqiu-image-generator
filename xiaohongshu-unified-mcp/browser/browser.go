package browser

import (
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/headless_browser"
	"xiaohongshu-unified-mcp/cookies"
)

func NewBrowser(headless bool) *headless_browser.Browser {

	opts := []headless_browser.Option{
		headless_browser.WithHeadless(headless),
	}

	// 加载 cookies
	cookiePath := cookies.GetCookiesFilePath()
	cookieLoader := cookies.NewLoadCookie(cookiePath)

	if data, err := cookieLoader.LoadCookies(); err == nil {
		opts = append(opts, headless_browser.WithCookies(string(data)))
		logrus.Debugf("loaded cookies from file successfully")
	} else {
		logrus.Warnf("failed to load cookies: %v", err)
	}

	return headless_browser.New(opts...)
}