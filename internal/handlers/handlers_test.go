package handlers

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Master-004/glutenscout/internal/models"
	"github.com/Master-004/glutenscout/internal/store"
	"github.com/Master-004/glutenscout/web"
)

// fakeStore stands in for the real database so handler tests run without
// Postgres. Set fields to control what each method returns.
type fakeStore struct {
	pingErr     error
	cities      []models.City
	restaurants map[int64][]models.RestaurantSummary // keyed by city id
	details     map[int64]models.RestaurantDetail    // keyed by restaurant id
	safety      map[int64]models.SafetyProfile       // keyed by restaurant id
	meals       map[int64][]models.Meal              // keyed by restaurant id
	reviews     map[int64][]models.Review            // keyed by meal id, newest first
	upvotes     map[int64]map[string]bool            // meal id -> visitor ids who upvoted
	err         error                                // if set, data methods fail with it
}

func (f *fakeStore) GetMeal(ctx context.Context, id int64) (models.MealDetail, error) {
	if f.err != nil {
		return models.MealDetail{}, f.err
	}
	for restID, meals := range f.meals {
		for _, m := range meals {
			if m.ID == id {
				d := f.details[restID]
				return models.MealDetail{Meal: m, RestaurantName: d.Name, CityName: d.CityName, CitySlug: d.CitySlug}, nil
			}
		}
	}
	return models.MealDetail{}, store.ErrNotFound
}

func (f *fakeStore) ListReviewsByMeal(ctx context.Context, mealID int64) ([]models.Review, error) {
	return f.reviews[mealID], f.err
}

// CreateReview mimics the real store: it validates, then saves newest-first.
func (f *fakeStore) CreateReview(ctx context.Context, mealID int64, in store.ReviewInput) (models.Review, error) {
	if f.err != nil {
		return models.Review{}, f.err
	}
	in, err := store.NormalizeReview(in)
	if err != nil {
		return models.Review{}, err
	}
	if f.reviews == nil {
		f.reviews = map[int64][]models.Review{}
	}
	r := models.Review{
		ID: int64(len(f.reviews[mealID]) + 1), MealID: mealID, AuthorName: in.Author, Body: in.Body,
		ReviewerType: in.ReviewerType, Rating: in.Rating, CreatedAt: time.Now(),
	}
	f.reviews[mealID] = append([]models.Review{r}, f.reviews[mealID]...)
	return r, nil
}

func (f *fakeStore) ToggleUpvote(ctx context.Context, mealID int64, visitorID string) (int, bool, error) {
	if f.err != nil {
		return 0, false, f.err
	}
	if _, err := f.GetMeal(ctx, mealID); err != nil {
		return 0, false, err
	}
	if f.upvotes == nil {
		f.upvotes = map[int64]map[string]bool{}
	}
	if f.upvotes[mealID] == nil {
		f.upvotes[mealID] = map[string]bool{}
	}
	voters := f.upvotes[mealID]
	if voters[visitorID] {
		delete(voters, visitorID)
	} else {
		voters[visitorID] = true
	}
	return len(voters), voters[visitorID], nil
}

func (f *fakeStore) MealStats(ctx context.Context, mealIDs []int64, visitorID string) (map[int64]models.MealStats, error) {
	if f.err != nil {
		return nil, f.err
	}
	stats := map[int64]models.MealStats{}
	for _, id := range mealIDs {
		var rating models.Rating
		sum := 0
		for _, r := range f.reviews[id] {
			if r.Rating > 0 {
				sum += r.Rating
				rating.RatingCount++
			}
		}
		if rating.RatingCount > 0 {
			rating.AvgRating = float64(sum) / float64(rating.RatingCount)
		}
		stats[id] = models.MealStats{
			MealID:  id,
			Rating:  rating,
			Upvotes: len(f.upvotes[id]),
			Reviews: len(f.reviews[id]),
			Upvoted: f.upvotes[id][visitorID],
		}
	}
	return stats, nil
}

