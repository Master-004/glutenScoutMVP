-- +goose Up

-- Phone number for the call-ahead script. Empty means "we don't have one".
ALTER TABLE restaurants ADD COLUMN phone TEXT NOT NULL DEFAULT '';

-- One safety profile per restaurant: the dated summary shown at the top of
-- its page. For now these rows come from seed data; later they can be
-- computed from restaurant answers and community reports.
CREATE TABLE safety_profiles (
    restaurant_id      BIGINT PRIMARY KEY REFERENCES restaurants(id) ON DELETE CASCADE,
    verified_on        DATE NOT NULL,          -- "Last verified {date}"
    answers_updated_on DATE,                   -- when the restaurant last answered our questions (NULL = never)
    conflict_note      TEXT NOT NULL DEFAULT '' -- reports disagreeing, in words; '' = no conflict
);

-- Precaution badges: one kitchen practice and how sure we are about it.
--   confirmed  = the restaurant and recent reviewers agree
--   conflict   = reports disagree, or a claim nobody has backed up
--   unverified = restaurant-only or out of date
CREATE TABLE precautions (
    id            BIGSERIAL PRIMARY KEY,
    restaurant_id BIGINT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    label         TEXT NOT NULL,   -- names a practice, never a verdict: "Separate prep area"
    status        TEXT NOT NULL,
    position      INT NOT NULL DEFAULT 0, -- display order on the page
    CONSTRAINT precautions_status_valid CHECK (status IN ('confirmed', 'conflict', 'unverified')),
    CONSTRAINT precautions_label_unique UNIQUE (restaurant_id, label)
);

-- Who wrote a review, and how they rated the food (stars are for food, never
-- gluten safety). Reviews posted before this migration have neither, so the
-- defaults ('' and 0) mean "not given"; the app requires both for new reviews.
ALTER TABLE reviews
    ADD COLUMN reviewer_type TEXT NOT NULL DEFAULT '',
    ADD COLUMN rating SMALLINT NOT NULL DEFAULT 0,
    ADD CONSTRAINT reviews_reviewer_type_valid CHECK (reviewer_type IN ('', 'celiac', 'strict', 'sensitive', 'restaurant')),
    ADD CONSTRAINT reviews_rating_range CHECK (rating BETWEEN 0 AND 5);

-- +goose Down
ALTER TABLE reviews
    DROP CONSTRAINT reviews_rating_range,
    DROP CONSTRAINT reviews_reviewer_type_valid,
    DROP COLUMN rating,
    DROP COLUMN reviewer_type;
DROP TABLE precautions;
DROP TABLE safety_profiles;
ALTER TABLE restaurants DROP COLUMN phone;
