-- +goose Up

-- Not every meal on a menu is gluten-free. Each meal now says how it stands:
--   gf            = made without gluten (a gluten-free option)
--   cross_contact = no gluten ingredients, but prepared with shared equipment
--                   (e.g. fries from a fryer shared with breaded chicken)
--   contains      = contains gluten
-- Existing meals were all listed as GF options, so they default to 'gf'.
ALTER TABLE meals
    ADD COLUMN gf_status TEXT NOT NULL DEFAULT 'gf',
    ADD COLUMN gf_note   TEXT NOT NULL DEFAULT '', -- why, in words: "Fried in the shared fryer."
    ADD CONSTRAINT meals_gf_status_valid CHECK (gf_status IN ('gf', 'cross_contact', 'contains'));

-- +goose Down
ALTER TABLE meals
    DROP CONSTRAINT meals_gf_status_valid,
    DROP COLUMN gf_note,
    DROP COLUMN gf_status;
