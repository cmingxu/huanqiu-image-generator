package main

import (
	"flag"

	"github.com/sirupsen/logrus"
)

func main() {
	var (
		headless bool
		port     string
	)
	flag.BoolVar(&headless, "headless", true, "Run browser in headless mode")
	flag.StringVar(&port, "port", ":18061", "Port to run MCP server on")
	flag.Parse()

	// Initialize browser service
	browserService := NewBrowserService(headless)

	// Create and start MCP server
	mcpServer := NewMCPServer(browserService)
	if err := mcpServer.Start(port); err != nil {
		logrus.Fatalf("failed to run MCP server: %v", err)
	}
}