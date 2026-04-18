package models

import "time"

type MarketStats struct {
	ID                int64
	AppID             int
	MarketHashName    string
	DisplayName       string
	Source            string

	CurrentPrice      float64   // actual CSFloat listing price
	CSFloatRefPrice   *float64  // item.scm.price
	CSFloatSecondBestPrice    *float64  // optional later
	CSFloatSecondRefPrice   *float64  // optional later
	//SteamRefPrice     *float64  // optional later

	Currency          string
	ListingCount      *int
	UpdatedAt         time.Time
	RawPayload        []byte
}

type TrackedItem struct {
	ID               int64
	AppID            int
	MarketHashName   string
	DisplayName      string
	Source           string
	IsActive         bool
	Priority         int
	LastPolledAt     *time.Time
	LastSuccessAt    *time.Time
	LastErrorAt      *time.Time
	LastErrorMessage *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
type PriceHistory struct {
	ID              int64
	AppID           int
	MarketHashName  string
	DisplayName     string
	Source          string
	CurrentPrice    float64
	CSFloatRefPrice *float64
	CSFloatSecondBestPrice *float64
	CSFloatSecondRefPrice *float64
	Currency        string
	CapturedAt      time.Time
}