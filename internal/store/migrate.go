package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Master-004/glutenscout/migrations"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" driver for database/sql
	"github.com/pressly/goose/v3"
)

// Migrate applies any migrations in migrations/ that haven't run yet.
// goose records applied versions in a goose_db_version table, so running
// this repeatedly is safe. It's shared by cmd/migrate and the store tests.
func Migrate(ctx context.Context, databaseURL string) error {
	// goose works with database/sql rather than pgxpool, so open a
	// short-lived database/sql connection just for migrating.
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return fmt.Errorf("create goose provider: %w", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
