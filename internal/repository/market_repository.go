package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/al4nzh/pollingservice.git/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SnipedDeal represents a good deal found by the sniper
type SnipedDeal struct {
    ID             int64
    ListingID      string
    AppID          int
    MarketHashName string
    DisplayName    string
    Price          float64
    DiscountFromRef    float64
    DiscountFromSecond float64
    Reason         string
    CapturedAt     time.Time
}

// SaveSnipedDeal saves a good deal found by the sniper
func (r *MarketRepository) SaveSnipedDeal(ctx context.Context, eval interface{}) error {
	e, ok := eval.(interface {
		ListingID() string
		AppID() int
		MarketHashName() string
		DisplayName() string
		DealPrice() float64
		DealDiscountFromRef() float64
		DealDiscountFromSecond() float64
		DealReason() string
	})
    if !ok {
        return fmt.Errorf("invalid evaluation type")
    }

    query := `
        INSERT INTO sniped_deals (
            listing_id, app_id, market_hash_name, display_name, price, discount_from_ref, discount_from_second, reason, captured_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        ON CONFLICT (listing_id) DO NOTHING
    `
    _, err := r.pool.Exec(ctx, query,
        e.ListingID(),
		e.AppID(),
		e.MarketHashName(),
		e.DisplayName(),
		e.DealPrice(),
		e.DealDiscountFromRef(),
		e.DealDiscountFromSecond(),
		e.DealReason(),
        time.Now(),
    )
    if err != nil {
        return fmt.Errorf("save sniped deal: %w", err)
    }
    return nil
}

// GetSnipedDeals returns all good deals since a given time
func (r *MarketRepository) GetSnipedDeals(ctx context.Context, since time.Time) ([]SnipedDeal, error) {
    query := `
        SELECT listing_id, app_id, market_hash_name, display_name, price, discount_from_ref, discount_from_second, reason, captured_at
        FROM sniped_deals
        WHERE captured_at >= $1
        ORDER BY captured_at DESC
    `
    rows, err := r.pool.Query(ctx, query, since)
    if err != nil {
        return nil, fmt.Errorf("get sniped deals: %w", err)
    }
    defer rows.Close()

	var deals []SnipedDeal
    for rows.Next() {
        var d SnipedDeal
        err := rows.Scan(&d.ListingID, &d.AppID, &d.MarketHashName, &d.DisplayName, &d.Price, &d.DiscountFromRef, &d.DiscountFromSecond, &d.Reason, &d.CapturedAt)
        if err != nil {
            return nil, fmt.Errorf("scan sniped deal: %w", err)
        }
        deals = append(deals, d)
    }

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sniped deals: %w", err)
	}

    return deals, nil
}

type MarketRepository struct {
	pool *pgxpool.Pool
}

func NewMarketRepository(pool *pgxpool.Pool) *MarketRepository {
	return &MarketRepository{pool: pool}
}

func (r *MarketRepository) GetActiveTrackedItems(ctx context.Context) ([]models.TrackedItem, error) {
	query := `
		SELECT id, app_id, market_hash_name, display_name, source,
		       is_active, priority,
		       last_polled_at, last_success_at, last_error_at, last_error_message,
		       created_at, updated_at
		FROM tracked_items
		WHERE is_active = true
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query tracked_items: %w", err)
	}
	defer rows.Close()

	var items []models.TrackedItem

	for rows.Next() {
		var item models.TrackedItem

		err := rows.Scan(
			&item.ID,
			&item.AppID,
			&item.MarketHashName,
			&item.DisplayName,
			&item.Source,
			&item.IsActive,
			&item.Priority,
			&item.LastPolledAt,
			&item.LastSuccessAt,
			&item.LastErrorAt,
			&item.LastErrorMessage,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan tracked_item: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *MarketRepository) UpsertMarketStats(ctx context.Context, stats models.MarketStats) error {
	query := `
		INSERT INTO items_market_stats (
			app_id,
			market_hash_name,
			display_name,
			source,
			current_price,
			csfloat_ref_price,
			csfloat_second_best_price,
			csfloat_second_ref_price,
			currency,
			listing_count,
			updated_at,
			raw_payload
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (app_id, market_hash_name)
		DO UPDATE SET
			display_name = EXCLUDED.display_name,
			source = EXCLUDED.source,
			current_price = EXCLUDED.current_price,
			csfloat_ref_price = EXCLUDED.csfloat_ref_price,
			csfloat_second_best_price = EXCLUDED.csfloat_second_best_price,
			csfloat_second_ref_price = EXCLUDED.csfloat_second_ref_price,
			currency = EXCLUDED.currency,
			listing_count = EXCLUDED.listing_count,
			updated_at = EXCLUDED.updated_at,
			raw_payload = EXCLUDED.raw_payload
		`

	_, err := r.pool.Exec(ctx, query,
		stats.AppID,
		stats.MarketHashName,
		stats.DisplayName,
		stats.Source,
		stats.CurrentPrice,
		stats.CSFloatRefPrice,
		stats.CSFloatSecondBestPrice,
		stats.CSFloatSecondRefPrice,
		stats.Currency,
		stats.ListingCount,
		stats.UpdatedAt,
		stats.RawPayload,
	)
	if err != nil {
		return fmt.Errorf("upsert market stats: %w", err)
	}

	return nil
}
func (r *MarketRepository) UpdateTrackedItemPollState(
	ctx context.Context,
	itemID int64,
	success bool,
	errMsg *string,
) error {
	now := time.Now()
	query := `
		UPDATE tracked_items
		SET last_polled_at = $2,
		    last_success_at = CASE WHEN $3 THEN $2 ELSE last_success_at END,
		    last_error_at = CASE WHEN NOT $3 THEN $2 ELSE last_error_at END,
		    last_error_message = $4
		WHERE id = $1
		`
	_, err := r.pool.Exec(ctx, query, itemID, now, success, errMsg)
	if err != nil {
		return fmt.Errorf("update tracked item poll state: %w", err)
	}
	return nil
}
func (r *MarketRepository) InsertPriceHistory(ctx context.Context, history models.PriceHistory) error {
	query := `
		INSERT INTO item_price_history (
			app_id,
			market_hash_name,
			display_name,
			source,
			current_price,
			csfloat_ref_price,
			csfloat_second_best_price,
			csfloat_second_ref_price,
			currency,
			captured_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`

	_, err := r.pool.Exec(ctx, query,
		history.AppID,
		history.MarketHashName,
		history.DisplayName,
		history.Source,
		history.CurrentPrice,
		history.CSFloatRefPrice,
		history.CSFloatSecondBestPrice,
		history.CSFloatSecondRefPrice,	
		history.Currency,
		history.CapturedAt,
	)
	if err != nil {
		return fmt.Errorf("insert price history: %w", err)
	}

	return nil
}