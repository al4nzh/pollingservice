package sniper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Fetcher interface {
    FetchListings(ctx context.Context) ([]Listing, error)
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


	// ListingsResponse represents the response from CSFloat API
type ListingsResponse struct {
	Data []Listing `json:"data"`
}

	// Listing represents a single listing from CSFloat
type Listing struct {
	ID    string `json:"id"`
	Price int64  `json:"price"`
	Item struct {
		AppID          int    `json:"appid"`
		MarketHashName string `json:"market_hash_name"`
		ItemName       string `json:"item_name"`
		WearName       string `json:"wear_name"`
	} `json:"item"`
	Reference struct {
		Price    int64 `json:"predicted_price"` // cents
		Quantity int   `json:"quantity"`
	} `json:"reference"`
}

func (f *CSFloatFetcher) FetchListings(ctx context.Context) ([]Listing, error) {
	endpoint, err := f.buildListingsURL()
	if err != nil {
		return nil, fmt.Errorf("build csfloat url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create csfloat request: %w", err)
	}

	if f.apiKey != "" {
		req.Header.Set("Authorization", f.apiKey)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("csfloat api error: %d - %s", resp.StatusCode, string(body))
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var listings ListingsResponse
	if err := json.Unmarshal(rawBody, &listings); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(listings.Data) == 0 {
		return nil, fmt.Errorf("no recent listings found")
	}

	return listings.Data, nil
}

func (f *CSFloatFetcher) buildListingsURL() (string, error) {
	u, err := url.Parse(f.baseURL + "/api/v1/listings")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("limit", "50")
	q.Set("sort_by", "most_recent")
	q.Set("type", "buy_now")
	u.RawQuery = q.Encode()

	return u.String(), nil
}