-- +goose Up

-- Backfill any existing NULL employee names before enforcing NOT NULL,
-- so this migration doesn't fail on existing data.
UPDATE employee SET name = 'Unnamed Employee ' || id WHERE name IS NULL;

ALTER TABLE employee ALTER COLUMN name SET NOT NULL;

ALTER TABLE restaurant ADD CONSTRAINT restaurant_name_unique UNIQUE (name);

-- +goose Down

ALTER TABLE restaurant DROP CONSTRAINT IF EXISTS restaurant_name_unique;

ALTER TABLE employee ALTER COLUMN name DROP NOT NULL;