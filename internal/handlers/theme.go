package handlers

import (
	"net/http"
	"net/url"
	"time"
)

// Light/dark theme switch, done on the server so it needs no custom
// JavaScript. The header has a tiny form that POSTs the choice here; we save
// it in a cookie and send the visitor back to the page they were on. Every
// page then renders <html data-theme="..."> and the CSS picks colors from it.
//
// "system" (no cookie) means "follow the phone/computer setting".

const themeCookieName = "gs_theme"

// themeFromRequest returns "light", "dark", or "system".
func themeFromRequest(r *http.Request) string {
	if c, err := r.Cookie(themeCookieName); err == nil && (c.Value == "light" || c.Value == "dark") {
		return c.Value
	}
	return "system"
}

func (h *Handler) setTheme(w http.ResponseWriter, r *http.Request) {
	switch theme := r.PostFormValue("theme"); theme {
	case "light", "dark":
		http.SetCookie(w, &http.Cookie{
			Name:     themeCookieName,
			Value:    theme,
			Path:     "/",
			MaxAge:   int((365 * 24 * time.Hour).Seconds()),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   r.TLS != nil,
		})
	default: // "system" or anything unexpected: forget the choice
		http.SetCookie(w, &http.Cookie{Name: themeCookieName, Value: "", Path: "/", MaxAge: -1})
	}
	// 303 See Other: the browser follows up with a GET, so a refresh won't re-POST.
	http.Redirect(w, r, backTo(r), http.StatusSeeOther)
}

// backTo picks where to send the visitor after switching themes: the page
// they came from (the Referer header), but only if it's on this site.
// Checking the host stops the redirect being abused to bounce people to
// another website. We use Referer (not the URL the page was first rendered
// at) because HTMX may have updated the address bar since, e.g. after picking
// a city.
func backTo(r *http.Request) string {
	ref, err := url.Parse(r.Referer())
	if err != nil || ref.Host != r.Host || ref.Path == "" || ref.Path[0] != '/' {
		return "/"
	}
	back := ref.Path
	if ref.RawQuery != "" {
		back += "?" + ref.RawQuery
	}
	return back
}
