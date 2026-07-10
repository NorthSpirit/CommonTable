-- name: GetRestaurant :one
SELECT * FROM restaurant
WHERE id = $1;

-- name: ListRestaurants :many
SELECT * FROM restaurant
ORDER BY name;

-- name: CreateRestaurant :one
INSERT INTO restaurant (name, address, phone, email, password_hash, logo_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateRestaurant :exec
UPDATE restaurant
SET name = $2, address = $3, phone = $4, email = $5, logo_url = $6
WHERE id = $1;

-- name: DeleteRestaurant :exec
DELETE FROM restaurant
WHERE id = $1;

-- name: GetRestaurantByName :one
SELECT * FROM restaurant
WHERE name = $1;

-- name: UpdateRestaurantPasswordHash :exec
UPDATE restaurant
SET password_hash = $2
WHERE id = $1;