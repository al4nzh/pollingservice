package sniper

import (
	"context"
       "sort"
	"sync"
	"time"

	"github.com/al4nzh/pollingservice.git/internal/repository"
)

type SniperService struct {
    repo    *repository.MarketRepository
    fetcher Fetcher
    mu      sync.Mutex
    seen    map[string]struct{} // deduplication: listing ID
}

func NewSniperService(repo *repository.MarketRepository, fetcher Fetcher) *SniperService {
       return &SniperService{
	       repo:    repo,
	       fetcher: fetcher,
	       seen:    make(map[string]struct{}),
       }
}

// RunOnce fetches, deduplicates, evaluates, and saves good deals
func (s *SniperService) RunOnce(ctx context.Context) ([]ListingEvaluation, error) {
       listings, err := s.fetcher.FetchListings(ctx)
       if err != nil {
              return nil, err
       }

       s.mu.Lock()
       newListings := make([]Listing, 0, len(listings))
       for _, listing := range listings {
              if _, ok := s.seen[listing.ID]; ok {
                     continue
              }
              s.seen[listing.ID] = struct{}{}
              newListings = append(newListings, listing)
       }
       s.mu.Unlock()

       groupedListings := make(map[string][]Listing)
       for _, listing := range newListings {
              key := listing.Item.MarketHashName
              if key == "" {
                     key = listing.Item.ItemName
              }
              groupedListings[key] = append(groupedListings[key], listing)
       }

       goodDeals := make([]ListingEvaluation, 0)
       for _, listingsForItem := range groupedListings {
              sort.Slice(listingsForItem, func(i, j int) bool {
                     return listingsForItem[i].Price < listingsForItem[j].Price
              })

              for _, evaluation := range EvaluateListings(listingsForItem) {
                     if !evaluation.IsDeal {
                            continue
                     }

                     goodDeals = append(goodDeals, evaluation)
                     _ = s.repo.SaveSnipedDeal(ctx, evaluation)
              }
       }

       return goodDeals, nil
}

// GetRecentDeals returns recent good deals for UI/alerts
func (s *SniperService) GetRecentDeals(ctx context.Context, since time.Time) ([]repository.SnipedDeal, error) {
	return s.repo.GetSnipedDeals(ctx, since)
}
