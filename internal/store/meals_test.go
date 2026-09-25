package store

import (
	"context"
	"errors"
	"testing"
)

func TestGetRestaurant(t *testing.T) {
	s := newTestStore(t)
	provo := insertCity(t, s, "Provo", "provo-ut")
	id := insertRestaurant(t, s, provo, "Canyon Hearth Pizza")

	r, err := s.GetRestaurant(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != id || r.Name != "Canyon Hearth Pizza" || r.CityName != "Provo" || r.CitySlug != "provo-ut" {
		t.Errorf("got %+v", r)
	}

	if _, err := s.GetRestaurant(context.Background(), id+999); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown id: err = %v, want ErrNotFound", err)
	}
}

func TestListMealsByRestaurant(t *testing.T) {
	s := newTestStore(t)
	provo := insertCity(t, s, "Provo", "provo-ut")
	pizza := insertRestaurant(t, s, provo, "Canyon Hearth Pizza")
	other := insertRestaurant(t, s, provo, "Other Place")
	empty := insertRestaurant(t, s, provo, "No Meals Yet")

	insertMeal(t, s, pizza, "Margherita", "/static/images/meals/pizza.svg")
	insertMeal(t, s, pizza, "Harvest Salad", "")
	insertMeal(t, s, other, "Not Mine", "")

	t.Run("sorted by name, only this restaurant", func(t *testing.T) {
		meals, err := s.ListMealsByRestaurant(context.Background(), pizza)
		if err != nil {
			t.Fatal(err)
		}
		if len(meals) != 2 {
			t.Fatalf("got %d meals, want 2: %+v", len(meals), meals)
		}
		if meals[0].Name != "Harvest Salad" || meals[1].Name != "Margherita" {
			t.Errorf("order = %s, %s; want Harvest Salad, Margherita", meals[0].Name, meals[1].Name)
		}
		if meals[0].PhotoURL != "" || meals[1].PhotoURL != "/static/images/meals/pizza.svg" {
			t.Errorf("photo urls not scanned correctly: %+v", meals)
		}
		if meals[1].RestaurantID != pizza || meals[1].Description == "" {
			t.Errorf("meal fields not filled in: %+v", meals[1])
		}
	})

	t.Run("restaurant with no meals", func(t *testing.T) {
		meals, err := s.ListMealsByRestaurant(context.Background(), empty)
		if err != nil {
			t.Fatal(err)
		}
		if len(meals) != 0 {
			t.Errorf("got %d meals, want 0", len(meals))
		}
	})
}
