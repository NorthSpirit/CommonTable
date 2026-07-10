-- name: GetRestaurant :one
SELECT * FROM restaurant
WHERE id = $1;

-- name: ListRestaurants :many
SELECT * FROM restaurant
ORDER BY name;

-- name: CreateRestaurant :one
INSERT INTO restaurant (name, address, phone, email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateRestaurant :exec
UPDATE restaurant
SET name = $2, address = $3, phone = $4, email = $5
WHERE id = $1;

-- name: DeleteRestaurant :exec
DELETE FROM restaurant
WHERE id = $1;