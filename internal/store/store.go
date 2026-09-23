// Package store is the only place that talks to the database.
// Handlers call Store methods and never write SQL themselves, which keeps
// SQL in one spot and makes handlers easy to test with a fake store.
package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a lookup (by id, slug, ...) matches no row.
// Handlers check for it with errors.Is to decide to show the 404 page.
var ErrNotFound = errors.New("not found")

// Store wraps a pgx connection pool. A pool (instead of one connection) lets
// many web requests query the database at the same time.
type Store struct {
	pool *pgxpool.Pool
}

// New connects to Postgres and verifies the connection works.
func New(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return &Store{pool: pool}, nil
}

// Ping checks the database is reachable (used by /healthz).
func (s *Store) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

// Close releases all pooled connections. Call it once on shutdown.
func (s *Store) Close() {
	s.pool.Close()
}
