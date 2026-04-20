package sniper

import (
	"context"
	"log"
	"time"
)

type SniperScheduler struct {
	service *SniperService
	interval time.Duration
	stopCh   chan struct{}
}

func NewSniperScheduler(service *SniperService, interval time.Duration) *SniperScheduler {
	return &SniperScheduler{
		service: service,
		interval: interval,
		stopCh: make(chan struct{}),
	}
}

func (s *SniperScheduler) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				deals, err := s.service.RunOnce(ctx)
				if err != nil {
					log.Printf("[sniper] error: %v", err)
					continue
				}
				if len(deals) > 0 {
					log.Printf("[sniper] found %d good recent deals", len(deals))
				}
			case <-s.stopCh:
				return
			}
		}
	}()
}

func (s *SniperScheduler) Stop() {
	close(s.stopCh)
}
