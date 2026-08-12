package repository

import (
	"context"
	"errors"
	"fmt"

	"flat-stalker/internal/cache"
	"flat-stalker/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Listings struct {
	pool   *pgxpool.Pool
	byChat *cache.LRU[[]models.Listing]
}

// Watch is a saved search URL owned by a Telegram user.
type Watch struct {
	ID     int64
	UserID int64
	ChatID int64
	URL    string
}

func NewListings(pool *pgxpool.Pool) *Listings {
	return &Listings{
		pool:   pool,
		byChat: cache.NewLRU[[]models.Listing](cache.DefaultSize, cache.DefaultTTL),
	}
}

func (r *Listings) Add(ctx context.Context, userID, chatID int64, url string) (*models.Listing, bool, error) {
	const q = `
INSERT INTO listings (user_id, url)
VALUES ($1, $2)
ON CONFLICT (user_id, url) DO NOTHING
RETURNING id, user_id, url, paused;
`
	listing := &models.Listing{}
	err := r.pool.QueryRow(ctx, q, userID, url).Scan(
		&listing.ID, &listing.UserID, &listing.URL, &listing.Paused,
	)
	if err == nil {
		r.byChat.Delete(chatID)
		return listing, true, nil
	}
	if err == pgx.ErrNoRows {
		return nil, false, nil
	}
	return nil, false, fmt.Errorf("add listing: %w", err)
}

func (r *Listings) ListByChatID(ctx context.Context, chatID int64) ([]models.Listing, error) {
	if cached, ok := r.byChat.Get(chatID); ok {
		return cloneListings(cached), nil
	}

	const q = `
SELECT l.id, l.user_id, l.url, l.paused
FROM listings l
JOIN users u ON u.id = l.user_id
WHERE u.chat_id = $1
ORDER BY l.id;
`
	rows, err := r.pool.Query(ctx, q, chatID)
	if err != nil {
		return nil, fmt.Errorf("list listings: %w", err)
	}
	defer rows.Close()

	out := make([]models.Listing, 0)
	for rows.Next() {
		var l models.Listing
		if err := rows.Scan(&l.ID, &l.UserID, &l.URL, &l.Paused); err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list listings rows: %w", err)
	}

	r.byChat.Put(chatID, cloneListings(out))
	return out, nil
}

func (r *Listings) CountByUserID(ctx context.Context, userID int64) (int, error) {
	const q = `SELECT COUNT(*) FROM listings WHERE user_id = $1`
	var n int
	if err := r.pool.QueryRow(ctx, q, userID).Scan(&n); err != nil {
		return 0, fmt.Errorf("count listings: %w", err)
	}
	return n, nil
}

func (r *Listings) SetPaused(ctx context.Context, listingID, chatID int64, paused bool) (*models.Listing, error) {
	const q = `
UPDATE listings l
SET paused = $3
FROM users u
WHERE l.id = $1
  AND l.user_id = u.id
  AND u.chat_id = $2
RETURNING l.id, l.user_id, l.url, l.paused;
`
	listing := &models.Listing{}
	err := r.pool.QueryRow(ctx, q, listingID, chatID, paused).Scan(
		&listing.ID, &listing.UserID, &listing.URL, &listing.Paused,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("set paused: %w", err)
	}
	r.byChat.Delete(chatID)
	return listing, nil
}

func (r *Listings) Delete(ctx context.Context, listingID, chatID int64) (bool, error) {
	const q = `
DELETE FROM listings l
USING users u
WHERE l.id = $1
  AND l.user_id = u.id
  AND u.chat_id = $2;
`
	tag, err := r.pool.Exec(ctx, q, listingID, chatID)
	if err != nil {
		return false, fmt.Errorf("delete listing: %w", err)
	}
	if tag.RowsAffected() > 0 {
		r.byChat.Delete(chatID)
	}
	return tag.RowsAffected() > 0, nil
}

func (r *Listings) ListWatches(ctx context.Context, userPlan string) ([]Watch, error) {
	const q = `
SELECT l.id, l.user_id, u.chat_id, l.url
FROM listings l
JOIN users u ON u.id = l.user_id
WHERE l.paused = false
  AND u.plan = $1
ORDER BY l.id;
`
	rows, err := r.pool.Query(ctx, q, userPlan)
	if err != nil {
		return nil, fmt.Errorf("list watches: %w", err)
	}
	defer rows.Close()

	watches := make([]Watch, 0)
	for rows.Next() {
		var w Watch
		if err := rows.Scan(&w.ID, &w.UserID, &w.ChatID, &w.URL); err != nil {
			return nil, fmt.Errorf("scan watch: %w", err)
		}
		watches = append(watches, w)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("watches rows: %w", err)
	}
	return watches, nil
}

func cloneListings(in []models.Listing) []models.Listing {
	out := make([]models.Listing, len(in))
	copy(out, in)
	return out
}
