-- +goose Up

-- Food type, so diners can filter: testers wanted to separate salads from
-- grain-based options, and said good desserts are the hardest to find.
-- Everything starts as an entree; seed data sets the rest.
ALTER TABLE meals
    ADD COLUMN food_type TEXT NOT NULL DEFAULT 'entree',
    ADD CONSTRAINT meals_food_type_valid CHECK (food_type IN ('salad', 'grain', 'entree', 'dessert'));

-- +goose Down
ALTER TABLE meals
    DROP CONSTRAINT meals_food_type_valid,
    DROP COLUMN food_type;
