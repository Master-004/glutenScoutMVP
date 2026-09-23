-- +goose Up
CREATE TABLE cities (
    id    BIGSERIAL PRIMARY KEY,
    name  TEXT NOT NULL,
    state TEXT NOT NULL,
    slug  TEXT NOT NULL UNIQUE -- used in URLs, e.g. denver-co
);

CREATE TABLE restaurants (
    id          BIGSERIAL PRIMARY KEY,
    city_id     BIGINT NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    address     TEXT NOT NULL,
    cuisine     TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE meals (
    id            BIGSERIAL PRIMARY KEY,
    restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    name          TEXT NOT NULL,
    description   TEXT NOT NULL DEFAULT '',
    photo_url     TEXT NOT NULL DEFAULT '', -- empty means "show the placeholder"
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Foreign keys aren't indexed automatically in Postgres; these make
-- "restaurants in city X" and "meals at restaurant Y" lookups fast.
CREATE INDEX restaurants_city_id_idx ON restaurants (city_id);
CREATE INDEX meals_restaurant_id_idx ON meals (restaurant_id);

-- +goose Down
DROP TABLE meals;
DROP TABLE restaurants;
DROP TABLE cities;
