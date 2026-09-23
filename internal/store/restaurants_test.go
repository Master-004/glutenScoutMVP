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

func TestListRestaurantsByCity(t *testing.T) {
	s := newTestStore(t)
	provo := insertCity(t, s, "Provo", "provo-ut")
	denver := insertCity(t, s, "Denver", "denver-co")
	empty := insertCity(t, s, "Emptyville", "emptyville-ut")

	zeta := insertRestaurant(t, s, provo, "Zeta Grill")
	insertMeal(t, s, zeta, "Burger", "")
	insertMeal(t, s, zeta, "Fries", "")
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
