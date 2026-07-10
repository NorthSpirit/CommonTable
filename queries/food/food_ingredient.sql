-- name: GetFoodIngredient :one
SELECT * FROM food_ingredient
WHERE id = $1;

-- name: ListFoodIngredients :many
SELECT * FROM food_ingredient
ORDER BY name;

-- name: CreateFoodIngredient :one
INSERT INTO food_ingredient (name, description)
VALUES ($1, $2)
    RETURNING *;

-- name: UpdateFoodIngredient :exec
UPDATE food_ingredient
SET name = $2, description = $3
WHERE id = $1;

-- name: DeleteFoodIngredient :exec
DELETE FROM food_ingredient
WHERE id = $1;