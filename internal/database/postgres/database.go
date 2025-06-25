// Package database consist of function to create
// database connection pool using pgxpool
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(dbURL string, maxConn int) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Database URL: %w", err)
	}

	config.MaxConns = int32(maxConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create db connection pool: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
