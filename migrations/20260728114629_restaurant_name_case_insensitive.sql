-- +goose Up

ALTER TABLE restaurant DROP CONSTRAINT IF EXISTS restaurant_name_unique;
CREATE UNIQUE INDEX restaurant_name_unique_ci ON restaurant (lower(name)) WHERE deleted_at IS NULL;

-- +goose Down

DROP INDEX IF EXISTS restaurant_name_unique_ci;
ALTER TABLE restaurant ADD CONSTRAINT restaurant_name_unique UNIQUE (name);