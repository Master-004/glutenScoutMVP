package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// testClient acts like one browser: it remembers cookies between requests,
// so it keeps the same anonymous visitor ID.
type testClient struct {
	t       *testing.T
	h       http.Handler
	cookies map[string]*http.Cookie
}

func newClient(t *testing.T, h http.Handler) *testClient {
	return &testClient{t: t, h: h, cookies: map[string]*http.Cookie{}}
}

func (c *testClient) do(req *http.Request) *httptest.ResponseRecorder {
	c.t.Helper()
	for _, ck := range c.cookies {
		req.AddCookie(ck)
	}
	rec := httptest.NewRecorder()
	c.h.ServeHTTP(rec, req)
	for _, ck := range rec.Result().Cookies() {
		c.cookies[ck.Name] = ck
	}
	return rec
}

func (c *testClient) get(path string) *httptest.ResponseRecorder {
	return c.do(httptest.NewRequest(http.MethodGet, path, nil))
}

// post submits a form. hx=true sends it the way HTMX would.
func (c *testClient) post(path string, form url.Values, hx bool) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if hx {
		req.Header.Set("HX-Request", "true")
	}
	return c.do(req)
}

// ---------- Visitor cookie ----------

func TestVisitorCookieIsSetOnce(t *testing.T) {
	c := newClient(t, newTestServer(t, seededStore()))

	first := c.get("/")
	ck := first.Result().Cookies()
	if len(ck) != 1 || ck[0].Name != visitorCookieName {
		t.Fatalf("first visit should set the visitor cookie, got %v", ck)
	}
	if !ck[0].HttpOnly || ck[0].SameSite != http.SameSiteLaxMode || !isValidVisitorID(ck[0].Value) {
		t.Errorf("cookie should be HttpOnly, SameSite=Lax, 32 hex chars: %+v", ck[0])
	}

	if again := c.get("/"); len(again.Result().Cookies()) != 0 {
		t.Error("returning visitor should keep their cookie, not get a new one")
	}
}

func TestVisitorCookieTamperedIsReplaced(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: visitorCookieName, Value: "'; DROP TABLE meals; --"})
	rec := httptest.NewRecorder()
	newTestServer(t, seededStore()).ServeHTTP(rec, req)

	ck := rec.Result().Cookies()
	if len(ck) != 1 || !isValidVisitorID(ck[0].Value) {
		t.Errorf("invalid cookie should be replaced with a fresh id, got %v", ck)
	}
}

// ---------- Meal page ----------

func TestMealPage(t *testing.T) {
	c := newClient(t, newTestServer(t, seededStore()))
	c.post("/meals/100/reviews", url.Values{"author_name": {"Jo"}, "body": {"Safe and delicious"}}, true)

	rec := c.get("/meals/100")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"<h1>GF Margherita</h1>", "Tomato and basil.", `src="/static/images/meals/pizza.svg"`,
		`href="/restaurants/10"`, "Canyon Hearth Pizza", // restaurant link
		`hx-post="/meals/100/upvote"`,                             // upvote button
		`id="reviews"`, "Reviews (1)", "Jo", "Safe and delicious", // review list
		`hx-post="/meals/100/reviews"`, `name="author_name"`, `name="body"`, // review form
		disclaimer,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("meal page missing %q", want)
		}
	}
}

func TestMealPageNotFound(t *testing.T) {
	h := newTestServer(t, seededStore())
	for _, path := range []string{"/meals/999", "/meals/abc", "/meals/0"} {
		if rec := get(t, h, path); rec.Code != http.StatusNotFound {
			t.Errorf("%s: status = %d, want 404", path, rec.Code)
		}
	}
}

// ---------- Upvotes ----------

func TestUpvoteToggle(t *testing.T) {
	h := newTestServer(t, seededStore())
	alice := newClient(t, h)
	alice.get("/") // first visit: get a visitor cookie

	rec := alice.post("/meals/100/upvote", nil, true)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<html") {
		t.Error("HTMX upvote should return only the button")
	}
	if !strings.Contains(body, `class="upvote-count">1<`) || !strings.Contains(body, "is-upvoted") || !strings.Contains(body, `aria-pressed="true"`) {
		t.Errorf("after upvote: want count 1, highlighted; got:\n%s", body)
	}

	// Tapping again removes it.
	body = alice.post("/meals/100/upvote", nil, true).Body.String()
	if !strings.Contains(body, `class="upvote-count">0<`) || strings.Contains(body, "is-upvoted") {
		t.Errorf("after second tap: want count 0, not highlighted; got:\n%s", body)
	}

	// Upvote again and refresh: the state survives a full page load.
	alice.post("/meals/100/upvote", nil, true)
	page := alice.get("/meals/100").Body.String()
	if !strings.Contains(page, `class="upvote-count">1<`) || !strings.Contains(page, "is-upvoted") {
		t.Error("refreshing the meal page should show the visitor's upvote")
	}
}

