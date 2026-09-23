// Package handlers turns HTTP requests into pages. Handlers stay thin:
// read the request, ask the store for data, render a template.
package handlers

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Master-004/glutenscout/internal/models"
)

// Store is everything handlers need from the database. It's an interface
// (rather than *store.Store) so tests can pass in a fake with canned data
// and run without Postgres. Later slices add methods here.
type Store interface {
	Ping(ctx context.Context) error

	ListCities(ctx context.Context) ([]models.City, error)
	GetCityBySlug(ctx context.Context, slug string) (models.City, error)
	ListRestaurantsByCity(ctx context.Context, cityID int64) ([]models.RestaurantSummary, error)
}

// Handler holds shared dependencies for all HTTP handlers.
type Handler struct {
	store  Store
	tmpl   *Templates
	static fs.FS
	log    *slog.Logger
}

// New creates a Handler. static is the file system served under /static/.
func New(store Store, tmpl *Templates, static fs.FS, log *slog.Logger) *Handler {
	return &Handler{store: store, tmpl: tmpl, static: static, log: log}
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

	// Catch-all: anything not matched above gets the friendly 404 page.
	mux.HandleFunc("/", h.notFound)

	return h.recoverPanics(mux)
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
