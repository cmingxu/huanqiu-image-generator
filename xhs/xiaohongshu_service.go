package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"xiaohongshu-unified/browser"
	"xiaohongshu-unified/configs"
	"xiaohongshu-unified/xiaohongshu"
)

// XiaohongshuService 小红书业务服务
type XiaohongshuService struct {
	headless bool
}

// NewXiaohongshuService 创建小红书服务实例
func NewXiaohongshuService(headless bool) *XiaohongshuService {
	configs.InitHeadless(headless)
	return &XiaohongshuService{
		headless: headless,
	}
}

// CheckLoginStatus 检查登录状态
func (s *XiaohongshuService) CheckLoginStatus(ctx context.Context) (*LoginStatusResponse, error) {
	logrus.Info("Checking Xiaohongshu login status")
	
	b := browser.NewBrowser(configs.IsHeadless())
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	loginAction := xiaohongshu.NewLogin(page)

	isLoggedIn, err := loginAction.CheckLoginStatus(ctx)
	if err != nil {
		return nil, err
	}

	response := &LoginStatusResponse{
		IsLoggedIn: isLoggedIn,
		Username:   configs.Username,
	}

	return response, nil
}

// PublishContent 发布内容
func (s *XiaohongshuService) PublishContent(ctx context.Context, req *PublishRequest) (*PublishResponse, error) {
	logrus.Infof("Publishing content: %s", req.Title)
	
	b := browser.NewBrowser(configs.IsHeadless())
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	// 构建发布内容
	content := xiaohongshu.PublishImageContent{
		Title:      req.Title,
		Content:    req.Content,
		ImagePaths: req.Images,
	}

	// 执行发布
	publishAction, err := xiaohongshu.NewPublishImageAction(page)
	if err != nil {
		return nil, err
	}
	err = publishAction.Publish(ctx, content)
	if err != nil {
		return nil, err
	}

	response := &PublishResponse{
		Title:   req.Title,
		Content: req.Content,
		Images:  len(req.Images),
		Status:  "发布完成",
	}

	return response, nil
}

// ListFeeds 获取Feeds列表
func (s *XiaohongshuService) ListFeeds(ctx context.Context) (*FeedsListResponse, error) {
	logrus.Info("Listing Xiaohongshu feeds")
	
	b := browser.NewBrowser(configs.IsHeadless())
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	feedsAction := xiaohongshu.NewFeedsListAction(page)
	feeds, err := feedsAction.GetFeedsList(ctx)
	if err != nil {
		return nil, err
	}

	// Convert xiaohongshu.Feed to main.Feed
	mainFeeds := make([]Feed, len(feeds))
	for i, feed := range feeds {
		mainFeeds[i] = Feed{
			ID:       feed.ID,
			Title:    feed.NoteCard.DisplayTitle,
			Content:  feed.NoteCard.DisplayTitle, // Using title as content for now
			Author:   feed.NoteCard.User.Nickname,
			Likes:    0, // TODO: Parse from InteractInfo
			Comments: 0, // TODO: Parse from InteractInfo
			URL:      "https://xiaohongshu.com/explore/" + feed.ID,
		}
	}

	response := &FeedsListResponse{
		Feeds: mainFeeds,
		Count: len(mainFeeds),
	}

	return response, nil
}

// SearchFeeds 搜索Feeds
func (s *XiaohongshuService) SearchFeeds(ctx context.Context, keyword string) (*SearchResponse, error) {
	logrus.Infof("Searching feeds with keyword: %s", keyword)
	
	b := browser.NewBrowser(configs.IsHeadless())
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	searchAction := xiaohongshu.NewSearchAction(page)
	results, err := searchAction.Search(ctx, keyword)
	if err != nil {
		return nil, err
	}

	// Convert xiaohongshu.Feed to main.Feed
	mainResults := make([]Feed, len(results))
	for i, feed := range results {
		mainResults[i] = Feed{
			ID:       feed.ID,
			Title:    feed.NoteCard.DisplayTitle,
			Content:  feed.NoteCard.DisplayTitle, // Using title as content for now
			Author:   feed.NoteCard.User.Nickname,
			Likes:    0, // TODO: Parse from InteractInfo
			Comments: 0, // TODO: Parse from InteractInfo
			URL:      "https://xiaohongshu.com/explore/" + feed.ID,
		}
	}

	response := &SearchResponse{
		Keyword: keyword,
		Results: mainResults,
		Total:   len(mainResults),
	}

	return response, nil
}