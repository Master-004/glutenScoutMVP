package store

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// NormalizeReview is pure logic, so these tests run without a database.
func TestNormalizeReview(t *testing.T) {
	tests := []struct {
		name          string
		author, body  string
		wantAuthor    string
		wantBody      string
		wantAuthorErr bool
		wantBodyErr   bool
	}{
		{"valid", "Sam", "Great and safe!", "Sam", "Great and safe!", false, false},
		{"trims whitespace", "  Sam  ", "  Tasty  ", "Sam", "Tasty", false, false},
		{"blank name becomes Anonymous", "   ", "Good", "Anonymous", "Good", false, false},
		{"empty body", "Sam", "", "Sam", "", false, true},
		{"whitespace-only body", "Sam", " \n\t ", "Sam", "", false, true},
		{"body at limit", "", strings.Repeat("a", 1000), "Anonymous", strings.Repeat("a", 1000), false, false},
		{"body too long", "", strings.Repeat("a", 1001), "Anonymous", strings.Repeat("a", 1001), false, true},
		{"body counts characters not bytes", "", strings.Repeat("é", 1000), "Anonymous", strings.Repeat("é", 1000), false, false},
		{"name at limit", strings.Repeat("n", 40), "ok", strings.Repeat("n", 40), "ok", false, false},
		{"name too long", strings.Repeat("n", 41), "ok", strings.Repeat("n", 41), "ok", true, false},
		{"both invalid", strings.Repeat("n", 41), "", strings.Repeat("n", 41), "", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			author, body, err := NormalizeReview(tt.author, tt.body)
			if author != tt.wantAuthor || body != tt.wantBody {
				t.Errorf("got (%q, %q), want (%q, %q)", author, body, tt.wantAuthor, tt.wantBody)
			}
			var inputErr *ReviewInputError
			if !tt.wantAuthorErr && !tt.wantBodyErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !errors.As(err, &inputErr) {
				t.Fatalf("err = %v, want *ReviewInputError", err)
			}
			if (inputErr.Author != "") != tt.wantAuthorErr {
				t.Errorf("author error = %q, want error: %v", inputErr.Author, tt.wantAuthorErr)
			}
			if (inputErr.Body != "") != tt.wantBodyErr {
				t.Errorf("body error = %q, want error: %v", inputErr.Body, tt.wantBodyErr)
			}
		})
	}
}

// seedMeal creates a city, restaurant, and meal, returning the meal id.
func seedMeal(t *testing.T, s *Store) int64 {
	t.Helper()
	city := insertCity(t, s, "Provo", "provo-ut")
	rest := insertRestaurant(t, s, city, "Canyon Hearth Pizza")
	return insertMeal(t, s, rest, "GF Margherita", "")
}

func TestGetMeal(t *testing.T) {
	s := newTestStore(t)
	id := seedMeal(t, s)

	m, err := s.GetMeal(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "GF Margherita" || m.RestaurantName != "Canyon Hearth Pizza" {
		t.Errorf("got %+v", m)
	}
	if _, err := s.GetMeal(context.Background(), id+999); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown meal: err = %v, want ErrNotFound", err)
	}
}

func TestCreateReview(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	meal := seedMeal(t, s)

	t.Run("valid review is saved trimmed with default author", func(t *testing.T) {
		r, err := s.CreateReview(ctx, meal, "  ", "  Safe for me!  ")
		if err != nil {
			t.Fatal(err)
		}
		if r.ID == 0 || r.MealID != meal || r.AuthorName != "Anonymous" || r.Body != "Safe for me!" || r.CreatedAt.IsZero() {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("invalid input is rejected and not saved", func(t *testing.T) {
		before, _ := s.ListReviewsByMeal(ctx, meal)
		for _, body := range []string{"", "   ", strings.Repeat("x", 1001)} {
			_, err := s.CreateReview(ctx, meal, "Sam", body)
			var inputErr *ReviewInputError
			if !errors.As(err, &inputErr) {
				t.Errorf("body of %d chars: err = %v, want *ReviewInputError", len(body), err)
			}
		}
		if _, err := s.CreateReview(ctx, meal, strings.Repeat("n", 41), "ok"); err == nil {
			t.Error("41-char name should be rejected")
		}
		after, _ := s.ListReviewsByMeal(ctx, meal)
		if len(after) != len(before) {
			t.Errorf("invalid reviews were saved: %d -> %d", len(before), len(after))
		}
	})

	t.Run("unknown meal", func(t *testing.T) {
		if _, err := s.CreateReview(ctx, meal+999, "Sam", "Hello"); !errors.Is(err, ErrNotFound) {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}

func TestListReviewsByMealNewestFirst(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	meal := seedMeal(t, s)
	for _, body := range []string{"first", "second", "third"} {
		if _, err := s.CreateReview(ctx, meal, "Sam", body); err != nil {
			t.Fatal(err)
		}
	}

	reviews, err := s.ListReviewsByMeal(ctx, meal)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range reviews {
		got = append(got, r.Body)
	}
	if strings.Join(got, ",") != "third,second,first" {
		t.Errorf("order = %v, want newest first", got)
	}
}

func TestToggleUpvote(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	meal := seedMeal(t, s)

	step := func(visitor string, wantCount int, wantUpvoted bool) {
		t.Helper()
		count, upvoted, err := s.ToggleUpvote(ctx, meal, visitor)
		if err != nil {
			t.Fatal(err)
		}
		if count != wantCount || upvoted != wantUpvoted {
			t.Errorf("%s: got (count %d, upvoted %v), want (%d, %v)", visitor, count, upvoted, wantCount, wantUpvoted)
		}
	}

	step("alice", 1, true)  // add
	step("alice", 0, false) // tap again removes it
	step("alice", 1, true)  // and adds it back
	step("bob", 2, true)    // a different visitor counts separately
	step("bob", 1, false)   // bob removing doesn't touch alice's upvote

	// One upvote per visitor: alice has exactly one row despite 3 toggles.
	stats, err := s.MealStats(ctx, []int64{meal}, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if st := stats[meal]; st.Upvotes != 1 || !st.Upvoted {
		t.Errorf("stats = %+v, want 1 upvote by alice", st)
	}

	if _, _, err := s.ToggleUpvote(ctx, meal+999, "alice"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown meal: err = %v, want ErrNotFound", err)
	}
}

func TestMealStats(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	city := insertCity(t, s, "Provo", "provo-ut")
	rest := insertRestaurant(t, s, city, "Place")
	a := insertMeal(t, s, rest, "A", "")
	b := insertMeal(t, s, rest, "B", "")

	for _, visitor := range []string{"alice", "bob"} {
		if _, _, err := s.ToggleUpvote(ctx, a, visitor); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.CreateReview(ctx, a, "", "yum"); err != nil {
		t.Fatal(err)
	}

	stats, err := s.MealStats(ctx, []int64{a, b}, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if st := stats[a]; st.Upvotes != 2 || st.Reviews != 1 || !st.Upvoted {
		t.Errorf("meal A = %+v, want 2 upvotes, 1 review, upvoted by bob", st)
	}
	if st := stats[b]; st.Upvotes != 0 || st.Reviews != 0 || st.Upvoted {
		t.Errorf("meal B = %+v, want all zero", st)
	}
}
