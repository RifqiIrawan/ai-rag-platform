package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates a connection pool without blocking on an initial ping,
// so the service can start even if Postgres is briefly unavailable.
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, databaseURL)
}
