-- name: CheckIdempotencyKey :one
SELECT key FROM idempotency_keys WHERE key = $1;

-- name: RecordIdempotencyKey :exec
INSERT INTO idempotency_keys (key) VALUES ($1);