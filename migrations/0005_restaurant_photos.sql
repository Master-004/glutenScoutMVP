-- +goose Up

-- Photos on a restaurant page, in three groups testers asked for:
--   food       = what the dishes look like
--   storefront = the outside of the restaurant
--   prep       = the kitchen / prep area (some diners judge care by how clean it looks)
-- For now these are local placeholder illustrations from seed data; there
-- are no real uploads yet.
CREATE TABLE restaurant_photos (
    id            BIGSERIAL PRIMARY KEY,
    restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    kind          TEXT NOT NULL,
    url           TEXT NOT NULL,           -- a path under /static/
    alt           TEXT NOT NULL,           -- alt text for screen readers
    position      INT NOT NULL DEFAULT 0,  -- display order within its group
    CONSTRAINT restaurant_photos_kind_valid CHECK (kind IN ('food', 'storefront', 'prep')),
    CONSTRAINT restaurant_photos_unique UNIQUE (restaurant_id, kind, url)
);

CREATE INDEX restaurant_photos_restaurant_id_idx ON restaurant_photos (restaurant_id);

-- +goose Down
DROP TABLE restaurant_photos;
