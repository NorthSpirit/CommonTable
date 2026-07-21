-- name: GetRestaurant :one
SELECT * FROM restaurant
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetRestaurantForDisplay :one
SELECT
    id,
    name,
    address,
    phone,
    email,
    logo_url,
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(description_alt, description) ELSE description END::text AS description
FROM restaurant
WHERE id = sqlc.arg(id)::integer AND deleted_at IS NULL;

-- name: ListRestaurants :many
SELECT * FROM restaurant
WHERE deleted_at IS NULL
ORDER BY name;

-- name: ListRestaurantsForDisplay :many
SELECT
    id,
    name,
    address,
    phone,
    email,
    logo_url,
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(description_alt, description) ELSE description END::text AS description
FROM restaurant
WHERE deleted_at IS NULL
ORDER BY name;

-- name: ListRestaurantsForAdmin :many
SELECT * FROM restaurant
ORDER BY deleted_at IS NULL DESC, name;

-- name: CreateRestaurant :one
INSERT INTO restaurant (name, address, phone, email, password_hash, logo_url, description, description_alt)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateRestaurant :exec
UPDATE restaurant
SET name = $2, address = $3, phone = $4, email = $5, logo_url = $6, description = $7, description_alt = $8
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteRestaurant :exec
UPDATE restaurant
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreRestaurant :exec
UPDATE restaurant
SET deleted_at = NULL
WHERE id = $1;

-- name: GetRestaurantByName :one
SELECT * FROM restaurant
WHERE name = $1 AND deleted_at IS NULL;

-- name: UpdateRestaurantPasswordHash :exec
UPDATE restaurant
SET password_hash = $2
WHERE id = $1 AND deleted_at IS NULL;