-- name: LinkRestaurantAllergen :exec
INSERT INTO restaurant_allergen (restaurant_id, allergen_id)
VALUES ($1, $2);

-- name: UnlinkRestaurantAllergen :exec
DELETE FROM restaurant_allergen
WHERE restaurant_id = $1 AND allergen_id = $2;

-- name: ListAllergensForRestaurant :many
SELECT a.*
FROM allergen a
         JOIN restaurant_allergen ra ON ra.allergen_id = a.id
WHERE ra.restaurant_id = $1
ORDER BY a.name;