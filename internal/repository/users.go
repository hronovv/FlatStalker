package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"flat-stalker/internal/cache"
	"flat-stalker/internal/models"
	"flat-stalker/internal/plan"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrBanned = errors.New("user is banned")

type Users struct {
	pool   *pgxpool.Pool
	byChat *cache.LRU[models.User]
}

func NewUsers(pool *pgxpool.Pool) *Users {
	return &Users{
		pool:   pool,
		byChat: cache.NewLRU[models.User](cache.DefaultSize, cache.DefaultTTL),
	}
}

func (r *Users) GetByChatID(ctx context.Context, chatID int64) (*models.User, error) {
	if cached, ok := r.byChat.Get(chatID); ok {
		user := cached
		return &user, nil
	}

	const q = `SELECT id, chat_id, plan, plan_expires_at FROM users WHERE chat_id = $1`
	user := &models.User{}
	err := r.pool.QueryRow(ctx, q, chatID).Scan(&user.ID, &user.ChatID, &user.Plan, &user.PlanExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	user.Plan = plan.Effective(user.Plan, user.PlanExpiresAt, time.Now())
	r.byChat.Put(chatID, *user)
	return user, nil
}

// CreateByChatID inserts a user or returns the existing one (chat_id is unique).
func (r *Users) CreateByChatID(ctx context.Context, chatID int64) (*models.User, error) {
	var banned bool
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM banned_users WHERE chat_id = $1)`, chatID).Scan(&banned); err != nil {
		return nil, fmt.Errorf("check ban: %w", err)
	}
	if banned {
		return nil, ErrBanned
	}

	const q = `
INSERT INTO users (chat_id)
VALUES ($1)
ON CONFLICT (chat_id) DO UPDATE SET chat_id = EXCLUDED.chat_id
RETURNING id, chat_id, plan, plan_expires_at;
`
	user := &models.User{}
	err := r.pool.QueryRow(ctx, q, chatID).Scan(&user.ID, &user.ChatID, &user.Plan, &user.PlanExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	user.Plan = plan.Effective(user.Plan, user.PlanExpiresAt, time.Now())
	r.byChat.Put(chatID, *user)
	return user, nil
}

func (r *Users) ApplyPlan(ctx context.Context, chatID int64, newPlan string, days int) error {
	newPlan = plan.Normalize(newPlan)
	if newPlan == plan.Free || days < 1 {
		return fmt.Errorf("invalid plan")
	}

	const q = `
UPDATE users
SET plan = $2,
    plan_expires_at = CASE
        WHEN plan = $2 AND plan_expires_at IS NOT NULL AND plan_expires_at > now()
            THEN plan_expires_at + ($3 * INTERVAL '1 day')
        ELSE now() + ($3 * INTERVAL '1 day')
    END
WHERE chat_id = $1
RETURNING id, chat_id, plan, plan_expires_at;
`
	user := models.User{}
	err := r.pool.QueryRow(ctx, q, chatID, newPlan, days).Scan(&user.ID, &user.ChatID, &user.Plan, &user.PlanExpiresAt)
	if err != nil {
		return fmt.Errorf("apply plan: %w", err)
	}
	user.Plan = plan.Effective(user.Plan, user.PlanExpiresAt, time.Now())
	r.byChat.Put(chatID, user)
	return nil
}
