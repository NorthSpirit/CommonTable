-- name: GetFoodMenu :one
SELECT * FROM food_menu
WHERE id = $1;

-- name: GetFoodMenuForDisplay :one
SELECT
    id,
    restaurant_id,
    CASE WHEN $2 = 'alt' THEN COALESCE(name_alt, name) ELSE name END AS name,
    CASE WHEN $2 = 'alt' THEN COALESCE(description_alt, description) ELSE description END AS description,
    image_url
FROM food_menu
WHERE id = $1;

-- name: ListFoodMenusByRestaurant :many
SELECT * FROM food_menu
WHERE restaurant_id = $1
ORDER BY name;

-- name: ListFoodMenusByRestaurantForDisplay :many
SELECT
    id,
    restaurant_id,
    CASE WHEN $2 = 'alt' THEN COALESCE(name_alt, name) ELSE name END AS name,
    CASE WHEN $2 = 'alt' THEN COALESCE(description_alt, description) ELSE description END AS description,
    image_url
FROM food_menu
WHERE restaurant_id = $1
ORDER BY name;

-- name: CreateFoodMenu :one
INSERT INTO food_menu (restaurant_id, name, name_alt, description, description_alt, image_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateFoodMenu :exec
UPDATE food_menu
SET name = $2, name_alt = $3, description = $4, description_alt = $5, image_url = $6
WHERE id = $1;

-- name: DeleteFoodMenu :exec
DELETE FROM food_menu
WHERE id = $1;