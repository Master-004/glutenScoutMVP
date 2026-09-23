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
	err         error                                // if set, data methods fail with it
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
// Emptyville (none).
func seededStore() *fakeStore {
	return &fakeStore{
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
	if strings.Contains(body, "restaurant-card") {
		t.Error("home page should start with an empty results area")
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
