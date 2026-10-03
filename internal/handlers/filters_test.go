package handlers

import (
	"strings"
	"testing"
)

// Tweak 5: food-type filter chips on the restaurant list and the meal list.

func TestRestaurantListFoodTypeFilter(t *testing.T) {
	h := newTestServer(t, seededStore())

	all := get(t, h, "/restaurants?city=provo-ut").Body.String()
	for _, want := range []string{
		`aria-label="Filter by food type"`,
		`href="/restaurants?city=provo-ut"`,
		`hx-get="/restaurants?city=provo-ut" hx-target="#results" hx-push-url="true"`,
		`href="/restaurants?city=provo-ut&amp;type=dessert"`,
		">Salads</a>", ">Grain-based</a>", ">Entrees</a>", ">Desserts</a>",
		"Canyon Hearth Pizza", "Solo Cafe",
	} {
		if !strings.Contains(all, want) {
			t.Errorf("unfiltered list missing %q", want)
		}
	}
	if !strings.Contains(all, `aria-current="true">All</a>`) {
		t.Error(`"All" should be selected with no filter`)
	}

	// Salads: only Canyon Hearth Pizza has a GF salad.
	salads := get(t, h, "/restaurants?city=provo-ut&type=salad").Body.String()
	if !strings.Contains(salads, "Canyon Hearth Pizza") || strings.Contains(salads, "Solo Cafe") {
		t.Error("Salads filter should show only restaurants with a GF salad")
	}
	if !strings.Contains(salads, `aria-current="true">Salads</a>`) {
		t.Error("the chosen chip should be marked selected")
	}
	if !strings.Contains(salads, `href="/restaurants/10?type=salad#meals"`) {
		t.Error("restaurant links should carry the filter into the meal list")
	}

	// Desserts: none here, so a message and a way back to all.
	desserts := get(t, h, "/restaurants?city=provo-ut&type=dessert").Body.String()
	if !strings.Contains(desserts, "No restaurants in Provo have a gluten-free dessert listed yet.") ||
		!strings.Contains(desserts, "Show all restaurants") {
		t.Error("an empty filter should explain itself and offer Show all")
	}

	// An unknown type is ignored (treated as All), never an error.
	if bad := get(t, h, "/restaurants?city=provo-ut&type=pizza"); bad.Code != 200 || !strings.Contains(bad.Body.String(), "Solo Cafe") {
		t.Error("unknown type should fall back to All")
	}
}

func TestMealListFoodTypeFilter(t *testing.T) {
	h := newTestServer(t, seededStore())

	all := get(t, h, "/restaurants/10").Body.String()
	if !strings.Contains(all, `<div id="meals">`) || !strings.Contains(all, `href="/restaurants/10?type=dessert#meals"`) {
		t.Error("restaurant page should have food-type chips linking back to #meals")
	}
	if !strings.Contains(all, "<h3>GF Margherita") || !strings.Contains(all, "<h3>Harvest Salad") {
		t.Fatal("unfiltered meal list should show every meal")
	}

	salads := get(t, h, "/restaurants/10?type=salad").Body.String()
	if !strings.Contains(salads, "<h3>Harvest Salad") || strings.Contains(salads, "<h3>GF Margherita") {
		t.Error("Salads should show only salads")
	}

	desserts := get(t, h, "/restaurants/10?type=dessert").Body.String()
	if !strings.Contains(desserts, "No gluten-free dessert listed here yet.") || !strings.Contains(desserts, "Show all meals") {
		t.Error("an empty meal filter should explain itself and offer Show all")
	}
}
