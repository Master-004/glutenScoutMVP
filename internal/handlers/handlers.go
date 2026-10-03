// Package handlers turns HTTP requests into pages. Handlers stay thin:
// read the request, ask the store for data, render a template.
package handlers

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Master-004/glutenscout/internal/models"
	"github.com/Master-004/glutenscout/internal/store"
)

// Store is everything handlers need from the database. It's an interface
// (rather than *store.Store) so tests can pass in a fake with canned data
// and run without Postgres.
type Store interface {
	Ping(ctx context.Context) error

	ListCities(ctx context.Context) ([]models.City, error)
	GetCityBySlug(ctx context.Context, slug string) (models.City, error)
	ListRestaurantsByCity(ctx context.Context, cityID int64) ([]models.RestaurantSummary, error)

	GetRestaurant(ctx context.Context, id int64) (models.RestaurantDetail, error)
	GetSafetyProfile(ctx context.Context, restaurantID int64) (models.SafetyProfile, error)
	ListRestaurantPhotos(ctx context.Context, restaurantID int64) ([]models.RestaurantPhoto, error)
	ListMealsByRestaurant(ctx context.Context, restaurantID int64) ([]models.Meal, error)

	GetMeal(ctx context.Context, id int64) (models.MealDetail, error)
	ListReviewsByMeal(ctx context.Context, mealID int64) ([]models.Review, error)
	CreateReview(ctx context.Context, mealID int64, in store.ReviewInput) (models.Review, error)
	ToggleUpvote(ctx context.Context, mealID int64, visitorID string) (count int, upvoted bool, err error)
	MealStats(ctx context.Context, mealIDs []int64, visitorID string) (map[int64]models.MealStats, error)
}

// Review rate limit: at most this many reviews per visitor per hour.
const (
	reviewLimit       = 5
	reviewLimitWindow = time.Hour
)

// Handler holds shared dependencies for all HTTP handlers.
type Handler struct {
	store         Store
	tmpl          *Templates
	static        fs.FS
	log           *slog.Logger
	reviewLimiter *rateLimiter
}

// New creates a Handler. static is the file system served under /static/.
func New(store Store, tmpl *Templates, static fs.FS, log *slog.Logger) *Handler {
	return &Handler{
		store:         store,
		tmpl:          tmpl,
		static:        static,
		log:           log,
		reviewLimiter: newRateLimiter(reviewLimit, reviewLimitWindow),
	}
}

// Routes registers every URL the app serves and returns the root handler.
// Go 1.22+ patterns include the HTTP method ("GET /path"), so a POST to a
// GET-only route automatically gets a 405 response.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", noDirListing(http.FileServerFS(h.static))))
	mux.HandleFunc("GET /healthz", h.health)
	// "{$}" means "exactly /". Without it, "/" would match every URL.
	mux.HandleFunc("GET /{$}", h.home)
	mux.HandleFunc("GET /restaurants", h.restaurants)
	mux.HandleFunc("GET /restaurants/{id}", h.restaurant)
	mux.HandleFunc("GET /meals/{id}", h.meal)
	mux.HandleFunc("POST /meals/{id}/upvote", h.upvote)
	mux.HandleFunc("POST /meals/{id}/reviews", h.createReview)
	mux.HandleFunc("POST /theme", h.setTheme)

	// Catch-all: anything not matched above gets the friendly 404 page.
	mux.HandleFunc("/", h.notFound)

	// Middleware wraps the router like layers of an onion; a request passes
	// through them from the outside in:
	//   recoverPanics -> cross-origin check -> visitor cookie -> routes
	return h.recoverPanics(h.crossOriginProtection(withVisitor(mux)))
}

// crossOriginProtection blocks POSTs that come from other websites (CSRF).
// Without it, a malicious page could make a visitor's browser submit reviews
// or upvotes to GlutenScout. Go's built-in check uses the browser's
// Sec-Fetch-Site and Origin headers. GET requests are always allowed.
func (h *Handler) crossOriginProtection(next http.Handler) http.Handler {
	cop := http.NewCrossOriginProtection()
	cop.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.log.Warn("blocked cross-origin request", "method", r.Method, "path", r.URL.Path, "origin", r.Header.Get("Origin"))
		http.Error(w, "This request came from another website and was blocked.", http.StatusForbidden)
	}))
	return cop.Handler(next)
}

// noDirListing hides directory listings (e.g. /static/images/) by returning 404.
func noDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") || r.URL.Path == "" {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// recoverPanics turns an unexpected panic in any handler into the friendly
// 500 page instead of a dropped connection.
func (h *Handler) recoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				h.log.Error("panic", "path", r.URL.Path, "panic", v)
				h.render(w, r, http.StatusInternalServerError, "500", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
