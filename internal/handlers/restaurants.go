package handlers

import (
	"errors"
	"net/http"

	"github.com/Master-004/glutenscout/internal/models"
	"github.com/Master-004/glutenscout/internal/store"
)

// homePage is the data for the home page template. When Selected is nil,
// no city has been picked yet and the results area stays empty.
type homePage struct {
	Cities      []models.City
	Selected    *models.City
	Restaurants []models.RestaurantSummary
	Fragment    bool // true when rendering only the list for an HTMX swap
}

// home shows the welcome line and the city dropdown.
func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	cities, err := h.store.ListCities(r.Context())
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	h.render(w, r, http.StatusOK, "home", homePage{Cities: cities})
}

// restaurants lists a city's restaurants: GET /restaurants?city={slug}.
//
// The same URL serves two kinds of responses:
//   - An HTMX request (the dropdown changed) gets only the list fragment,
//     which HTMX swaps into #results without reloading the page.
//   - A normal visit (shared link, bookmark, refresh) gets the full page with
//     the city pre-selected, so the URL works on its own.
func (h *Handler) restaurants(w http.ResponseWriter, r *http.Request) {
	// Tell browser caches the response depends on this header; otherwise a
	// cached fragment could be shown for a full-page visit (or vice versa).
	w.Header().Add("Vary", "HX-Request")

	slug := r.URL.Query().Get("city")
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	city, err := h.store.GetCityBySlug(r.Context(), slug)
	if errors.Is(err, store.ErrNotFound) {
		h.notFound(w, r)
		return
	}
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	list, err := h.store.ListRestaurantsByCity(r.Context(), city.ID)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	data := homePage{Selected: &city, Restaurants: list}

	if isFragmentRequest(r) {
		data.Fragment = true
		h.renderPartial(w, r, http.StatusOK, "home", "restaurant_list", data)
		return
	}

	// Full page: also needs the dropdown's city list.
	if data.Cities, err = h.store.ListCities(r.Context()); err != nil {
		h.serverError(w, r, err)
		return
	}
	h.render(w, r, http.StatusOK, "home", data)
}

// isFragmentRequest reports whether HTMX asked for just a page fragment.
// HTMX sets HX-Request on every request it makes. The exception is history
// restore (user pressed Back and HTMX's cache missed): then HTMX needs the
// full page, and it signals that with HX-History-Restore-Request.
func isFragmentRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true" &&
		r.Header.Get("HX-History-Restore-Request") != "true"
}
