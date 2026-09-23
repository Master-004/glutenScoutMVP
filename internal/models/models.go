// Package models holds plain Go structs that mirror our database tables.
// They have no behavior — the store fills them in and templates read them.
package models

import "time"

// City is a place users can pick to find restaurants.
type City struct {
	ID    int64
	Name  string
	State string
	Slug  string // URL-friendly id, e.g. "denver-co"
}

// Restaurant is a place with gluten-free options.
type Restaurant struct {
	ID          int64
	CityID      int64
	Name        string
	Address     string
	Cuisine     string
	Description string
	CreatedAt   time.Time
}

// Meal is a gluten-free dish a restaurant offers.
type Meal struct {
	ID           int64
	RestaurantID int64
	Name         string
	Description  string
	PhotoURL     string // empty means "show the placeholder image"
	CreatedAt    time.Time
}
