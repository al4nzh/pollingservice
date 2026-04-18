package service

import (
	"context"
	"log"
	"time"
)

func (s *PollingService) StartScheduler(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("scheduler started, interval=%s\n", interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler stopped")
			return

		case <-ticker.C:
			log.Println("starting scheduled poll run")

			result, err := s.RunOnce(ctx)
			if err != nil {
				log.Printf("scheduled poll failed: %v\n", err)
				continue
			}

			log.Printf(
				"scheduled poll finished: tracked=%d success=%d failed=%d\n",
				result.TrackedItems,
				result.Success,
				result.Failed,
			)
		}
	}
}