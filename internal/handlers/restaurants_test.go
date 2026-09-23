package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// getHX performs a GET as HTMX would, with the HX-Request header set.
func getHX(t *testing.T, h http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRestaurantsFullPage(t *testing.T) {
	rec := get(t, newTestServer(t, seededStore()), "/restaurants?city=provo-ut")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		// Full page, with layout and the dropdown pre-selected.
		"<html", `<option value="provo-ut" selected>`, disclaimer,
		// List heading and card fields.
		"GF-friendly restaurants in Provo",
		"Canyon Hearth Pizza", "Pizza", "410 Timber Lane", "3 GF meals",
		// Card links to the detail page; one meal uses the singular.
		`href="/restaurants/10"`, "1 GF meal<",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("full page missing %q", want)
		}
	}
	if got := rec.Header().Get("Vary"); got != "HX-Request" {
		t.Errorf("Vary = %q, want HX-Request", got)
	}
}

func TestRestaurantsHTMXPartial(t *testing.T) {
	rec := getHX(t, newTestServer(t, seededStore()), "/restaurants?city=provo-ut")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<html") || strings.Contains(body, disclaimer) {
		t.Error("HTMX request should get only the list fragment, not the full layout")
	}
	for _, want := range []string{
		"<title>Restaurants in Provo", // HTMX uses this to update the tab title
		"GF-friendly restaurants in Provo", "Canyon Hearth Pizza", "3 GF meals",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("partial missing %q", want)
		}
	}
}

func TestRestaurantsHistoryRestoreGetsFullPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/restaurants?city=provo-ut", nil)
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-History-Restore-Request", "true")
	rec := httptest.NewRecorder()
	newTestServer(t, seededStore()).ServeHTTP(rec, req)

	if !strings.Contains(rec.Body.String(), "<html") {
		t.Error("history restore should get the full page")
	}
}

func TestRestaurantsEmptyCity(t *testing.T) {
	h := newTestServer(t, seededStore())
	for name, rec := range map[string]*httptest.ResponseRecorder{
		"full page": get(t, h, "/restaurants?city=emptyville-ut"),
		"partial":   getHX(t, h, "/restaurants?city=emptyville-ut"),
	} {
		if rec.Code != http.StatusOK {
			t.Errorf("%s: status = %d, want 200", name, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "No restaurants here yet.") {
			t.Errorf("%s: missing empty-state message", name)
		}
	}
}

func TestRestaurantsUnknownCity(t *testing.T) {
	rec := get(t, newTestServer(t, seededStore()), "/restaurants?city=atlantis")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Page not found") {
		t.Error("unknown city should show the friendly 404 page")
	}
}

func TestRestaurantsMissingCityRedirectsHome(t *testing.T) {
	rec := get(t, newTestServer(t, seededStore()), "/restaurants")
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/" {
		t.Errorf("got %d to %q, want 303 to /", rec.Code, rec.Header().Get("Location"))
	}
}

func TestRestaurantsStoreErrorShowsFriendly500(t *testing.T) {
	st := seededStore()
	st.err = errors.New("pq: relation does not exist")
	rec := get(t, newTestServer(t, st), "/restaurants?city=provo-ut")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Something went wrong") || strings.Contains(body, "relation does not exist") {
		t.Error("500 page should be friendly and never show the raw error")
	}
}

func TestRestaurantsEscapesContent(t *testing.T) {
	st := seededStore()
	st.restaurants[2][0].Name = "<script>alert(1)</script>"
	body := get(t, newTestServer(t, st), "/restaurants?city=provo-ut").Body.String()
	if strings.Contains(body, "<script>alert(1)") {
		t.Error("restaurant name was not HTML-escaped")
	}
}
