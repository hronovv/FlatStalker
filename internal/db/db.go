package db

import (
	"context"
	"fmt"
	"log"

	"flat-stalker/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, database config.Database) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(database.URL)
	if err != nil {
		return nil, fmt.Errorf("db: parse url: %w", err)
	}

	cfg.MaxConns = database.MaxConns
	cfg.MinConns = database.MinConns
	cfg.MaxConnLifetime = database.MaxConnLifetime
	cfg.MaxConnIdleTime = database.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("db: connect: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, database.PingTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return pool, nil
}

func MustNewPool(ctx context.Context, database config.Database) *pgxpool.Pool {
	pool, err := NewPool(ctx, database)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}
