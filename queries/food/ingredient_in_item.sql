-- name: AddIngredientToItem :exec
INSERT INTO ingredient_in_item (ingredient_id, item_id, measurement, amount)
VALUES ($1, $2, $3, $4);

-- name: UpdateIngredientInItem :exec
UPDATE ingredient_in_item
SET measurement = $3, amount = $4
WHERE ingredient_id = $1 AND item_id = $2;

-- name: RemoveIngredientFromItem :exec
DELETE FROM ingredient_in_item
WHERE ingredient_id = $1 AND item_id = $2;

-- name: ListIngredientsForItem :many
SELECT fi.id, fi.name, fi.description, iii.measurement, iii.amount
FROM food_ingredient fi
         JOIN ingredient_in_item iii ON iii.ingredient_id = fi.id
WHERE iii.item_id = $1
ORDER BY fi.name;

-- name: ListItemsUsingIngredient :many
SELECT fitem.id, fitem.name, iii.measurement, iii.amount
FROM food_items fitem
         JOIN ingredient_in_item iii ON iii.item_id = fitem.id
WHERE iii.ingredient_id = $1
ORDER BY fitem.name;