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

	"github.com/Master-004/glutenscout/web"
)

// fakeStore stands in for the real database so handler tests run without
// Postgres. Set fields to control what each method returns.
type fakeStore struct {
	pingErr error
}

func (f *fakeStore) Ping(ctx context.Context) error { return f.pingErr }

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
	rec := get(t, newTestServer(t, &fakeStore{}), "/")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"GlutenScout — coming soon", disclaimer, "/static/css/app.css", "/static/js/htmx.min.js"} {
		if !strings.Contains(body, want) {
			t.Errorf("home page missing %q", want)
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
