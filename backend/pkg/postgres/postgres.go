package postgres

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func New(dsn string) (*sqlx.DB, error) {
	pool, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	pool.SetMaxOpenConns(25)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(5 * time.Minute)
	pool.SetConnMaxIdleTime(5 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = pool.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("postgres.New: ping: %w", err)
	}
	return pool, nil
}