func TestUpvoteDifferentVisitorsCountSeparately(t *testing.T) {
	h := newTestServer(t, seededStore())
	alice, bob := newClient(t, h), newClient(t, h)
	alice.get("/")
	bob.get("/")

	alice.post("/meals/100/upvote", nil, true)
	body := bob.post("/meals/100/upvote", nil, true).Body.String()
	if !strings.Contains(body, `class="upvote-count">2<`) {
		t.Errorf("two visitors should make 2 upvotes; got:\n%s", body)
	}

	// Bob sees alice's upvote in the count, but his own button state is his.
	bob.post("/meals/100/upvote", nil, true) // bob removes his
	page := bob.get("/meals/100").Body.String()
	if !strings.Contains(page, `class="upvote-count">1<`) || strings.Contains(page, "is-upvoted") {
		t.Error("bob should see 1 upvote (alice's), not highlighted for him")
	}
}

func TestUpvoteWithoutHTMXRedirects(t *testing.T) {
	c := newClient(t, newTestServer(t, seededStore()))
	rec := c.post("/meals/100/upvote", nil, false)
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/meals/100" {
		t.Errorf("got %d to %q, want 303 to /meals/100", rec.Code, rec.Header().Get("Location"))
	}
}

func TestUpvoteUnknownMeal(t *testing.T) {
	c := newClient(t, newTestServer(t, seededStore()))
	for _, path := range []string{"/meals/999/upvote", "/meals/abc/upvote"} {
		if rec := c.post(path, nil, true); rec.Code != http.StatusNotFound {
			t.Errorf("%s: status = %d, want 404", path, rec.Code)
		}
	}
}

// ---------- Reviews ----------

func TestCreateReviewSuccess(t *testing.T) {
	st := seededStore()
	c := newClient(t, newTestServer(t, st))
	c.post("/meals/100/reviews", url.Values{"body": {"Older review"}}, true)

	rec := c.post("/meals/100/reviews", url.Values{"author_name": {"  Sam "}, "body": {"Newest review"}}, true)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<html") {
		t.Error("HTMX post should return only the reviews section")
	}
	if !strings.Contains(body, "Reviews (2)") || !strings.Contains(body, "Thanks! Your review was posted.") {
		t.Error("response should show the updated count and a success message")
	}
	// Newest first: the new review appears above the older one.
	if strings.Index(body, "Newest review") > strings.Index(body, "Older review") {
		t.Error("new review should be at the top of the list")
	}
	// The form is cleared.
	if !strings.Contains(body, `placeholder="Anonymous" value=""`) || strings.Contains(body, ">Newest review</textarea>") {
		t.Error("form should be empty after a successful post")
	}
	// Name trimmed; blank name defaulted.
	if st.reviews[100][0].AuthorName != "Sam" || st.reviews[100][1].AuthorName != "Anonymous" {
		t.Errorf("authors = %q, %q", st.reviews[100][0].AuthorName, st.reviews[100][1].AuthorName)
	}
}

func TestCreateReviewValidationErrorsKeepInput(t *testing.T) {
	tests := []struct {
		name, author, body, wantMsg, wantKept string
	}{
		{"empty review", "Sam", "   ", "Please write a review before submitting.", `value="Sam"`},
		{"too long", "Sam", strings.Repeat("x", 1001), "Reviews can be at most 1000 characters", strings.Repeat("x", 1001) + "</textarea>"},
		{"name too long", strings.Repeat("n", 41), "Nice", "Display name can be at most 40 characters", ">Nice</textarea>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := seededStore()
			c := newClient(t, newTestServer(t, st))
			rec := c.post("/meals/100/reviews", url.Values{"author_name": {tt.author}, "body": {tt.body}}, true)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Errorf("status = %d, want 422", rec.Code)
			}
			body := rec.Body.String()
			if !strings.Contains(body, tt.wantMsg) {
				t.Errorf("missing error message %q", tt.wantMsg)
			}
			if !strings.Contains(body, tt.wantKept) {
				t.Errorf("typed input not kept (want %q)", tt.wantKept[:min(len(tt.wantKept), 40)])
			}
			if len(st.reviews[100]) != 0 {
				t.Error("invalid review was saved")
			}
		})
	}
}

