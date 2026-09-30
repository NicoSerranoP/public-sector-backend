package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"public-sector-backend/internal/config"
)

var database *pgxpool.Pool

func Initialize(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, config.DATABASE_URL)

	if err != nil {
		return fmt.Errorf("create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	database = pool

	return nil
}

func Get() *pgxpool.Pool {
	if database == nil {
		panic("database not initialized. Call Initialize first")
	}

	return database
}
