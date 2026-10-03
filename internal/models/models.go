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
	Phone       string    `db:"phone"` // for the call-ahead script; "" if unknown
	CreatedAt   time.Time `db:"created_at"`
}

// Rating is the average food rating (1–5 stars) across reviews that gave one.
// Stars are for food and experience only, never gluten safety.
type Rating struct {
	AvgRating   float64 `db:"avg_rating"`   // 0 when nobody has rated yet
	RatingCount int     `db:"rating_count"` // how many reviews gave a rating
	ReviewCount int     `db:"review_count"` // all reviews, rated or not (0 = "No reviews yet")
}

// RestaurantSummary is a restaurant plus how many GF meals it lists, its food
// rating, and when its safety profile was last verified (nil = no profile),
// used for the rows on the city results list.
type RestaurantSummary struct {
	Restaurant            // embedded: its fields are scanned as if declared here
	Rating                // embedded too: avg_rating, rating_count
	MealCount  int        `db:"meal_count"`
	VerifiedOn *time.Time `db:"verified_on"` // a pointer so a missing profile (NULL) scans as nil
}

// RestaurantDetail is a restaurant plus its city's name and slug (for the
// "Back to {City} restaurants" link) and its food rating.
type RestaurantDetail struct {
	Restaurant
	Rating
	CityName string `db:"city_name"`
	CitySlug string `db:"city_slug"`
}

// Precaution evidence states. They match the CHECK constraint in the database.
const (
	EvidenceConfirmed  = "confirmed"  // restaurant and recent reviewers agree
	EvidenceConflict   = "conflict"   // reports disagree
	EvidenceUnverified = "unverified" // restaurant-only or out of date
)

// Precaution is one kitchen practice and how sure we are about it,
// shown as a PrecautionBadge.
type Precaution struct {
	Label  string `db:"label"`  // names a practice, e.g. "Separate prep area"
	Status string `db:"status"` // one of the Evidence* constants
}

// SafetyProfile is the dated summary at the top of a restaurant page.
// A zero SafetyProfile (Found == false) means we have no information yet.
type SafetyProfile struct {
	Found            bool
	VerifiedOn       time.Time
	AnswersUpdatedOn *time.Time // when the restaurant last answered; nil = never
	ConflictNote     string     // disagreement in words; "" = none
	Precautions      []Precaution
}

// Meal gluten statuses. They match the CHECK constraint in the database.
const (
	GFOption     = "gf"            // made without gluten
	CrossContact = "cross_contact" // no gluten ingredients, but shared equipment
	ContainsGF   = "contains"      // contains gluten
)

// Meal is a dish a restaurant offers. GFStatus says whether it's a
// gluten-free option; most listed meals are, but menus also have dishes
// with gluten or a cross-contact risk.
type Meal struct {
	ID           int64     `db:"id"`
	RestaurantID int64     `db:"restaurant_id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	PhotoURL     string    `db:"photo_url"` // empty means "show the placeholder image"
	GFStatus     string    `db:"gf_status"` // one of GFOption, CrossContact, ContainsGF
	GFNote       string    `db:"gf_note"`   // why, in words; "" if nothing to add
	CreatedAt    time.Time `db:"created_at"`
}

// MealDetail is a meal plus its restaurant's name and city, for the meal
// page's links back to the restaurant and to the city's restaurant list.
type MealDetail struct {
	Meal
	RestaurantName string `db:"restaurant_name"`
	CityName       string `db:"city_name"`
	CitySlug       string `db:"city_slug"`
}

// MealStats is the community activity on one meal, as seen by one visitor.
type MealStats struct {
	MealID  int64 `db:"meal_id"`
	Rating        // avg_rating, rating_count
	Upvotes int   `db:"upvotes"`
	Reviews int   `db:"reviews"`
	Upvoted bool  `db:"upvoted"` // has the current visitor upvoted this meal?
}

// Reviewer types: who wrote a review. They match the CHECK constraint in the
// database ("" is only for reviews posted before reviewer types existed).
const (
	ReviewerCeliac     = "celiac"
	ReviewerStrict     = "strict"
	ReviewerSensitive  = "sensitive"
	ReviewerRestaurant = "restaurant"
)

// Review is a community member's written review of a meal.
type Review struct {
	ID           int64     `db:"id"`
	MealID       int64     `db:"meal_id"`
	AuthorName   string    `db:"author_name"`
	Body         string    `db:"body"`
	ReviewerType string    `db:"reviewer_type"` // one of the Reviewer* constants
	Rating       int       `db:"rating"`        // 1–5 food stars; 0 = none given
	CreatedAt    time.Time `db:"created_at"`
}
