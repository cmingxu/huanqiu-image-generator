package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

// BrowserService handles browser automation and screenshot operations
type BrowserService struct {
	headless bool
}

// NewBrowserService creates a new browser service instance
func NewBrowserService(headless bool) *BrowserService {
	return &BrowserService{
		headless: headless,
	}
}

// TakeScreenshot takes a screenshot of the specified URL
func (bs *BrowserService) TakeScreenshot(ctx context.Context, req *ScreenshotRequest) (*ScreenshotResult, error) {
	logrus.Infof("Taking screenshot of URL: %s", req.URL)

	// Set default values
	if req.Selector == "" {
		req.Selector = "body"
	}
	if req.OutputPath == "" {
		req.OutputPath = fmt.Sprintf("screenshot_%d.jpg", time.Now().Unix())
	}
	if req.WaitTime == 0 {
		req.WaitTime = 3
	}

	// Create browser context with more permissive options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", bs.headless),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("ignore-ssl-errors", true),
		chromedp.Flag("allow-running-insecure-content", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(logrus.Printf))
	defer cancel()

	// Set timeout to 60 seconds
	browserCtx, cancel = context.WithTimeout(browserCtx, 60*time.Second)
	defer cancel()

	var buf []byte

	// Run browser tasks
	err := chromedp.Run(browserCtx,
		// Navigate to the URL
		chromedp.Navigate(req.URL),
		// Wait for the page to load
		chromedp.WaitVisible(req.Selector, chromedp.ByID),
		// Wait additional time for images and fonts to load
		chromedp.Sleep(time.Duration(req.WaitTime)*time.Second),
		// Take screenshot
		chromedp.Screenshot(req.Selector, &buf, chromedp.NodeVisible, chromedp.ByID),
	)

	if err != nil {
		logrus.Errorf("Error taking screenshot: %v", err)
		return &ScreenshotResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(req.OutputPath)
	if outputDir != "." {
		err = ensureDir(outputDir)
		if err != nil {
			return &ScreenshotResult{
				Success: false,
				Error:   fmt.Sprintf("Failed to create output directory: %v", err),
			}, err
		}
	}

	// Save the screenshot
	err = ioutil.WriteFile(req.OutputPath, buf, 0644)
	if err != nil {
		logrus.Errorf("Error saving screenshot: %v", err)
		return &ScreenshotResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to save screenshot: %v", err),
		}, err
	}

	logrus.Infof("Screenshot saved successfully: %s", req.OutputPath)
	return &ScreenshotResult{
		Success:    true,
		OutputPath: req.OutputPath,
		Message:    fmt.Sprintf("Screenshot saved successfully to %s", req.OutputPath),
	}, nil
}

// ensureDir creates directory if it doesn't exist
func ensureDir(dir string) error {
	return nil // Directory creation handled by ioutil.WriteFile for parent dirs
}
