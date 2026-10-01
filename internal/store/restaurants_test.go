package store

import (
	"context"
	"errors"
	"testing"
)

func TestListCitiesSortedByName(t *testing.T) {
	s := newTestStore(t)
	insertCity(t, s, "Provo", "provo-ut")
	insertCity(t, s, "Denver", "denver-co")

	cities, err := s.ListCities(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(cities) != 2 || cities[0].Name != "Denver" || cities[1].Name != "Provo" {
		t.Fatalf("got %+v, want Denver then Provo", cities)
	}
}

func TestGetCityBySlug(t *testing.T) {
	s := newTestStore(t)
	id := insertCity(t, s, "Provo", "provo-ut")

	city, err := s.GetCityBySlug(context.Background(), "provo-ut")
	if err != nil {
		t.Fatal(err)
	}
	if city.ID != id || city.Name != "Provo" {
		t.Errorf("got %+v", city)
	}

	if _, err := s.GetCityBySlug(context.Background(), "nowhere"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown slug: err = %v, want ErrNotFound", err)
	}
}

func TestGetSafetyProfile(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	city := insertCity(t, s, "Provo", "provo-ut")
	withProfile := insertRestaurant(t, s, city, "Checked Kitchen")
	without := insertRestaurant(t, s, city, "Unknown Kitchen")

	if _, err := s.pool.Exec(ctx, `
		INSERT INTO safety_profiles (restaurant_id, verified_on, conflict_note) VALUES ($1, '2026-09-12', 'Reports disagree.');
	`, withProfile); err != nil {
		t.Fatal(err)
	}
	// Inserted out of order to check that position decides display order.
	if _, err := s.pool.Exec(ctx, `
		INSERT INTO precautions (restaurant_id, label, status, position) VALUES
		  ($1, 'Shared fryer?', 'conflict', 2),
		  ($1, 'Gloves changed', 'confirmed', 1)`, withProfile); err != nil {
		t.Fatal(err)
	}

	p, err := s.GetSafetyProfile(ctx, withProfile)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Found || p.VerifiedOn.Format("2006-01-02") != "2026-09-12" || p.AnswersUpdatedOn != nil || p.ConflictNote != "Reports disagree." {
		t.Errorf("profile = %+v", p)
	}
	if len(p.Precautions) != 2 || p.Precautions[0].Label != "Gloves changed" || p.Precautions[1].Status != "conflict" {
		t.Errorf("precautions = %+v, want Gloves changed then Shared fryer?", p.Precautions)
	}

	// No profile is a normal state, not an error.
	if p, err := s.GetSafetyProfile(ctx, without); err != nil || p.Found {
		t.Errorf("no profile: got %+v, %v; want zero profile, nil", p, err)
	}

	// The list shows the verified date too.
	list, err := s.ListRestaurantsByCity(ctx, city)
	if err != nil {
		t.Fatal(err)
	}
	if list[0].VerifiedOn == nil || list[1].VerifiedOn != nil {
		t.Errorf("verified dates = %v, %v; want set, nil", list[0].VerifiedOn, list[1].VerifiedOn)
	}

	// The database rejects an unknown evidence state.
	if _, err := s.pool.Exec(ctx, `INSERT INTO precautions (restaurant_id, label, status) VALUES ($1, 'X', 'safe')`, withProfile); err == nil {
		t.Error(`status "safe" should be rejected by the CHECK constraint`)
	}
}

func TestListRestaurantsByCity(t *testing.T) {
	s := newTestStore(t)
	provo := insertCity(t, s, "Provo", "provo-ut")
	denver := insertCity(t, s, "Denver", "denver-co")
	empty := insertCity(t, s, "Emptyville", "emptyville-ut")

	zeta := insertRestaurant(t, s, provo, "Zeta Grill")
	insertMeal(t, s, zeta, "Burger", "")
	insertMeal(t, s, zeta, "Fries", "")
	// Meals that aren't GF options must not count toward "N GF meals".
	mustInsert(t, s, `INSERT INTO meals (restaurant_id, name, gf_status) VALUES ($1, 'Onion Rings', 'contains') RETURNING id`, zeta)
	mustInsert(t, s, `INSERT INTO meals (restaurant_id, name, gf_status) VALUES ($1, 'Shared Fries', 'cross_contact') RETURNING id`, zeta)
	insertRestaurant(t, s, provo, "Alpha Cafe") // no meals
	insertRestaurant(t, s, denver, "Denver Diner")

	t.Run("sorted by name with meal counts, only this city", func(t *testing.T) {
		got, err := s.ListRestaurantsByCity(context.Background(), provo)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Fatalf("got %d restaurants, want 2: %+v", len(got), got)
		}
		if got[0].Name != "Alpha Cafe" || got[0].MealCount != 0 {
			t.Errorf("first = %s (%d meals), want Alpha Cafe (0)", got[0].Name, got[0].MealCount)
		}
		if got[1].Name != "Zeta Grill" || got[1].MealCount != 2 {
			t.Errorf("second = %s (%d meals), want Zeta Grill (2)", got[1].Name, got[1].MealCount)
		}
		if got[1].CityID != provo || got[1].Address == "" {
			t.Errorf("restaurant fields not filled in: %+v", got[1])
		}
	})

	t.Run("ratings and verified date; reviews don't inflate meal counts", func(t *testing.T) {
		ctx := context.Background()
		burger := insertMeal(t, s, insertRestaurant(t, s, denver, "Rated Place"), "Burger", "")
		for _, rating := range []int{5, 4, 4} {
			if _, err := s.CreateReview(ctx, burger, ReviewInput{Body: "ok", ReviewerType: "celiac", Rating: rating}); err != nil {
				t.Fatal(err)
			}
		}
		got, err := s.ListRestaurantsByCity(ctx, denver)
		if err != nil {
			t.Fatal(err)
		}
		var rated, diner = got[1], got[0] // sorted: Denver Diner, Rated Place
		if rated.MealCount != 1 || rated.RatingCount != 3 || rated.AvgRating < 4.33 || rated.AvgRating > 4.34 {
			t.Errorf("Rated Place = %+v, want 1 meal, 3 ratings, avg 4.33", rated)
		}
		if diner.RatingCount != 0 || diner.AvgRating != 0 || diner.VerifiedOn != nil {
			t.Errorf("Denver Diner = %+v, want no ratings and no verified date", diner)
		}
	})

	t.Run("city with no restaurants", func(t *testing.T) {
		got, err := s.ListRestaurantsByCity(context.Background(), empty)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Errorf("got %d restaurants, want 0", len(got))
		}
	})
}
