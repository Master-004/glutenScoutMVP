package handlers

import (
	"strings"
	"testing"

	"github.com/Master-004/glutenscout/internal/models"
)

// Tweak 6: "Good for mixed groups" on the list and the restaurant page.
func TestGoodForMixedGroupsLabel(t *testing.T) {
	h := newTestServer(t, seededStore())

	list := get(t, h, "/restaurants?city=provo-ut").Body.String()
	pizza := strings.Index(list, "Canyon Hearth Pizza")
	solo := strings.Index(list, "Solo Cafe")
	label := strings.Index(list, "Good for mixed groups")
	// Exactly one label, on Canyon Hearth Pizza's row (before Solo Cafe's).
	if strings.Count(list, "Good for mixed groups") != 1 || label < pizza || label > solo {
		t.Error("only group-friendly restaurants should show the label, on their own row")
	}

	if rest := get(t, h, "/restaurants/10").Body.String(); !strings.Contains(rest, "Good for mixed groups") ||
		!strings.Contains(rest, "Plenty for people who aren't gluten-free too.") {
		t.Error("restaurant page should show the group label when true")
	}
	if rest := get(t, h, "/restaurants/11").Body.String(); strings.Contains(rest, "Good for mixed groups") {
		t.Error("restaurant page shouldn't show the label when false")
	}
	// The section with options for people who aren't GF is still there.
	st := seededStore()
	st.meals[10] = append(st.meals[10], models.Meal{ID: 104, RestaurantID: 10, Name: "Garlic Knots", GFStatus: models.ContainsGF})
	if rest := get(t, newTestServer(t, st), "/restaurants/10").Body.String(); !strings.Contains(rest, "Contains gluten (1)") {
		t.Error("the contains-gluten section should still show")
	}
}

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
	if strings.Contains(salads, " GF meal") || !strings.Contains(salads, " GF salad") {
		t.Error("a filtered row should count the filtered type, not meals")
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

func TestRestaurantBackLinkKeepsFilter(t *testing.T) {
	h := newTestServer(t, seededStore())
	if body := get(t, h, "/restaurants/10?type=salad").Body.String(); !strings.Contains(body, `href="/restaurants?city=provo-ut&type=salad#results"`) {
		t.Error("back link should keep the food-type filter")
	}
	if body := get(t, h, "/restaurants/10").Body.String(); !strings.Contains(body, `href="/restaurants?city=provo-ut#results"`) {
		t.Error("back link should have no filter when none is active")
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