func (f *fakeStore) GetRestaurant(ctx context.Context, id int64) (models.RestaurantDetail, error) {
	if f.err != nil {
		return models.RestaurantDetail{}, f.err
	}
	r, ok := f.details[id]
	if !ok {
		return models.RestaurantDetail{}, store.ErrNotFound
	}
	return r, nil
}

func (f *fakeStore) GetSafetyProfile(ctx context.Context, restaurantID int64) (models.SafetyProfile, error) {
	return f.safety[restaurantID], f.err // missing key = zero profile, like the real store
}

func (f *fakeStore) ListMealsByRestaurant(ctx context.Context, restaurantID int64) ([]models.Meal, error) {
	return f.meals[restaurantID], f.err
}

func (f *fakeStore) Ping(ctx context.Context) error { return f.pingErr }

func (f *fakeStore) ListCities(ctx context.Context) ([]models.City, error) {
	return f.cities, f.err
}

func (f *fakeStore) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	if f.err != nil {
		return models.City{}, f.err
	}
	for _, c := range f.cities {
		if c.Slug == slug {
			return c, nil
		}
	}
	return models.City{}, store.ErrNotFound
}

func (f *fakeStore) ListRestaurantsByCity(ctx context.Context, cityID int64) ([]models.RestaurantSummary, error) {
	return f.restaurants[cityID], f.err
}

// seededStore returns a fake with two cities: Provo (two restaurants) and
// Emptyville (none). Restaurant 10 has two meals (one without a photo) and a
// safety profile with a conflict; restaurant 11 has no meals, no profile and
// no phone.
func seededStore() *fakeStore {
	pizza := models.Restaurant{ID: 10, CityID: 2, Name: "Canyon Hearth Pizza", Cuisine: "Pizza", Address: "410 Timber Lane", Description: "Wood-fired pizza.", Phone: "(801) 555-0101"}
	solo := models.Restaurant{ID: 11, CityID: 2, Name: "Solo Cafe", Cuisine: "Cafe", Address: "1 Main St"}
	answered := time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)
	return &fakeStore{
		details: map[int64]models.RestaurantDetail{
			10: {Restaurant: pizza, CityName: "Provo", CitySlug: "provo-ut", Rating: models.Rating{AvgRating: 4.3, RatingCount: 3}},
			11: {Restaurant: solo, CityName: "Provo", CitySlug: "provo-ut"},
		},
		safety: map[int64]models.SafetyProfile{
			10: {
				Found:            true,
				VerifiedOn:       time.Date(2026, 9, 12, 0, 0, 0, 0, time.UTC),
				AnswersUpdatedOn: &answered,
				ConflictNote:     "2 reports disagree about the oven.",
				Precautions: []models.Precaution{
					{Label: "Separate prep area", Status: models.EvidenceConfirmed},
					{Label: "Shared oven?", Status: models.EvidenceConflict},
					{Label: "GF hoisin on request", Status: models.EvidenceUnverified},
				},
			},
		},
		meals: map[int64][]models.Meal{
			10: {
				{ID: 100, RestaurantID: 10, Name: "GF Margherita", Description: "Tomato and basil.", PhotoURL: "/static/images/meals/pizza.svg"},
				{ID: 101, RestaurantID: 10, Name: "Harvest Salad", Description: "Greens and squash."},
			},
		},
		cities: []models.City{
			{ID: 1, Name: "Emptyville", State: "UT", Slug: "emptyville-ut"},
			{ID: 2, Name: "Provo", State: "UT", Slug: "provo-ut"},
		},
		restaurants: map[int64][]models.RestaurantSummary{
			2: {
				{Restaurant: models.Restaurant{ID: 10, CityID: 2, Name: "Canyon Hearth Pizza", Cuisine: "Pizza", Address: "410 Timber Lane"}, MealCount: 3},
				{Restaurant: models.Restaurant{ID: 11, CityID: 2, Name: "Solo Cafe", Cuisine: "Cafe", Address: "1 Main St"}, MealCount: 1},
			},
		},
	}
}

