package service

import (
	"context"
	"fmt"

	"github.com/al4nzh/pollingservice.git/internal/fetcher"
	"github.com/al4nzh/pollingservice.git/internal/models"
	"github.com/al4nzh/pollingservice.git/internal/repository"
)

type PollRunResult struct {
	TrackedItems int `json:"tracked_items"`
	Success      int `json:"success"`
	Failed       int `json:"failed"`
}

type PollingService struct {
	repo    *repository.MarketRepository
	fetcher fetcher.Fetcher
}

func NewPollingService(repo *repository.MarketRepository, fetcher fetcher.Fetcher) *PollingService {
	return &PollingService{
		repo:    repo,
		fetcher: fetcher,
	}
}

func (s *PollingService) RunOnce(ctx context.Context) (PollRunResult, error) {
	items, err := s.repo.GetActiveTrackedItems(ctx)
	if err != nil {
		return PollRunResult{}, fmt.Errorf("get active tracked items: %w", err)
	}

	result := PollRunResult{
		TrackedItems: len(items),
	}

	for _, item := range items {
		stats, err := s.fetcher.Fetch(ctx, item)
		if err != nil {
			msg := err.Error()
			_ = s.repo.UpdateTrackedItemPollState(ctx, item.ID, false, &msg)
			result.Failed++
			continue
		}

		if err := s.repo.UpsertMarketStats(ctx, stats); err != nil {
			msg := err.Error()
			_ = s.repo.UpdateTrackedItemPollState(ctx, item.ID, false, &msg)
			result.Failed++
			continue
		}

		history := models.PriceHistory{
			AppID:                  stats.AppID,
			MarketHashName:         stats.MarketHashName,
			DisplayName:            stats.DisplayName,
			Source:                 stats.Source,
			CurrentPrice:           stats.CurrentPrice,
			CSFloatRefPrice:        stats.CSFloatRefPrice,
			CSFloatSecondBestPrice: stats.CSFloatSecondBestPrice,
			CSFloatSecondRefPrice:  stats.CSFloatSecondRefPrice,
			Currency:               stats.Currency,
		}

		if err := s.repo.InsertPriceHistory(ctx, history); err != nil {
			msg := err.Error()
			_ = s.repo.UpdateTrackedItemPollState(ctx, item.ID, false, &msg)
			result.Failed++
			continue
		}

		if err := s.repo.UpdateTrackedItemPollState(ctx, item.ID, true, nil); err != nil {
			result.Failed++
			continue
		}

		result.Success++
	}

	return result, nil
}