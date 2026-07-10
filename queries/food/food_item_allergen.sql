-- name: AddAllergenToItem :exec
INSERT INTO food_item_allergen (item_id, allergen_id)
VALUES ($1, $2);

-- name: RemoveAllergenFromItem :exec
DELETE FROM food_item_allergen
WHERE item_id = $1 AND allergen_id = $2;

-- name: ListAllergensForItem :many
SELECT a.*
FROM allergen a
         JOIN food_item_allergen fia ON fia.allergen_id = a.id
WHERE fia.item_id = $1
ORDER BY a.name;

-- name: ListItemsWithAllergen :many
SELECT fi.id, fi.restaurant_id, fi.name, fi.description, fi.is_available, fi.image_url, fi.prep_time_mins, fi.instructions
FROM food_items fi
         JOIN food_item_allergen fia ON fia.item_id = fi.id
WHERE fia.allergen_id = $1
ORDER BY fi.name;