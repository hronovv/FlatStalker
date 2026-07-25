package repository

import (
	"context"
	"fmt"

	"flat-stalker/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Users struct {
	pool *pgxpool.Pool
}

func NewUsers(pool *pgxpool.Pool) *Users {
	return &Users{pool: pool}
}

// CreateByChatID inserts a user or returns the existing one (chat_id is unique).
func (r *Users) CreateByChatID(ctx context.Context, chatID int64) (*models.User, error) {
	const q = `
INSERT INTO users (chat_id)
VALUES ($1)
ON CONFLICT (chat_id) DO UPDATE SET chat_id = EXCLUDED.chat_id
RETURNING id, chat_id;
`
	user := &models.User{}
	err := r.pool.QueryRow(ctx, q, chatID).Scan(&user.ID, &user.ChatID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}