func TestCreateReviewWithoutHTMX(t *testing.T) {
	c := newClient(t, newTestServer(t, seededStore()))

	ok := c.post("/meals/100/reviews", url.Values{"body": {"Plain form post"}}, false)
	if ok.Code != http.StatusSeeOther || ok.Header().Get("Location") != "/meals/100#reviews" {
		t.Errorf("success: got %d to %q, want 303 to /meals/100#reviews", ok.Code, ok.Header().Get("Location"))
	}

	bad := c.post("/meals/100/reviews", url.Values{"body": {""}}, false)
	if bad.Code != http.StatusUnprocessableEntity || !strings.Contains(bad.Body.String(), "<html") {
		t.Errorf("error without HTMX should re-render the full page with 422, got %d", bad.Code)
	}
}

func TestCreateReviewEscapesHTML(t *testing.T) {
	c := newClient(t, newTestServer(t, seededStore()))
	attack := `<script>alert("xss")</script>`
	c.post("/meals/100/reviews", url.Values{"author_name": {attack}, "body": {attack}}, true)

	body := c.get("/meals/100").Body.String()
	if strings.Contains(body, attack) {
		t.Fatal("review HTML was not escaped")
	}
	if !strings.Contains(body, "&lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt;") {
		t.Error("escaped review text should still be shown as plain text")
	}
}

func TestCreateReviewRateLimit(t *testing.T) {
	st := seededStore()
	h := newTestServer(t, st)
	c := newClient(t, h)
	c.get("/")

	for i := 1; i <= reviewLimit; i++ {
		if rec := c.post("/meals/100/reviews", url.Values{"body": {"ok"}}, true); rec.Code != http.StatusOK {
			t.Fatalf("review %d: status = %d, want 200", i, rec.Code)
		}
	}
	rec := c.post("/meals/100/reviews", url.Values{"body": {"one too many"}}, true)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Please wait a bit") || !strings.Contains(rec.Body.String(), "one too many</textarea>") {
		t.Error("rate-limit response should explain and keep the typed text")
	}
	if n := len(st.reviews[100]); n != reviewLimit {
		t.Errorf("saved %d reviews, want %d", n, reviewLimit)
	}

	// A different visitor is unaffected.
	other := newClient(t, h)
	if rec := other.post("/meals/100/reviews", url.Values{"body": {"hi"}}, true); rec.Code != http.StatusOK {
		t.Errorf("other visitor: status = %d, want 200", rec.Code)
	}
}

func TestCreateReviewUnknownMeal(t *testing.T) {
	c := newClient(t, newTestServer(t, seededStore()))
	if rec := c.post("/meals/999/reviews", url.Values{"body": {"hi"}}, true); rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

// ---------- Cross-site request protection ----------

func TestCrossSitePostsAreBlocked(t *testing.T) {
	st := seededStore()
	h := newTestServer(t, st)

	for _, path := range []string{"/meals/100/upvote", "/meals/100/reviews"} {
		for name, setHeaders := range map[string]func(*http.Request){
			"Sec-Fetch-Site cross-site": func(r *http.Request) { r.Header.Set("Sec-Fetch-Site", "cross-site") },
			"Origin from another site":  func(r *http.Request) { r.Header.Set("Origin", "https://evil.example") },
		} {
			req := httptest.NewRequest(http.MethodPost, path, strings.NewReader("body=pwned"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			setHeaders(req)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != http.StatusForbidden {
				t.Errorf("%s with %s: status = %d, want 403", path, name, rec.Code)
			}
		}
	}
	if len(st.reviews[100]) != 0 || len(st.upvotes[100]) != 0 {
		t.Error("blocked requests must not change data")
	}

	// Same-origin browser requests still work.
	req := httptest.NewRequest(http.MethodPost, "/meals/100/upvote", nil)
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Errorf("same-origin POST: status = %d, want 303", rec.Code)
	}
}

// ---------- Restaurant page stats ----------

func TestRestaurantPageShowsMealStats(t *testing.T) {
	h := newTestServer(t, seededStore())
	c := newClient(t, h)
	c.get("/")
	c.post("/meals/100/upvote", nil, true)
	c.post("/meals/100/reviews", url.Values{"body": {"great"}}, true)

	body := c.get("/restaurants/10").Body.String()
	for _, want := range []string{
		`href="/meals/100"`,                   // card links to the meal page
		`class="upvote-count">1<`,             // upvote count on the card
		`hx-post="/meals/100/upvote"`,         // upvote button on the card
		`href="/meals/100#reviews">1 review<`, // singular review count
		`href="/meals/101#reviews">0 reviews<`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("restaurant page missing %q", want)
		}
	}
}
