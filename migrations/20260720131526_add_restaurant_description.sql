-- +goose Up

ALTER TABLE restaurant ADD COLUMN description text;
ALTER TABLE restaurant ADD COLUMN description_alt text;

-- +goose Down

ALTER TABLE restaurant DROP COLUMN IF EXISTS description_alt;
ALTER TABLE restaurant DROP COLUMN IF EXISTS description;