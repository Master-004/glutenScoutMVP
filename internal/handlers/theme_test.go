package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestThemeSwitch(t *testing.T) {
	c := newClient(t, newTestServer(t, seededStore()))

	// No choice yet: follow the device, so no data-theme attribute.
	page := c.get("/").Body.String()
	if strings.Contains(page, "data-theme=") || !strings.Contains(page, `value="system" aria-pressed="true"`) {
		t.Error("default should be Auto with no data-theme attribute")
	}

	// Pick dark: cookie saved, sent back to the page we came from.
	req := themeRequest("dark", "http://example.com/restaurants?city=provo-ut")
	rec := c.do(req)
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/restaurants?city=provo-ut" {
		t.Errorf("got %d to %q, want 303 back to the referring page", rec.Code, rec.Header().Get("Location"))
	}
	page = c.get("/").Body.String()
	if !strings.Contains(page, `<html lang="en" data-theme="dark">`) || !strings.Contains(page, `value="dark" aria-pressed="true"`) {
		t.Error("after picking dark, pages should render data-theme=\"dark\"")
	}

	// Back to Auto: cookie cleared.
	c.do(themeRequest("system", "http://example.com/"))
	if page = c.get("/").Body.String(); strings.Contains(page, "data-theme=") {
		t.Error("Auto should remove data-theme")
	}
}

func TestThemeIgnoresBadInput(t *testing.T) {
	c := newClient(t, newTestServer(t, seededStore()))
	c.do(themeRequest(`"><script>`, ""))
	if page := c.get("/").Body.String(); strings.Contains(page, "data-theme=") {
		t.Error("unknown theme values should fall back to Auto")
	}
}

func TestThemeRedirectStaysOnSite(t *testing.T) {
	h := newTestServer(t, seededStore())
	for _, ref := range []string{"https://evil.example/phish", "", "not a url\x7f"} {
		rec := newClient(t, h).do(themeRequest("light", ref))
		if loc := rec.Header().Get("Location"); loc != "/" {
			t.Errorf("Referer %q: redirected to %q, want /", ref, loc)
		}
	}
}

// themeRequest builds the POST the header switch sends. httptest requests
// have Host example.com, so that's the "same site" Referer host.
func themeRequest(theme, referer string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "http://example.com/theme", strings.NewReader(url.Values{"theme": {theme}}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
	return req
}
