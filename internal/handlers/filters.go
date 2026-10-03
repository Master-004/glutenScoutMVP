package handlers

import (
	"net/url"

	"github.com/Master-004/glutenscout/internal/models"
)

// Food-type filter chips, shared by the restaurant list and the restaurant
// page's meal list. The filter is a plain ?type= query parameter, so each
// chip is just a link (no JavaScript): it works with Back/Forward, can be
// bookmarked, and HTMX can fetch it as a fragment.

// foodType is one filter option.
type foodType struct {
	Value string // query value and models.Food* constant
	Label string // chip text, e.g. "Desserts"
	Noun  string // for empty-state messages, e.g. "dessert"
}

// foodTypes are the filter options, in chip order (after "All").
var foodTypes = []foodType{
	{models.FoodSalad, "Salads", "salad"},
	{models.FoodGrain, "Grain-based", "grain-based meal"},
	{models.FoodEntree, "Entrees", "entree"},
	{models.FoodDessert, "Desserts", "dessert"},
}

// parseFoodType returns the matching option for a ?type= value, or nil for
// "" and anything unknown (which means "All": a bad value never errors).
func parseFoodType(v string) *foodType {
	for i := range foodTypes {
		if foodTypes[i].Value == v {
			return &foodTypes[i]
		}
	}
	return nil
}

// filterChip is one chip link in the template.
type filterChip struct {
	Label    string
	URL      string
	Selected bool
	HXTarget string // if set, HTMX swaps only this element instead of loading the page
}

// filterChips builds the "All" chip plus one per food type. path and query
// are the page's URL without "type" (e.g. "/restaurants" + city=provo-ut);
// fragment is appended to each link (e.g. "#meals") to keep the list in view;
// hxTarget ("" for none) lets HTMX swap just the results.
func filterChips(path string, query url.Values, current *foodType, fragment, hxTarget string) []filterChip {
	link := func(value string) string {
		q := url.Values{}
		for k, v := range query {
			q[k] = v
		}
		q.Del("type")
		if value != "" {
			q.Set("type", value)
		}
		u := path
		if enc := q.Encode(); enc != "" {
			u += "?" + enc
		}
		return u + fragment
	}
	chips := []filterChip{{Label: "All", URL: link(""), Selected: current == nil, HXTarget: hxTarget}}
	for _, ft := range foodTypes {
		chips = append(chips, filterChip{
			Label: ft.Label, URL: link(ft.Value), HXTarget: hxTarget,
			Selected: current != nil && current.Value == ft.Value,
		})
	}
	return chips
}

// filterMeals keeps the cards whose meal matches the food type (nil = all).
func filterMeals(cards []mealCard, ft *foodType) []mealCard {
	if ft == nil {
		return cards
	}
	var out []mealCard
	for _, c := range cards {
		if c.Meal.FoodType == ft.Value {
			out = append(out, c)
		}
	}
	return out
}
