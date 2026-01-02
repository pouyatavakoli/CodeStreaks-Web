package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/pouyatavakoli/CodeStreaks-web/internal/service"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron        *cron.Cron
	syncService service.SyncService
	interval    int // seconds
}

func NewScheduler(syncService service.SyncService, interval int) *Scheduler {
	return &Scheduler{
		cron:        cron.New(cron.WithSeconds()),
		syncService: syncService,
		interval:    interval,
	}
}

func (s *Scheduler) Start() error {
	// Convert seconds to cron expression
	cronExpr := fmt.Sprintf("@every %ds", s.interval)

	_, err := s.cron.AddFunc(cronExpr, func() {
		log.Println("Starting scheduled sync...")
		startTime := time.Now()

		if err := s.syncService.SyncAllUsers(); err != nil {
			log.Printf("Scheduled sync failed: %v", err)
		} else {
			duration := time.Since(startTime)
			log.Printf("Scheduled sync completed successfully in %v", duration)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to schedule sync job: %w", err)
	}

	s.cron.Start()
	log.Printf("Scheduler started with interval: %d seconds", s.interval)

	// Run initial sync
	log.Println("Running initial sync...")
	go func() {
		if err := s.syncService.SyncAllUsers(); err != nil {
			log.Printf("Initial sync failed: %v", err)
		}
	}()

	return nil
}

func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Scheduler stopped")
}
