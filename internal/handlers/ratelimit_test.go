package handlers

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	clock := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	l := newRateLimiter(5, time.Hour)
	l.now = func() time.Time { return clock }

	for i := 1; i <= 5; i++ {
		if !l.Allow("alice") {
			t.Fatalf("action %d should be allowed", i)
		}
		clock = clock.Add(time.Minute)
	}
	if l.Allow("alice") {
		t.Error("6th action within the hour should be denied")
	}
	if !l.Allow("bob") {
		t.Error("other visitors have their own limit")
	}

	// The first action was at 12:00, so at 13:00:01 one slot frees up.
	clock = time.Date(2026, 1, 1, 13, 0, 1, 0, time.UTC)
	if !l.Allow("alice") {
		t.Error("should be allowed again once the oldest action leaves the window")
	}
	if l.Allow("alice") {
		t.Error("only one slot should have freed up")
	}
}
