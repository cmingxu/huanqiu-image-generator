package main

import (
	"fmt"
	"log"
	"os"

	"xiaohongshu-content/internal/config"
	"xiaohongshu-content/internal/orchestrator"
)

func main() {
	fmt.Println("Starting Xiaohongshu Content Generator...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create orchestrator
	orch := orchestrator.New(cfg)

	// Run the content generation and posting workflow
	if err := orch.Run(); err != nil {
		log.Fatalf("Failed to run orchestrator: %v", err)
	}

	fmt.Println("Content generation and posting completed successfully!")
	os.Exit(0)
}