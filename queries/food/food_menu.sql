-- name: GetFoodMenu :one
SELECT * FROM food_menu
WHERE id = $1;

-- name: ListFoodMenusByRestaurant :many
SELECT * FROM food_menu
WHERE restaurant_id = $1
ORDER BY name;

-- name: CreateFoodMenu :one
INSERT INTO food_menu (restaurant_id, name, description, image_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateFoodMenu :exec
UPDATE food_menu
SET name = $2, description = $3, image_url = $4
WHERE id = $1;

-- name: DeleteFoodMenu :exec
DELETE FROM food_menu
WHERE id = $1;