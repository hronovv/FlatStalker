package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SeenAds struct {
	pool *pgxpool.Pool
}

func NewSeenAds(pool *pgxpool.Pool) *SeenAds {
	return &SeenAds{pool: pool}
}

func (r *SeenAds) CountByListing(ctx context.Context, listingID int64) (int, error) {
	const q = `SELECT COUNT(*) FROM seen_ads WHERE listing_id = $1`
	var n int
	if err := r.pool.QueryRow(ctx, q, listingID).Scan(&n); err != nil {
		return 0, fmt.Errorf("count seen ads: %w", err)
	}
	return n, nil
}

func (r *SeenAds) ExistingIDs(ctx context.Context, listingID int64, adIDs []int64) (map[int64]struct{}, error) {
	out := make(map[int64]struct{})
	if len(adIDs) == 0 {
		return out, nil
	}

	const q = `
SELECT ad_id
FROM seen_ads
WHERE listing_id = $1 AND ad_id = ANY($2);
`
	rows, err := r.pool.Query(ctx, q, listingID, adIDs)
	if err != nil {
		return nil, fmt.Errorf("existing seen ads: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan seen ad: %w", err)
		}
		out[id] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("seen ads rows: %w", err)
	}
	return out, nil
}

func (r *SeenAds) MarkMany(ctx context.Context, listingID int64, adIDs []int64) error {
	if len(adIDs) == 0 {
		return nil
	}

	const q = `
INSERT INTO seen_ads (listing_id, ad_id)
SELECT $1, unnest($2::bigint[])
ON CONFLICT (listing_id, ad_id) DO NOTHING;
`
	if _, err := r.pool.Exec(ctx, q, listingID, adIDs); err != nil {
		return fmt.Errorf("mark seen ads: %w", err)
	}
	return nil
}
