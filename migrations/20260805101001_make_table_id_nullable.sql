-- +goose Up

ALTER TABLE food_orders ALTER COLUMN table_id DROP NOT NULL;

-- +goose Down

ALTER TABLE food_orders ALTER COLUMN table_id SET NOT NULL;