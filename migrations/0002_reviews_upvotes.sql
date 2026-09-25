-- +goose Up
CREATE TABLE reviews (
    id          BIGSERIAL PRIMARY KEY,
    meal_id     BIGINT NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    author_name TEXT NOT NULL,
    body        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- The app validates these too; the database is a last line of defense.
    CONSTRAINT reviews_author_name_length CHECK (char_length(author_name) BETWEEN 1 AND 40),
    CONSTRAINT reviews_body_length        CHECK (char_length(body) BETWEEN 1 AND 1000)
);

-- Serves "reviews for meal X, newest first" straight from the index.
CREATE INDEX reviews_meal_id_created_at_idx ON reviews (meal_id, created_at DESC);

-- One row per (meal, visitor). The composite primary key is what enforces
-- "one upvote per meal per visitor": a second insert would conflict.
CREATE TABLE meal_upvotes (
    meal_id    BIGINT NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    visitor_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (meal_id, visitor_id)
);

-- +goose Down
DROP TABLE meal_upvotes;
DROP TABLE reviews;
