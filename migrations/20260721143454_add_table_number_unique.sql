-- +goose Up

ALTER TABLE tables ADD CONSTRAINT tables_number_unique UNIQUE (number);

-- +goose Down

ALTER TABLE tables DROP CONSTRAINT IF EXISTS tables_number_unique;