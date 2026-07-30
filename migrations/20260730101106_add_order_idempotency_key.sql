-- +goose Up

ALTER TABLE food_orders ADD COLUMN idempotency_key varchar(64);
CREATE UNIQUE INDEX food_orders_idempotency_key_unique ON food_orders (idempotency_key) WHERE idempotency_key IS NOT NULL;

-- +goose Down

DROP INDEX IF EXISTS food_orders_idempotency_key_unique;
ALTER TABLE food_orders DROP COLUMN IF EXISTS idempotency_key;