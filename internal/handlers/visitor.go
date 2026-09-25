package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

// There are no user accounts in the MVP. Instead, each browser gets a random
// anonymous "visitor ID" in a cookie. It's what lets us allow one upvote per
// meal per visitor and rate-limit reviews. A new browser (or incognito
// window, or cleared cookies) counts as a new visitor. That's acceptable for
// an MVP, and accounts can replace it later.

const (
	visitorCookieName = "gs_visitor"
	visitorIDBytes    = 16 // 128 random bits: practically impossible to guess
)

type visitorKey struct{} // private type so no other package can clash with our context key

// withVisitor makes sure every request has a visitor ID, creating one (and
// its cookie) on the first visit.
func withVisitor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := ""
		if c, err := r.Cookie(visitorCookieName); err == nil && isValidVisitorID(c.Value) {
			id = c.Value
		} else {
			id = newVisitorID()
			http.SetCookie(w, &http.Cookie{
				Name:     visitorCookieName,
				Value:    id,
				Path:     "/",
				MaxAge:   int((365 * 24 * time.Hour).Seconds()),
				HttpOnly: true,                 // page scripts can't read it
				SameSite: http.SameSiteLaxMode, // not sent on cross-site POSTs
				Secure:   r.TLS != nil,         // HTTPS-only when served over HTTPS
			})
		}
		ctx := context.WithValue(r.Context(), visitorKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// visitorID returns the current request's visitor ID (set by withVisitor).
func visitorID(r *http.Request) string {
	id, _ := r.Context().Value(visitorKey{}).(string)
	return id
}

func newVisitorID() string {
	b := make([]byte, visitorIDBytes)
	_, _ = rand.Read(b) // crypto/rand.Read never returns an error on supported platforms
	return hex.EncodeToString(b)
}

// isValidVisitorID accepts only IDs shaped like ours (32 lowercase hex chars),
// so a tampered cookie can't smuggle odd values into the database.
func isValidVisitorID(s string) bool {
	if len(s) != visitorIDBytes*2 {
		return false
	}
	for _, c := range s {
		if !('0' <= c && c <= '9' || 'a' <= c && c <= 'f') {
			return false
		}
	}
	return true
}
