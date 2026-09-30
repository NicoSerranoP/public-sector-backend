package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"public-sector-backend/internal/config"
)

var database *pgxpool.Pool

func createTables(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS voters (
			id CHAR(10) PRIMARY KEY,
			address CHAR(42) NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)

	return err
}

func Initialize(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, config.DATABASE_URL)

	if err != nil {
		return fmt.Errorf("create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	if err := createTables(ctx, pool); err != nil {
		pool.Close()
		return fmt.Errorf("create tables: %w", err)
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
