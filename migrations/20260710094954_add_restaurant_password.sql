-- +goose Up

ALTER TABLE restaurant ADD COLUMN password_hash varchar(255);

-- +goose Down

ALTER TABLE restaurant DROP COLUMN IF EXISTS password_hash;