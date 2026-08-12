package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Ban struct {
	Reason string
}

type Bans struct {
	pool *pgxpool.Pool
}

func NewBans(pool *pgxpool.Pool) *Bans {
	return &Bans{pool: pool}
}

func (r *Bans) Get(ctx context.Context, chatID int64) (*Ban, error) {
	const q = `SELECT reason FROM banned_users WHERE chat_id = $1`
	var reason *string
	err := r.pool.QueryRow(ctx, q, chatID).Scan(&reason)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("check ban: %w", err)
	}
	text := ""
	if reason != nil {
		text = strings.TrimSpace(*reason)
	}
	return &Ban{Reason: text}, nil
}

func (r *Bans) IsBanned(ctx context.Context, chatID int64) (bool, error) {
	ban, err := r.Get(ctx, chatID)
	return ban != nil, err
}
