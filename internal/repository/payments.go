package repository

import (
	"context"
	"errors"
	"fmt"

	"flat-stalker/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Payments struct {
	pool *pgxpool.Pool
}

func NewPayments(pool *pgxpool.Pool) *Payments {
	return &Payments{pool: pool}
}

func (r *Payments) Create(ctx context.Context, userID int64, planName string, days, amountKop int) (*models.Payment, error) {
	const q = `
INSERT INTO payments (user_id, payload, plan, days, amount_kop)
VALUES ($1, replace(gen_random_uuid()::text, '-', ''), $2, $3, $4)
RETURNING id, user_id, payload, plan, days, amount_kop, status;
`
	p := &models.Payment{}
	err := r.pool.QueryRow(ctx, q, userID, planName, days, amountKop).Scan(
		&p.ID, &p.UserID, &p.Payload, &p.Plan, &p.Days, &p.AmountKop, &p.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}
	return p, nil
}

func (r *Payments) GetByPayload(ctx context.Context, payload string) (*models.Payment, error) {
	const q = `
SELECT p.id, p.user_id, u.chat_id, p.payload, p.plan, p.days, p.amount_kop, p.status
FROM payments p
JOIN users u ON u.id = p.user_id
WHERE p.payload = $1
`
	p := &models.Payment{}
	err := r.pool.QueryRow(ctx, q, payload).Scan(
		&p.ID, &p.UserID, &p.ChatID, &p.Payload, &p.Plan, &p.Days, &p.AmountKop, &p.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}
	return p, nil
}

func (r *Payments) MarkPaid(ctx context.Context, payload, telegramChargeID string) (*models.Payment, error) {
	const q = `
UPDATE payments p
SET status = 'paid',
    telegram_charge_id = $2,
    paid_at = now()
FROM users u
WHERE p.payload = $1
  AND p.user_id = u.id
  AND p.status = 'pending'
RETURNING p.id, p.user_id, u.chat_id, p.payload, p.plan, p.days, p.amount_kop, p.status;
`
	p := &models.Payment{}
	err := r.pool.QueryRow(ctx, q, payload, telegramChargeID).Scan(
		&p.ID, &p.UserID, &p.ChatID, &p.Payload, &p.Plan, &p.Days, &p.AmountKop, &p.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("mark paid: %w", err)
	}
	return p, nil
}
