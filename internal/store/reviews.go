package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Master-004/glutenscout/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Review input limits (in characters, not bytes, so emoji and accents count as one).
const (
	MaxAuthorNameLen = 40
	MaxReviewBodyLen = 1000
	DefaultAuthor    = "Anonymous"
)

// ReviewInputError describes what's wrong with submitted review fields.
// Each field holds a user-friendly message, or "" if that field is fine.
type ReviewInputError struct {
	Author string
	Body   string
}

func (e *ReviewInputError) Error() string {
	return strings.TrimSpace("invalid review: " + e.Author + " " + e.Body)
}

// NormalizeReview trims the inputs, fills in the default author name, and
// checks the length rules. It returns the cleaned values, or a
// *ReviewInputError that handlers can show next to the form fields.
func NormalizeReview(author, body string) (string, string, error) {
	author = strings.TrimSpace(author)
	body = strings.TrimSpace(body)
	if author == "" {
		author = DefaultAuthor
	}

	var e ReviewInputError
	if n := utf8.RuneCountInString(author); n > MaxAuthorNameLen {
		e.Author = fmt.Sprintf("Display name can be at most %d characters (yours is %d).", MaxAuthorNameLen, n)
	}
	switch n := utf8.RuneCountInString(body); {
	case n == 0:
		e.Body = "Please write a review before submitting."
	case n > MaxReviewBodyLen:
		e.Body = fmt.Sprintf("Reviews can be at most %d characters (yours is %d).", MaxReviewBodyLen, n)
	}
	if e.Author != "" || e.Body != "" {
		return author, body, &e
	}
	return author, body, nil
}

// GetMeal returns one meal with its restaurant's name.
// It returns ErrNotFound if no meal has that id.
func (s *Store) GetMeal(ctx context.Context, id int64) (models.MealDetail, error) {
	const q = `
		SELECT m.id, m.restaurant_id, m.name, m.description, m.photo_url, m.created_at,
		       r.name AS restaurant_name
		FROM meals m
		JOIN restaurants r ON r.id = m.restaurant_id
		WHERE m.id = $1`
	rows, err := s.pool.Query(ctx, q, id)
	if err != nil {
		return models.MealDetail{}, fmt.Errorf("get meal %d: %w", id, err)
	}
	m, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.MealDetail])
	if errors.Is(err, pgx.ErrNoRows) {
		return models.MealDetail{}, ErrNotFound
	}
	if err != nil {
		return models.MealDetail{}, fmt.Errorf("get meal %d: %w", id, err)
	}
	return m, nil
}

// ListReviewsByMeal returns a meal's reviews, newest first.
func (s *Store) ListReviewsByMeal(ctx context.Context, mealID int64) ([]models.Review, error) {
	const q = `
		SELECT id, meal_id, author_name, body, created_at
		FROM reviews
		WHERE meal_id = $1
		ORDER BY created_at DESC, id DESC`
	rows, err := s.pool.Query(ctx, q, mealID)
	if err != nil {
		return nil, fmt.Errorf("list reviews for meal %d: %w", mealID, err)
	}
	reviews, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Review])
	if err != nil {
		return nil, fmt.Errorf("list reviews for meal %d: %w", mealID, err)
	}
	return reviews, nil
}

// CreateReview validates and saves a review. It returns a *ReviewInputError
// for bad input, or ErrNotFound if the meal doesn't exist.
func (s *Store) CreateReview(ctx context.Context, mealID int64, authorName, body string) (models.Review, error) {
	authorName, body, err := NormalizeReview(authorName, body)
	if err != nil {
		return models.Review{}, err
	}
	const q = `
		INSERT INTO reviews (meal_id, author_name, body)
		VALUES ($1, $2, $3)
		RETURNING id, meal_id, author_name, body, created_at`
	rows, err := s.pool.Query(ctx, q, mealID, authorName, body)
	if err != nil {
		return models.Review{}, fmt.Errorf("create review for meal %d: %w", mealID, err)
	}
	review, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Review])
	if isForeignKeyViolation(err) {
		return models.Review{}, ErrNotFound
	}
	if err != nil {
		return models.Review{}, fmt.Errorf("create review for meal %d: %w", mealID, err)
	}
	return review, nil
}

// ToggleUpvote adds this visitor's upvote to a meal, or removes it if they
// already upvoted. It returns the meal's new upvote count and whether the
// visitor now has an upvote on it. Unknown meals return ErrNotFound.
func (s *Store) ToggleUpvote(ctx context.Context, mealID int64, visitorID string) (count int, upvoted bool, err error) {
	// A transaction makes the delete/insert/count happen as one unit, so the
	// count we return reflects this toggle.
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, false, fmt.Errorf("toggle upvote: %w", err)
	}
	defer tx.Rollback(ctx) // no-op after a successful Commit

	// Try removing an existing upvote first. If nothing was removed, the
	// visitor hadn't upvoted yet, so add one.
	tag, err := tx.Exec(ctx, `DELETE FROM meal_upvotes WHERE meal_id = $1 AND visitor_id = $2`, mealID, visitorID)
	if err != nil {
		return 0, false, fmt.Errorf("toggle upvote: %w", err)
	}
	if tag.RowsAffected() == 0 {
		// ON CONFLICT DO NOTHING: if a double-tap raced us, don't error.
		_, err = tx.Exec(ctx,
			`INSERT INTO meal_upvotes (meal_id, visitor_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			mealID, visitorID)
		if isForeignKeyViolation(err) {
			return 0, false, ErrNotFound
		}
		if err != nil {
			return 0, false, fmt.Errorf("toggle upvote: %w", err)
		}
		upvoted = true
	}

	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM meal_upvotes WHERE meal_id = $1`, mealID).Scan(&count); err != nil {
		return 0, false, fmt.Errorf("toggle upvote: count: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, false, fmt.Errorf("toggle upvote: commit: %w", err)
	}
	return count, upvoted, nil
}

// MealStats returns upvote and review counts for several meals at once, plus
// whether visitorID has upvoted each one. One query for the whole list avoids
// running a query per meal card. The map is keyed by meal id.
func (s *Store) MealStats(ctx context.Context, mealIDs []int64, visitorID string) (map[int64]models.MealStats, error) {
	const q = `
		SELECT m.id AS meal_id,
		       (SELECT COUNT(*) FROM meal_upvotes u WHERE u.meal_id = m.id) AS upvotes,
		       (SELECT COUNT(*) FROM reviews r WHERE r.meal_id = m.id) AS reviews,
		       EXISTS (SELECT 1 FROM meal_upvotes u WHERE u.meal_id = m.id AND u.visitor_id = $2) AS upvoted
		FROM meals m
		WHERE m.id = ANY($1)`
	rows, err := s.pool.Query(ctx, q, mealIDs, visitorID)
	if err != nil {
		return nil, fmt.Errorf("meal stats: %w", err)
	}
	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.MealStats])
	if err != nil {
		return nil, fmt.Errorf("meal stats: %w", err)
	}
	stats := make(map[int64]models.MealStats, len(list))
	for _, st := range list {
		stats[st.MealID] = st
	}
	return stats, nil
}

// isForeignKeyViolation reports whether err is Postgres error 23503, which
// here means the meal id we referenced doesn't exist.
func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
