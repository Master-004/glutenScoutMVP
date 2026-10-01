package handlers

import (
	"strings"
	"testing"
)

func TestStars(t *testing.T) {
	tests := []struct {
		rating any
		want   string
	}{
		{0.0, "empty empty empty empty empty"},
		{4, "full full full full empty"},      // int ratings from single reviews
		{4.3, "full full full full half"},     // rounds to the nearest half
		{4.2, "full full full full empty"},    // 4.2 -> 4.0
		{2.5, "full full half empty empty"},
		{5.0, "full full full full full"},
	}
	for _, tt := range tests {
		if got := strings.Join(stars(tt.rating), " "); got != tt.want {
			t.Errorf("stars(%v) = %q, want %q", tt.rating, got, tt.want)
		}
	}
	if got := ratingText(4.33); got != "4.5" {
		t.Errorf("ratingText(4.33) = %q, want 4.5", got)
	}
}

func TestAsPhrase(t *testing.T) {
	for in, want := range map[string]string{
		"Shared oven?":         "shared oven",
		"GF hoisin on request": "GF hoisin on request",
		"Dedicated toaster":    "dedicated toaster",
	} {
		if got := asPhrase(in); got != want {
			t.Errorf("asPhrase(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestInitial(t *testing.T) {
	tests := []struct{ name, want string }{
		{"Sam", "S"},
		{"anonymous", "A"},
		{"  jo", "J"},
		{"Émile", "É"}, // multi-byte first letter stays whole
		{"", "?"},
	}
	for _, tt := range tests {
		if got := initial(tt.name); got != tt.want {
			t.Errorf("initial(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
