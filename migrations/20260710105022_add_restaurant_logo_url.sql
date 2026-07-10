-- +goose Up

ALTER TABLE restaurant ADD COLUMN logo_url text;

-- +goose Down

ALTER TABLE restaurant DROP COLUMN IF EXISTS logo_url;