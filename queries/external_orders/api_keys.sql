-- name: CreateApiKey :one
INSERT INTO restaurant_api_keys (restaurant_id, key_hash, key_prefix, label)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetActiveApiKeyByHash :one
SELECT * FROM restaurant_api_keys
WHERE key_hash = $1 AND revoked_at IS NULL;

-- name: ListApiKeysByRestaurant :many
SELECT * FROM restaurant_api_keys
WHERE restaurant_id = $1
ORDER BY created_at DESC;

-- name: RevokeApiKey :execrows
UPDATE restaurant_api_keys
SET revoked_at = now()
WHERE id = $1 AND restaurant_id = $2 AND revoked_at IS NULL;
