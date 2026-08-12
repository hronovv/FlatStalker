package repository

import (
	"context"
	"fmt"

	"flat-stalker/internal/cache"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Bans struct {
	pool   *pgxpool.Pool
	byChat *cache.LRU[bool]
}

func NewBans(pool *pgxpool.Pool) *Bans {
	return &Bans{
		pool:   pool,
		byChat: cache.NewLRU[bool](cache.DefaultSize, cache.DefaultTTL),
	}
}

func (r *Bans) IsBanned(ctx context.Context, chatID int64) (bool, error) {
	if banned, ok := r.byChat.Get(chatID); ok {
		return banned, nil
	}

	const q = `SELECT EXISTS(SELECT 1 FROM banned_users WHERE chat_id = $1)`
	var banned bool
	if err := r.pool.QueryRow(ctx, q, chatID).Scan(&banned); err != nil {
		return false, fmt.Errorf("check ban: %w", err)
	}
	if banned {
		r.byChat.Put(chatID, true)
	}
	return banned, nil
}
