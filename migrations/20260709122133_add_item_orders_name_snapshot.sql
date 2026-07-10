-- +goose Up

ALTER TABLE item_orders ADD COLUMN item_name_snapshot varchar(255);
ALTER TABLE item_orders ADD COLUMN item_description_snapshot text;

-- +goose Down

ALTER TABLE item_orders DROP COLUMN IF EXISTS item_description_snapshot;
ALTER TABLE item_orders DROP COLUMN IF EXISTS item_name_snapshot;