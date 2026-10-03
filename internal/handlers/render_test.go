package handlers

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/Master-004/glutenscout/web"
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

// Every meal image the app can point to must exist in the embedded files,
// in whichever folder useIllustratedMeals selects.
func TestMealImagesExist(t *testing.T) {
	originals, err := fs.Glob(web.FS, "static/images/meals/*.svg")
	if err != nil || len(originals) == 0 {
		t.Fatalf("no original meal images found: %v", err)
	}
	for _, orig := range originals {
		url := mealImage("/" + orig)
		if _, err := fs.Stat(web.FS, strings.TrimPrefix(url, "/")); err != nil {
			t.Errorf("%s maps to %s, which doesn't exist", orig, url)
		}
	}
	// Restaurant photo illustrations (storefronts, prep areas) exist too.
	for _, p := range []string{"storefront-1", "storefront-2", "prep-1", "prep-2"} {
		if _, err := fs.Stat(web.FS, "static/images/places/"+p+".svg"); err != nil {
			t.Errorf("missing place image %s: %v", p, err)
		}
	}
	if got := mealImage("/static/images/places/prep-1.svg"); got != "/static/images/places/prep-1.svg" {
		t.Errorf("place images should pass through mealImage unchanged, got %q", got)
	}
	if got := mealImage("https://example.com/photo.jpg"); got != "https://example.com/photo.jpg" {
		t.Errorf("other URLs should pass through unchanged, got %q", got)
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
