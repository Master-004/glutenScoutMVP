package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/Master-004/glutenscout/internal/models"
	"github.com/jackc/pgx/v5"
)

// ListCities returns every city, sorted by name, for the city dropdown.
func (s *Store) ListCities(ctx context.Context) ([]models.City, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, state, slug FROM cities ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}
	// CollectRows reads every row and closes rows for us.
	cities, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.City])
	if err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}
	return cities, nil
}

// GetCityBySlug finds a city by its URL slug (e.g. "denver-co").
// It returns ErrNotFound if no city has that slug.
func (s *Store) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, state, slug FROM cities WHERE slug = $1`, slug)
	if err != nil {
		return models.City{}, fmt.Errorf("get city %q: %w", slug, err)
	}
	city, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.City])
	if errors.Is(err, pgx.ErrNoRows) {
		return models.City{}, ErrNotFound
	}
	if err != nil {
		return models.City{}, fmt.Errorf("get city %q: %w", slug, err)
	}
	return city, nil
}

// GetRestaurant returns one restaurant with its city's name and slug.
// It returns ErrNotFound if no restaurant has that id.
func (s *Store) GetRestaurant(ctx context.Context, id int64) (models.RestaurantDetail, error) {
	const q = `
		SELECT r.id, r.city_id, r.name, r.address, r.cuisine, r.description, r.created_at,
		       c.name AS city_name, c.slug AS city_slug
		FROM restaurants r
		JOIN cities c ON c.id = r.city_id
		WHERE r.id = $1`
	rows, err := s.pool.Query(ctx, q, id)
	if err != nil {
		return models.RestaurantDetail{}, fmt.Errorf("get restaurant %d: %w", id, err)
	}
	r, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.RestaurantDetail])
	if errors.Is(err, pgx.ErrNoRows) {
		return models.RestaurantDetail{}, ErrNotFound
	}
	if err != nil {
		return models.RestaurantDetail{}, fmt.Errorf("get restaurant %d: %w", id, err)
	}
	return r, nil
}

// ListRestaurantsByCity returns a city's restaurants sorted by name, each with
// a count of its GF meals.
func (s *Store) ListRestaurantsByCity(ctx context.Context, cityID int64) ([]models.RestaurantSummary, error) {
	// LEFT JOIN (not JOIN) keeps restaurants with zero meals; COUNT(m.id)
	// counts only matched meals, so those restaurants get 0 instead of 1.
	const q = `
		SELECT r.id, r.city_id, r.name, r.address, r.cuisine, r.description, r.created_at,
		       COUNT(m.id) AS meal_count
		FROM restaurants r
		LEFT JOIN meals m ON m.restaurant_id = r.id
		WHERE r.city_id = $1
		GROUP BY r.id
		ORDER BY r.name, r.id`
	rows, err := s.pool.Query(ctx, q, cityID)
	if err != nil {
		return nil, fmt.Errorf("list restaurants for city %d: %w", cityID, err)
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.RestaurantSummary])
	if err != nil {
		return nil, fmt.Errorf("list restaurants for city %d: %w", cityID, err)
	}
	return list, nil
}
