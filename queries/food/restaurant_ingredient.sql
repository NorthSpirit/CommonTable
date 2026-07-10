-- name: LinkRestaurantIngredient :exec
INSERT INTO restaurant_ingredient (restaurant_id, ingredient_id, is_shared)
VALUES ($1, $2, $3);

-- name: UpdateRestaurantIngredientShared :exec
UPDATE restaurant_ingredient
SET is_shared = $3
WHERE restaurant_id = $1 AND ingredient_id = $2;

-- name: UnlinkRestaurantIngredient :exec
DELETE FROM restaurant_ingredient
WHERE restaurant_id = $1 AND ingredient_id = $2;

-- name: ListIngredientsForRestaurant :many
SELECT fi.id, fi.name, fi.description, ri.is_shared
FROM food_ingredient fi
         JOIN restaurant_ingredient ri ON ri.ingredient_id = fi.id
WHERE ri.restaurant_id = $1
ORDER BY fi.name;