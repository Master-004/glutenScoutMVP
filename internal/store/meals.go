package store

import (
	"context"
	"fmt"

	"github.com/Master-004/glutenscout/internal/models"
	"github.com/jackc/pgx/v5"
)

// ListMealsByRestaurant returns a restaurant's GF meals sorted by name.
func (s *Store) ListMealsByRestaurant(ctx context.Context, restaurantID int64) ([]models.Meal, error) {
	const q = `
		SELECT id, restaurant_id, name, description, photo_url, created_at
		FROM meals
		WHERE restaurant_id = $1
		ORDER BY name, id`
	rows, err := s.pool.Query(ctx, q, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list meals for restaurant %d: %w", restaurantID, err)
	}
	meals, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Meal])
	if err != nil {
		return nil, fmt.Errorf("list meals for restaurant %d: %w", restaurantID, err)
	}
	return meals, nil
}
