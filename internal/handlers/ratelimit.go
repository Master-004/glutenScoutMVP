package handlers

import (
	"sync"
	"time"
)

// rateLimiter allows each key (a visitor ID) at most `limit` actions per
// sliding `window`. It lives in memory, so it resets when the server
// restarts and isn't shared between multiple servers. That's fine for the
// MVP; a shared store (Postgres or Redis) would be needed at scale.
type rateLimiter struct {
	mu     sync.Mutex // handlers run concurrently, so guard the map
	limit  int
	window time.Duration
	hits   map[string][]time.Time
	now    func() time.Time // swappable in tests
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{limit: limit, window: window, hits: make(map[string][]time.Time), now: time.Now}
}

// Allow records an action for key and reports whether it's within the limit.
// Denied attempts aren't recorded, so they don't extend the wait.
func (l *rateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)

	// Keep only hits inside the window. Old ones are dropped here, whenever
	// the key is used again.
	recent := l.hits[key][:0]
	for _, t := range l.hits[key] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= l.limit {
		l.hits[key] = recent
		return false
	}
	l.hits[key] = append(recent, now)
	return true
}
