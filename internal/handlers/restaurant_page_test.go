package handlers

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestRestaurantPageFound(t *testing.T) {
	rec := get(t, newTestServer(t, seededStore()), "/restaurants/10")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		// Header info.
		"<h1>Canyon Hearth Pizza</h1>", "Pizza", "410 Timber Lane", "Wood-fired pizza.",
		// Back link returns to the right city list.
		`href="/restaurants?city=provo-ut"`, "Back to Provo restaurants",
		// Meal cards with name, description, and photo.
		"GF Margherita", "Tomato and basil.",
		// Images are lazy and use the meal name as alt text.
		`src="/static/images/meals/pizza.svg" alt="GF Margherita" loading="lazy"`,
		disclaimer,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("restaurant page missing %q", want)
		}
	}
}

func TestRestaurantPageMissingPhotoUsesPlaceholder(t *testing.T) {
	body := get(t, newTestServer(t, seededStore()), "/restaurants/10").Body.String()
	// Harvest Salad has no photo_url, so it should get the placeholder.
	want := `src="/static/images/meals/placeholder.svg" alt="Harvest Salad"`
	if !strings.Contains(body, want) {
		t.Error("meal without a photo should show the placeholder image")
	}
	if strings.Contains(body, `src=""`) {
		t.Error("no image should have an empty src")
	}
}

func TestRestaurantPageNoMeals(t *testing.T) {
	rec := get(t, newTestServer(t, seededStore()), "/restaurants/11")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "No GF meals listed yet.") {
		t.Error("missing empty-state message")
	}
}

func TestRestaurantPageNotFoundAndBadIDs(t *testing.T) {
	h := newTestServer(t, seededStore())
	for _, path := range []string{
		"/restaurants/999",                  // unknown id
		"/restaurants/abc",                  // not a number
		"/restaurants/-1",                   // negative
		"/restaurants/0",                    // ids start at 1
		"/restaurants/1.5",                  // not an integer
		"/restaurants/99999999999999999999", // too big for int64
	} {
		rec := get(t, h, path)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s: status = %d, want 404", path, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "Page not found") {
			t.Errorf("%s: should show the friendly 404 page", path)
		}
	}
}

func TestRestaurantPageStoreError(t *testing.T) {
	st := seededStore()
	st.err = errors.New("connection reset")
	rec := get(t, newTestServer(t, st), "/restaurants/10")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "connection reset") {
		t.Error("raw error leaked to the page")
	}
}
