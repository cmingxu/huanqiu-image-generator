package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
)

const defaultText = `
8 月 3 日入园人数: <span style="color: #ff0000; font-weight: bold;">19999</span><br/>天气晴朗适合游玩
`

func main() {
	// Define command line flags
	addr := flag.String("addr", "http://localhost:3000", "Address of the Next.js project")
	text := flag.String("text", defaultText, "Text to display on the image")
	output := flag.String("output", "untitled.jpg", "Output file path and name (.jpg)")
	image := flag.String("image", "/assets/6.jpg", "Background image path")

	flag.Parse()

	// Build URL with query parameters
	baseURL := *addr
	params := url.Values{}
	params.Add("image", *image)
	params.Add("text", *text)
	params.Add("autoExport", "false")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	println(fullURL)

	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte

	// Run tasks
	err := chromedp.Run(ctx,
		// Navigate to the URL
		chromedp.Navigate(fullURL),
		// Wait for the page to load
		chromedp.WaitVisible(`#exportable`, chromedp.ByID),
		// Wait additional time for images and fonts to load
		chromedp.Sleep(3*time.Second),
		// Take screenshot of the specific element
		chromedp.Screenshot(`#exportable`, &buf, chromedp.NodeVisible, chromedp.ByID),
	)

	if err != nil {
		log.Fatal("Error running chromedp tasks:", err)
	}

	// Save the screenshot
	err = ioutil.WriteFile(*output, buf, 0644)
	if err != nil {
		log.Fatal("Error saving screenshot:", err)
	}

	fmt.Printf("Screenshot saved as %s\n", *output)
}

