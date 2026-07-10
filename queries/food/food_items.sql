-- name: GetFoodItem :one
SELECT id, restaurant_id, name, description, is_available, image_url, prep_time_mins, instructions
FROM food_items
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListFoodItemsByRestaurant :many
SELECT id, restaurant_id, name, description, is_available, image_url, prep_time_mins, instructions
FROM food_items
WHERE restaurant_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: ListAvailableFoodItemsByRestaurant :many
SELECT id, restaurant_id, name, description, is_available, image_url, prep_time_mins, instructions
FROM food_items
WHERE restaurant_id = $1 AND is_available = true AND deleted_at IS NULL
ORDER BY name;

-- name: CreateFoodItem :one
INSERT INTO food_items (restaurant_id, name, description, is_available, image_url, prep_time_mins, instructions)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, restaurant_id, name, description, is_available, image_url, prep_time_mins, instructions;

-- name: UpdateFoodItem :exec
UPDATE food_items
SET name = $2, description = $3, is_available = $4, image_url = $5, prep_time_mins = $6, instructions = $7
WHERE id = $1 AND deleted_at IS NULL;

-- name: SetFoodItemAvailability :exec
UPDATE food_items
SET is_available = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteFoodItem :exec
UPDATE food_items
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreFoodItem :exec
UPDATE food_items
SET deleted_at = NULL
WHERE id = $1;