// Package models holds plain Go structs that mirror our database tables.
// They have no behavior — the store fills them in and templates read them.
//
// The `db:"..."` tags tell pgx which SQL column goes into which field when
// the store scans rows with pgx.RowToStructByName.
package models

import "time"

// City is a place users can pick to find restaurants.
type City struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	State string `db:"state"`
	Slug  string `db:"slug"` // URL-friendly id, e.g. "denver-co"
}

// Restaurant is a place with gluten-free options.
type Restaurant struct {
	ID          int64     `db:"id"`
	CityID      int64     `db:"city_id"`
	Name        string    `db:"name"`
	Address     string    `db:"address"`
	Cuisine     string    `db:"cuisine"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

// RestaurantSummary is a restaurant plus how many GF meals it lists,
// used for the cards on the city results list.
type RestaurantSummary struct {
	Restaurant     // embedded: its fields are scanned as if declared here
	MealCount  int `db:"meal_count"`
}

// RestaurantDetail is a restaurant plus its city's name and slug, used on the
// restaurant page for the "Back to {City} restaurants" link.
type RestaurantDetail struct {
	Restaurant
	CityName string `db:"city_name"`
	CitySlug string `db:"city_slug"`
}

// Meal is a gluten-free dish a restaurant offers.
type Meal struct {
	ID           int64     `db:"id"`
	RestaurantID int64     `db:"restaurant_id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	PhotoURL     string    `db:"photo_url"` // empty means "show the placeholder image"
	CreatedAt    time.Time `db:"created_at"`
}

// MealDetail is a meal plus its restaurant's name, for the meal page.
type MealDetail struct {
	Meal
	RestaurantName string `db:"restaurant_name"`
}

// MealStats is the community activity on one meal, as seen by one visitor.
type MealStats struct {
	MealID  int64 `db:"meal_id"`
	Upvotes int   `db:"upvotes"`
	Reviews int   `db:"reviews"`
	Upvoted bool  `db:"upvoted"` // has the current visitor upvoted this meal?
}

// Review is a community member's written review of a meal.
type Review struct {
	ID         int64     `db:"id"`
	MealID     int64     `db:"meal_id"`
	AuthorName string    `db:"author_name"`
	Body       string    `db:"body"`
	CreatedAt  time.Time `db:"created_at"`
}
