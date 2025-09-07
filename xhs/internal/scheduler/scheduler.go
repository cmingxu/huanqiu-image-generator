package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"xiaohongshu-unified/internal/orchestrator"
)

// Scheduler handles daily content generation and publishing
type Scheduler struct {
	orch *orchestrator.Orchestrator
	ctx  context.Context
	cancel context.CancelFunc
}

// New creates a new scheduler instance
func New(orch *orchestrator.Orchestrator) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		orch:   orch,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins the daily scheduling loop
func (s *Scheduler) Start() error {
	log.Println("🕐 Starting daily scheduler for 8pm Beijing time...")
	
	// Beijing timezone
	beijingTZ, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return fmt.Errorf("failed to load Beijing timezone: %w", err)
	}

	// Calculate next 8pm Beijing time
	nextRun := s.getNext8PMBeijing(beijingTZ)
	log.Printf("📅 Next scheduled run: %s", nextRun.Format("2006-01-02 15:04:05 MST"))

	for {
		select {
		case <-s.ctx.Done():
			log.Println("🛑 Scheduler stopped")
			return nil
		default:
			now := time.Now().In(beijingTZ)
			
			// Check if it's time to run (within 1 minute window)
			if s.isTimeToRun(now, nextRun) {
				log.Printf("⏰ Executing scheduled run at %s", now.Format("2006-01-02 15:04:05 MST"))
				
				// Run the workflow
				if err := s.orch.Run(); err != nil {
					log.Printf("❌ Scheduled workflow failed: %v", err)
				} else {
					log.Println("✅ Scheduled workflow completed successfully")
				}
				
				// Calculate next run (tomorrow 8pm)
				nextRun = s.getNext8PMBeijing(beijingTZ)
				log.Printf("📅 Next scheduled run: %s", nextRun.Format("2006-01-02 15:04:05 MST"))
			}
			
			// Sleep for 30 seconds before checking again
			time.Sleep(30 * time.Second)
		}
	}
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	log.Println("🛑 Stopping scheduler...")
	s.cancel()
}

// getNext8PMBeijing calculates the next 8pm Beijing time
func (s *Scheduler) getNext8PMBeijing(beijingTZ *time.Location) time.Time {
	now := time.Now().In(beijingTZ)
	
	// Create 8pm today in Beijing timezone
	target := time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, beijingTZ)
	
	// If 8pm today has already passed, schedule for tomorrow
	if now.After(target) {
		target = target.Add(24 * time.Hour)
	}
	
	return target
}

// isTimeToRun checks if current time is within the execution window
func (s *Scheduler) isTimeToRun(now, target time.Time) bool {
	// Allow execution within 1 minute window (20:00:00 - 20:00:59)
	return now.After(target) && now.Before(target.Add(1*time.Minute))
}

// RunOnce executes the workflow immediately (for testing)
func (s *Scheduler) RunOnce() error {
	log.Println("🚀 Running workflow immediately...")
	return s.orch.Run()
}

// GetNextRunTime returns the next scheduled run time
func (s *Scheduler) GetNextRunTime() (time.Time, error) {
	beijingTZ, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load Beijing timezone: %w", err)
	}
	return s.getNext8PMBeijing(beijingTZ), nil
}