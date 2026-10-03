package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/Master-004/glutenscout/internal/models"
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

func TestRestaurantPageSafetyProfile(t *testing.T) {
	body := get(t, newTestServer(t, seededStore()), "/restaurants/10").Body.String()
	for _, want := range []string{
		// Every safety claim has a date.
		"Last verified Sep 12, 2026", "updated Aug 30, 2026",
		// Badges carry icon + words + color: a class per state plus hidden words.
		`class="badge badge-confirmed"`, "Confirmed: </span>Separate prep area",
		`class="badge badge-conflict"`, "Reports disagree: </span>Shared oven?",
		`class="badge badge-unverified"`, "Unverified: </span>GF hoisin on request",
		// Conflicts in words.
		"2 reports disagree about the oven.",
		// The call-ahead script is a sheet opened by a button, with a nudge
		// because something isn't confirmed.
		`popovertarget="call-ahead">Call ahead</button>`, `id="call-ahead" class="sheet" popover`,
		"Some details aren't confirmed.",
		"Ask about shared oven, GF hoisin on request.",
		`href="tel:+18015550101"`, "Call (801) 555-0101",
		// Food rating, rounded to half stars, always with the count.
		"4.5", "3 ratings",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("restaurant page missing %q", want)
		}
	}
	// Copy rule: never call anything "safe" ("Safety profile" is fine).
	if regexp.MustCompile(`(?i)\bsafe\b`).MatchString(body) {
		t.Error(`page copy should never say "safe"`)
	}
}

func TestRestaurantPageGroupsMealsByGlutenStatus(t *testing.T) {
	st := seededStore()
	st.meals[10] = append(st.meals[10],
		models.Meal{ID: 102, RestaurantID: 10, Name: "House Fries", GFStatus: models.CrossContact, GFNote: "Shared fryer."},
		models.Meal{ID: 103, RestaurantID: 10, Name: "Garlic Knots", GFStatus: models.ContainsGF, GFNote: "Wheat dough."},
	)
	body := get(t, newTestServer(t, st), "/restaurants/10").Body.String()

	gf := strings.Index(body, "Gluten-free options")
	cross := strings.Index(body, "Cross-contact risk</h2>")
	contains := strings.Index(body, "Contains gluten (1)")
	if gf < 0 || cross < 0 || contains < 0 || !(gf < cross && cross < contains) {
		t.Fatalf("want sections in order GF, cross-contact, contains; got %d, %d, %d", gf, cross, contains)
	}
	// Each meal lands in its own group.
	if i := strings.Index(body, "GF Margherita"); i < gf || i > cross {
		t.Error("GF Margherita should be in the GF section")
	}
	if i := strings.Index(body, "House Fries"); i < cross || i > contains {
		t.Error("House Fries should be in the cross-contact section")
	}
	if i := strings.Index(body, "Garlic Knots"); i < contains {
		t.Error("Garlic Knots should be in the contains-gluten section")
	}
	// Status in icon + words + color, plus the reason.
	for _, want := range []string{
		`class="badge badge-confirmed"`, "GF option",
		`class="badge badge-conflict"`, `class="badge badge-contains"`, "Contains gluten",
		"Shared fryer.", "Wheat dough.",
		`<details class="meal-group contains-group">`, // folded away by default
	} {
		if !strings.Contains(body, want) {
			t.Errorf("restaurant page missing %q", want)
		}
	}
}

func TestRestaurantPageWithoutSafetyProfile(t *testing.T) {
	body := get(t, newTestServer(t, seededStore()), "/restaurants/11").Body.String()
	for _, want := range []string{
		"No safety information yet", `popovertarget="call-ahead"`,
		"We don't have a phone number for this restaurant yet.",
		// Tweak 2: an empty review area is always explained.
		"No reviews yet. Be the first to review a meal here.",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("restaurant page missing %q", want)
		}
	}
	if strings.Contains(body, "Last verified") {
		t.Error("no profile means no verified date")
	}
}

// Tweak 2: "No reviews yet" on the list, the restaurant page and meals.
func TestNoReviewsYetLabels(t *testing.T) {
	h := newTestServer(t, seededStore())

	// The fake list rows have zero reviews.
	if list := get(t, h, "/restaurants?city=provo-ut").Body.String(); !strings.Contains(list, "No reviews yet") {
		t.Error("restaurant list should label restaurants with zero reviews")
	}

	// Restaurant 10 has reviews: no empty-state message there.
	if rest := get(t, h, "/restaurants/10").Body.String(); strings.Contains(rest, "Be the first to review a meal here.") {
		t.Error("restaurant with reviews shouldn't say it has none")
	}

	// Meal with zero reviews: card and meal page both explain it.
	if rest := get(t, h, "/restaurants/10").Body.String(); !strings.Contains(rest, `<span class="rating rating-none">No reviews yet</span>`) {
		t.Error("meal card with zero reviews should say No reviews yet")
	}
	if meal := get(t, h, "/meals/100").Body.String(); !strings.Contains(meal, "No reviews yet. Be the first to share how it went.") {
		t.Error("meal page with zero reviews should explain the empty list")
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
