-- +goose Up

ALTER TABLE tables ADD COLUMN deleted_at timestamp;

ALTER TABLE tables DROP CONSTRAINT IF EXISTS tables_number_unique;
CREATE UNIQUE INDEX tables_number_unique_active ON tables (number) WHERE deleted_at IS NULL;

-- +goose Down

DROP INDEX IF EXISTS tables_number_unique_active;
ALTER TABLE tables ADD CONSTRAINT tables_number_unique UNIQUE (number);
ALTER TABLE tables DROP COLUMN IF EXISTS deleted_at;