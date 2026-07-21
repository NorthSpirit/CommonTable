-- +goose Up

ALTER TABLE restaurant ADD COLUMN deleted_at timestamp;

-- +goose Down

ALTER TABLE restaurant DROP COLUMN IF EXISTS deleted_at;