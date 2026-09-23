package store

import (
	"context"
	"os"
	"testing"
)

// newTestStore connects to the test database, applies migrations, and empties
// every table so each test starts from a clean slate.
//
// Store tests need a real Postgres, so they skip unless TEST_DATABASE_URL is
// set (see README). That keeps `go test ./...` working without Docker.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping database test")
	}
	ctx := context.Background()
	if err := Migrate(ctx, url); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	s, err := New(ctx, url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(s.Close)

	// CASCADE also clears tables that reference these (meals, reviews, ...).
	if _, err := s.pool.Exec(ctx, `TRUNCATE cities RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	return s
}

// mustInsert runs a setup INSERT and returns the id it produces
// (statements should end with RETURNING id).
func mustInsert(t *testing.T, s *Store, sql string, args ...any) int64 {
	t.Helper()
	var id int64
	if err := s.pool.QueryRow(context.Background(), sql, args...).Scan(&id); err != nil {
		t.Fatalf("insert %q: %v", sql, err)
	}
	return id
}

func insertCity(t *testing.T, s *Store, name, slug string) int64 {
	t.Helper()
	return mustInsert(t, s, `INSERT INTO cities (name, state, slug) VALUES ($1, 'UT', $2) RETURNING id`, name, slug)
}

func insertRestaurant(t *testing.T, s *Store, cityID int64, name string) int64 {
	t.Helper()
	return mustInsert(t, s,
		`INSERT INTO restaurants (city_id, name, address, cuisine) VALUES ($1, $2, '1 Test St', 'Test') RETURNING id`,
		cityID, name)
}

func insertMeal(t *testing.T, s *Store, restaurantID int64, name, photoURL string) int64 {
	t.Helper()
	return mustInsert(t, s,
		`INSERT INTO meals (restaurant_id, name, description, photo_url) VALUES ($1, $2, 'Tasty', $3) RETURNING id`,
		restaurantID, name, photoURL)
}