// newTestServer builds the real routes with the real templates and a fake store.
func newTestServer(t *testing.T, st Store) http.Handler {
	t.Helper()
	tmpl, err := LoadTemplates(web.FS)
	if err != nil {
		t.Fatalf("load templates: %v", err)
	}
	static, err := fs.Sub(web.FS, "static")
	if err != nil {
		t.Fatal(err)
	}
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	return New(st, tmpl, static, quiet).Routes()
}

// get performs a GET request against h and returns the recorded response.
func get(t *testing.T, h http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

const disclaimer = "GlutenScout info comes from the community and may be out of date."

func TestHealthz(t *testing.T) {
	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
		wantBody   string
	}{
		{"db ok", nil, http.StatusOK, "ok"},
		{"db down", errors.New("connection refused"), http.StatusServiceUnavailable, "database unavailable"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := get(t, newTestServer(t, &fakeStore{pingErr: tt.pingErr}), "/healthz")
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := rec.Body.String(); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
			// The raw DB error must never reach the client.
			if strings.Contains(rec.Body.String(), "connection refused") {
				t.Error("response leaked the internal error")
			}
		})
	}
}

func TestHomePage(t *testing.T) {
	rec := get(t, newTestServer(t, seededStore()), "/")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"Find gluten-free-friendly restaurants near you",
		`<option value="provo-ut"`, "Provo, UT",
		`hx-get="/restaurants"`, `hx-target="#results"`, `hx-push-url="true"`,
		disclaimer, "/static/css/app.css", "/static/js/htmx.min.js",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("home page missing %q", want)
		}
	}
	// No city picked yet, so no results should render.
	if strings.Contains(body, "restaurant-row") {
		t.Error("home page should start with an empty results area")
	}
}

// Tweak 1: "Find restaurants" is the first screen's one obvious, green
// action, and every screen points to its next step and a way back.
func TestNextStepsAreObvious(t *testing.T) {
	h := newTestServer(t, seededStore())

	home := get(t, h, "/").Body.String()
	if !strings.Contains(home, `<button class="button button-large" type="submit">Find restaurants</button>`) {
		t.Error("home page needs a visible primary Find restaurants button")
	}
	if strings.Count(home, `class="button`) != 1 {
		t.Error("home page should have exactly one primary button")
	}
	if !strings.Contains(home, `name="city" required`) {
		t.Error("submitting without a city should ask for one")
	}

	list := get(t, h, "/restaurants?city=provo-ut").Body.String()
	for _, want := range []string{"Tap a restaurant to see its gluten-free meals.", "See meals ›"} {
		if !strings.Contains(list, want) {
			t.Errorf("restaurant list missing %q", want)
		}
	}

	rest := get(t, h, "/restaurants/10").Body.String()
	if !strings.Contains(rest, "Tap a meal to see reviews") || !strings.Contains(rest, "Back to Provo restaurants") {
		t.Error("restaurant page needs a next step and a way back to the list")
	}

	meal := get(t, h, "/meals/100").Body.String()
	for _, want := range []string{
		`href="/restaurants?city=provo-ut">← Provo restaurants</a>`, // back to the list
		`href="/restaurants/10">Canyon Hearth Pizza</a>`,            // back to the restaurant
		">Write a review</button>",                                  // the next step
	} {
		if !strings.Contains(meal, want) {
			t.Errorf("meal page missing %q", want)
		}
	}
}

func TestUnknownURLShowsFriendly404(t *testing.T) {
	rec := get(t, newTestServer(t, &fakeStore{}), "/this/does/not/exist")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Page not found") || !strings.Contains(body, disclaimer) {
		t.Error("404 page should use the layout and show a friendly message")
	}
}

func TestStaticFiles(t *testing.T) {
	h := newTestServer(t, &fakeStore{})
	if rec := get(t, h, "/static/css/app.css"); rec.Code != http.StatusOK {
		t.Errorf("app.css status = %d, want 200", rec.Code)
	}
	if rec := get(t, h, "/static/images/"); rec.Code != http.StatusNotFound {
		t.Errorf("directory listing status = %d, want 404", rec.Code)
	}
}
