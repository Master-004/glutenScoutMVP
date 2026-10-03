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

// ratingColumns computes a food rating from reviews joined as "rv". Reviews
// with rating 0 (none given) are left out of both the average and the count.
// COALESCE turns "no ratings at all" (NULL) into 0.
const ratingColumns = `
	COALESCE(AVG(rv.rating) FILTER (WHERE rv.rating > 0), 0)::float8 AS avg_rating,
	COUNT(rv.id) FILTER (WHERE rv.rating > 0) AS rating_count,
	COUNT(rv.id) AS review_count`

// GetRestaurant returns one restaurant with its city's name and slug, and its
// food rating across all of its meals' reviews.
// It returns ErrNotFound if no restaurant has that id.
func (s *Store) GetRestaurant(ctx context.Context, id int64) (models.RestaurantDetail, error) {
	q := `
		SELECT r.id, r.city_id, r.name, r.address, r.cuisine, r.description, r.phone, r.good_for_groups, r.created_at,
		       c.name AS city_name, c.slug AS city_slug,` + ratingColumns + `
		FROM restaurants r
		JOIN cities c ON c.id = r.city_id
		LEFT JOIN meals m ON m.restaurant_id = r.id
		LEFT JOIN reviews rv ON rv.meal_id = m.id
		WHERE r.id = $1
		GROUP BY r.id, c.id`
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
// a count of its GF options (not meals with gluten or a cross-contact risk),
// its food rating, and its last-verified date.
//
// foodType ("" for all) keeps only restaurants with at least one GF option of
// that type: "Desserts" means "somewhere I can get a gluten-free dessert".
func (s *Store) ListRestaurantsByCity(ctx context.Context, cityID int64, foodType string) ([]models.RestaurantSummary, error) {
	// LEFT JOIN (not JOIN) keeps restaurants with zero meals, reviews or
	// safety profile. Joining reviews repeats each meal once per review, so
	// meals are counted with DISTINCT.
	q := `
		SELECT r.id, r.city_id, r.name, r.address, r.cuisine, r.description, r.phone, r.good_for_groups, r.created_at,
		       COUNT(DISTINCT m.id) FILTER (WHERE m.gf_status = 'gf') AS meal_count,
		       sp.verified_on,` + ratingColumns + `
		FROM restaurants r
		LEFT JOIN meals m ON m.restaurant_id = r.id
		LEFT JOIN reviews rv ON rv.meal_id = m.id
		LEFT JOIN safety_profiles sp ON sp.restaurant_id = r.id
		WHERE r.city_id = $1
		  AND ($2 = '' OR EXISTS (
		        SELECT 1 FROM meals f
		        WHERE f.restaurant_id = r.id AND f.gf_status = 'gf' AND f.food_type = $2))
		GROUP BY r.id, sp.verified_on
		ORDER BY r.name, r.id`
	rows, err := s.pool.Query(ctx, q, cityID, foodType)
	if err != nil {
		return nil, fmt.Errorf("list restaurants for city %d: %w", cityID, err)
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.RestaurantSummary])
	if err != nil {
		return nil, fmt.Errorf("list restaurants for city %d: %w", cityID, err)
	}
	return list, nil
}

// ListRestaurantPhotos returns a restaurant's photos, grouped by kind and in
// display order.
func (s *Store) ListRestaurantPhotos(ctx context.Context, restaurantID int64) ([]models.RestaurantPhoto, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT kind, url, alt FROM restaurant_photos
		WHERE restaurant_id = $1
		ORDER BY kind, position, id`, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list photos for restaurant %d: %w", restaurantID, err)
	}
	photos, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.RestaurantPhoto])
	if err != nil {
		return nil, fmt.Errorf("list photos for restaurant %d: %w", restaurantID, err)
	}
	return photos, nil
}

// GetSafetyProfile returns a restaurant's safety profile and its precaution
// badges in display order. A restaurant without a profile gets a zero
// SafetyProfile (Found == false), not an error: "no information yet" is a
// normal state the page shows.
func (s *Store) GetSafetyProfile(ctx context.Context, restaurantID int64) (models.SafetyProfile, error) {
	p := models.SafetyProfile{}
	err := s.pool.QueryRow(ctx, `
		SELECT verified_on, answers_updated_on, conflict_note
		FROM safety_profiles WHERE restaurant_id = $1`, restaurantID).
		Scan(&p.VerifiedOn, &p.AnswersUpdatedOn, &p.ConflictNote)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.SafetyProfile{}, nil
	}
	if err != nil {
		return models.SafetyProfile{}, fmt.Errorf("get safety profile %d: %w", restaurantID, err)
	}
	p.Found = true

	rows, err := s.pool.Query(ctx, `
		SELECT label, status FROM precautions
		WHERE restaurant_id = $1
		ORDER BY position, id`, restaurantID)
	if err != nil {
		return models.SafetyProfile{}, fmt.Errorf("list precautions %d: %w", restaurantID, err)
	}
	if p.Precautions, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.Precaution]); err != nil {
		return models.SafetyProfile{}, fmt.Errorf("list precautions %d: %w", restaurantID, err)
	}
	return p, nil
}
