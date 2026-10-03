-- +goose Up

-- Testers all eat out with people who aren't gluten-free. This flag marks
-- restaurants with plenty for everyone, so groups can spot them in the list.
ALTER TABLE restaurants ADD COLUMN good_for_groups BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE restaurants DROP COLUMN good_for_groups;
