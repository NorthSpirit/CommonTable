-- +goose Up

ALTER TABLE food_items ADD COLUMN deleted_at timestamp;

-- +goose Down

ALTER TABLE food_items DROP COLUMN IF EXISTS deleted_at;