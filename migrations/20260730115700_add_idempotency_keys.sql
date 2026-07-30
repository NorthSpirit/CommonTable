-- +goose Up

CREATE TABLE idempotency_keys (
                                  key varchar(64) PRIMARY KEY,
                                  created_at timestamp NOT NULL DEFAULT now()
);

-- +goose Down

DROP TABLE IF EXISTS idempotency_keys;