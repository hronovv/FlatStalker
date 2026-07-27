package repository

import (
	"context"
	"fmt"

	"flat-stalker/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Listings struct {
	pool *pgxpool.Pool
}

// Watch is a saved search URL owned by a Telegram user.
type Watch struct {
	ID     int64
	UserID int64
	ChatID int64
	URL    string
}

func NewListings(pool *pgxpool.Pool) *Listings {
	return &Listings{pool: pool}
}

func (r *Listings) Add(ctx context.Context, userID int64, url string) (*models.Listing, bool, error) {
	const q = `
INSERT INTO listings (user_id, url)
VALUES ($1, $2)
ON CONFLICT (user_id, url) DO NOTHING
RETURNING id, user_id, url;
`
	listing := &models.Listing{}
	err := r.pool.QueryRow(ctx, q, userID, url).Scan(&listing.ID, &listing.UserID, &listing.URL)
	if err == nil {
		return listing, true, nil
	}
	if err == pgx.ErrNoRows {
		return nil, false, nil
	}
	return nil, false, fmt.Errorf("add listing: %w", err)
}

func (r *Listings) ListURLsByChatID(ctx context.Context, chatID int64) ([]string, error) {
	const q = `
SELECT l.url
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

	urls := make([]string, 0)
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list listings rows: %w", err)
	}
	return urls, nil
}

func (r *Listings) ListWatches(ctx context.Context) ([]Watch, error) {
	const q = `
SELECT l.id, l.user_id, u.chat_id, l.url
FROM listings l
JOIN users u ON u.id = l.user_id
ORDER BY l.id;
`
	rows, err := r.pool.Query(ctx, q)
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
