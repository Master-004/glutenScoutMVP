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
			got, err := NormalizeReview(review(tt.author, tt.body))
			if got.Author != tt.wantAuthor || got.Body != tt.wantBody {
				t.Errorf("got (%q, %q), want (%q, %q)", got.Author, got.Body, tt.wantAuthor, tt.wantBody)
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

// review builds a valid ReviewInput (a celiac reviewer, 5 stars) so tests
// about the name and body don't have to repeat the other fields.
func review(author, body string) ReviewInput {
	return ReviewInput{Author: author, Body: body, ReviewerType: "celiac", Rating: 5}
}

func TestNormalizeReviewTypeAndRating(t *testing.T) {
	tests := []struct {
		name          string
		reviewerType  string
		rating        int
		wantRating    int
		wantTypeErr   bool
		wantRatingErr bool
	}{
		{"diner with rating", "strict", 4, 4, false, false},
		{"missing type", "", 4, 4, true, false},
		{"unknown type", "chef", 4, 4, true, false},
		{"diner without rating", "celiac", 0, 0, false, true},
		{"rating too high", "sensitive", 6, 6, false, true},
		{"restaurant rating is dropped", "restaurant", 5, 0, false, false},
		{"restaurant needs no rating", "restaurant", 0, 0, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeReview(ReviewInput{Body: "ok", ReviewerType: tt.reviewerType, Rating: tt.rating})
			if got.Rating != tt.wantRating {
				t.Errorf("rating = %d, want %d", got.Rating, tt.wantRating)
			}
			var e *ReviewInputError
			errors.As(err, &e)
			if e == nil {
				e = &ReviewInputError{}
			}
			if (e.ReviewerType != "") != tt.wantTypeErr || (e.Rating != "") != tt.wantRatingErr {
				t.Errorf("errors = %+v, want type err %v, rating err %v", e, tt.wantTypeErr, tt.wantRatingErr)
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
		r, err := s.CreateReview(ctx, meal, review("  ", "  Safe for me!  "))
		if err != nil {
			t.Fatal(err)
		}
		if r.ID == 0 || r.MealID != meal || r.AuthorName != "Anonymous" || r.Body != "Safe for me!" ||
			r.ReviewerType != "celiac" || r.Rating != 5 || r.CreatedAt.IsZero() {
			t.Errorf("got %+v", r)
		}
	})

	t.Run("invalid input is rejected and not saved", func(t *testing.T) {
		before, _ := s.ListReviewsByMeal(ctx, meal)
		for _, body := range []string{"", "   ", strings.Repeat("x", 1001)} {
			_, err := s.CreateReview(ctx, meal, review("Sam", body))
			var inputErr *ReviewInputError
			if !errors.As(err, &inputErr) {
				t.Errorf("body of %d chars: err = %v, want *ReviewInputError", len(body), err)
			}
		}
		if _, err := s.CreateReview(ctx, meal, review(strings.Repeat("n", 41), "ok")); err == nil {
			t.Error("41-char name should be rejected")
		}
		after, _ := s.ListReviewsByMeal(ctx, meal)
		if len(after) != len(before) {
			t.Errorf("invalid reviews were saved: %d -> %d", len(before), len(after))
		}
	})

	t.Run("unknown meal", func(t *testing.T) {
		if _, err := s.CreateReview(ctx, meal+999, review("Sam", "Hello")); !errors.Is(err, ErrNotFound) {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}

func TestListReviewsByMealNewestFirst(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	meal := seedMeal(t, s)
	for _, body := range []string{"first", "second", "third"} {
		if _, err := s.CreateReview(ctx, meal, review("Sam", body)); err != nil {
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
	if _, err := s.CreateReview(ctx, a, review("", "yum")); err != nil {
		t.Fatal(err)
	}
	// A 4-star review, plus a restaurant reply that has no rating and must
	// not drag the average down.
	if _, err := s.CreateReview(ctx, a, ReviewInput{Body: "good", ReviewerType: "strict", Rating: 4}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreateReview(ctx, a, ReviewInput{Body: "thanks!", ReviewerType: "restaurant"}); err != nil {
		t.Fatal(err)
	}

	stats, err := s.MealStats(ctx, []int64{a, b}, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if st := stats[a]; st.Upvotes != 2 || st.Reviews != 3 || !st.Upvoted || st.AvgRating != 4.5 || st.RatingCount != 2 {
		t.Errorf("meal A = %+v, want 2 upvotes, 3 reviews, upvoted by bob, 4.5 avg from 2 ratings", st)
	}
	if st := stats[b]; st.Upvotes != 0 || st.Reviews != 0 || st.Upvoted {
		t.Errorf("meal B = %+v, want all zero", st)
	}
}
