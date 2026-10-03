package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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
	FoodType    *foodType    // the active food-type filter; nil = all
	Chips       []filterChip // the food-type filter chips
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

	ft := parseFoodType(r.URL.Query().Get("type"))
	typeValue := ""
	if ft != nil {
		typeValue = ft.Value
	}
	list, err := h.store.ListRestaurantsByCity(r.Context(), city.ID, typeValue)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	data := homePage{
		Selected:    &city,
		Restaurants: list,
		FoodType:    ft,
		Chips:       filterChips("/restaurants", url.Values{"city": {city.Slug}}, ft, "", "#results"),
	}

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

// restaurantPage is the data for the restaurant detail page.
type restaurantPage struct {
	Restaurant models.RestaurantDetail
	Safety     models.SafetyProfile
	// Meals grouped by gluten status, so GF options are listed first and
	// meals with gluten never get mixed in with them.
	GFMeals       []mealCard
	CrossMeals    []mealCard // no gluten ingredients, but shared equipment
	ContainsMeals []mealCard
	Photos        []photoGroup
	FoodType      *foodType    // the active food-type filter; nil = all
	Chips         []filterChip // the food-type filter chips for the meal list
}

// HasMeals reports whether any meal is left after filtering.
func (p restaurantPage) HasMeals() bool {
	return len(p.GFMeals)+len(p.CrossMeals)+len(p.ContainsMeals) > 0
}

// photoGroup is one labeled group in the restaurant's photo section.
// Photos may be empty; the template then shows a "No photo yet" tile.
type photoGroup struct {
	Label  string
	Photos []models.RestaurantPhoto
}

// First returns the group's first photo, or nil if it has none (the
// template shows a fallback tile for nil).
func (g photoGroup) First() *models.RestaurantPhoto {
	if len(g.Photos) == 0 {
		return nil
	}
	return &g.Photos[0]
}

// groupPhotos always returns the three groups, in the order testers asked
// for them, so the layout looks the same even when a group has no photo.
func groupPhotos(photos []models.RestaurantPhoto) []photoGroup {
	groups := []photoGroup{{Label: "Food"}, {Label: "Storefront"}, {Label: "Prep area"}}
	for _, p := range photos {
		switch p.Kind {
		case models.PhotoFood:
			groups[0].Photos = append(groups[0].Photos, p)
		case models.PhotoStorefront:
			groups[1].Photos = append(groups[1].Photos, p)
		case models.PhotoPrep:
			groups[2].Photos = append(groups[2].Photos, p)
		}
	}
	return groups
}

// groupMeals splits meal cards by gluten status, keeping their order.
// An unknown or empty status counts as a GF option (the database default).
func groupMeals(cards []mealCard) (gf, cross, contains []mealCard) {
	for _, c := range cards {
		switch c.Meal.GFStatus {
		case models.CrossContact:
			cross = append(cross, c)
		case models.ContainsGF:
			contains = append(contains, c)
		default:
			gf = append(gf, c)
		}
	}
	return gf, cross, contains
}

// NeedsCallAhead reports whether anything in the safety profile is
// conflicting or unverified (or there's no profile at all). The design system
// says to point diners to the call-ahead script whenever evidence is thin.
func (p restaurantPage) NeedsCallAhead() bool {
	if !p.Safety.Found || p.Safety.ConflictNote != "" {
		return true
	}
	for _, pr := range p.Safety.Precautions {
		if pr.Status != models.EvidenceConfirmed {
			return true
		}
	}
	return false
}

// ToConfirm lists the precautions to ask about when calling ahead: every
// badge that isn't confirmed.
func (p restaurantPage) ToConfirm() []models.Precaution {
	var out []models.Precaution
	for _, pr := range p.Safety.Precautions {
		if pr.Status != models.EvidenceConfirmed {
			out = append(out, pr)
		}
	}
	return out
}

// restaurant shows one restaurant and its GF meals: GET /restaurants/{id}.
func (h *Handler) restaurant(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		// "/restaurants/abc" isn't a real page, so treat it as not found.
		h.notFound(w, r)
		return
	}

	rest, err := h.store.GetRestaurant(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		h.notFound(w, r)
		return
	}
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	safety, err := h.store.GetSafetyProfile(r.Context(), id)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	meals, err := h.store.ListMealsByRestaurant(r.Context(), id)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	cards, err := h.mealCards(r.Context(), meals, visitorID(r))
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	photos, err := h.store.ListRestaurantPhotos(r.Context(), id)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	ft := parseFoodType(r.URL.Query().Get("type"))
	page := restaurantPage{
		Restaurant: rest,
		Safety:     safety,
		Photos:     groupPhotos(photos),
		FoodType:   ft,
		// "#meals" keeps the meal list in view after tapping a chip.
		Chips: filterChips(fmt.Sprintf("/restaurants/%d", id), nil, ft, "#meals", ""),
	}
	page.GFMeals, page.CrossMeals, page.ContainsMeals = groupMeals(filterMeals(cards, ft))
	h.render(w, r, http.StatusOK, "restaurant", page)
}

// parseID turns a URL path segment into a positive database id.
// Anything else (letters, negatives, decimals, too big) is rejected.
func parseID(s string) (int64, bool) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

// isFragmentRequest reports whether HTMX asked for just a page fragment.
// HTMX sets HX-Request on every request it makes. The exception is history
// restore (user pressed Back and HTMX's cache missed): then HTMX needs the
// full page, and it signals that with HX-History-Restore-Request.
func isFragmentRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true" &&
		r.Header.Get("HX-History-Restore-Request") != "true"
}
