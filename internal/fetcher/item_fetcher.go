package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/al4nzh/pollingservice.git/internal/models"
)

type Fetcher interface {
	Fetch(ctx context.Context, item models.TrackedItem) (models.MarketStats, error)
}

type CSFloatFetcher struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewCSFloatFetcher(client *http.Client, apiKey string) *CSFloatFetcher {
	return &CSFloatFetcher{
		client:  client,
		baseURL: "https://csfloat.com",
		apiKey:  apiKey,
	}
}

type ListingsResponse struct {
	Data []Listing `json:"data"`
}

type Listing struct {
    ID        string `json:"id"`
    Price     int64  `json:"price"`
    Item      struct {
        MarketHashName string `json:"market_hash_name"`
        ItemName       string `json:"item_name"`
        WearName       string `json:"wear_name"`
    } `json:"item"`
	Reference struct {
		Price          int64   `json:"predicted_price"` // cents
		Quantity       int     `json:"quantity"`
	} `json:"reference"`
}


func (f *CSFloatFetcher) Fetch(ctx context.Context, item models.TrackedItem) (models.MarketStats, error) {
	endpoint, err := f.buildListingsURL(item)
	if err != nil {
		return models.MarketStats{}, fmt.Errorf("build csfloat url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return models.MarketStats{}, fmt.Errorf("create csfloat request: %w", err)
	}

	if f.apiKey != "" {
		req.Header.Set("Authorization", f.apiKey)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return models.MarketStats{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return models.MarketStats{}, fmt.Errorf("csfloat api error: %d - %s", resp.StatusCode, string(body))
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.MarketStats{}, fmt.Errorf("read body: %w", err)
	}

	var listings ListingsResponse
	if err := json.Unmarshal(rawBody, &listings); err != nil {
		return models.MarketStats{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(listings.Data) == 0 {
		return models.MarketStats{}, fmt.Errorf("no listings found for item: %s", item.MarketHashName)
	}

	best := listings.Data[0]
	currPrice := float64(best.Price) / 100.0

	var refPrice *float64
	if best.Reference.Price > 0 {
		v := float64(best.Reference.Price) / 100.0
		refPrice = &v
	}

	var secondBestPrice *float64
	if len(listings.Data) > 1 {
		v := float64(listings.Data[1].Price) / 100.0
		secondBestPrice = &v
	}

	var secondRefPrice *float64
	if len(listings.Data) > 1 && listings.Data[1].Reference.Price > 0 {
		v := float64(listings.Data[1].Reference.Price) / 100.0
		secondRefPrice = &v
	}

	displayName := item.DisplayName
	if displayName == "" {
		displayName = best.Item.MarketHashName
	}

	var listingCount *int
	if best.Reference.Quantity > 0 {
		q := best.Reference.Quantity
		listingCount = &q
	}

	return models.MarketStats{
		AppID:                 item.AppID,
		MarketHashName:        item.MarketHashName,
		DisplayName:           displayName,
		Source:                "csfloat",
		CurrentPrice:          currPrice,
		CSFloatRefPrice:       refPrice,
		CSFloatSecondBestPrice: secondBestPrice,
		CSFloatSecondRefPrice: secondRefPrice,
		Currency:              "USD",
		ListingCount:          listingCount,
		UpdatedAt:             time.Now(),
		RawPayload:            rawBody,
	}, nil
}

func (f *CSFloatFetcher) buildListingsURL(item models.TrackedItem) (string, error) {
	u, err := url.Parse(f.baseURL + "/api/v1/listings")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("market_hash_name", item.MarketHashName)
	q.Set("limit", "2")
	q.Set("sort_by", "best_deal")
	q.Set("type", "buy_now")
	u.RawQuery = q.Encode()

	return u.String(), nil
}