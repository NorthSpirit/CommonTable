-- name: AddFoodItemToMenu :exec
INSERT INTO food_item_to_menu (menu_id, item_id)
VALUES ($1, $2);

-- name: RemoveFoodItemFromMenu :exec
DELETE FROM food_item_to_menu
WHERE menu_id = $1 AND item_id = $2;

-- name: ListItemsInMenu :many
SELECT fi.id, fi.restaurant_id, fi.name, fi.description, fi.is_available, fi.image_url, fi.prep_time_mins
FROM food_items fi
         JOIN food_item_to_menu fitm ON fitm.item_id = fi.id
WHERE fitm.menu_id = $1
ORDER BY fi.name;

-- name: ListMenusForItem :many
SELECT fm.*
FROM food_menu fm
         JOIN food_item_to_menu fitm ON fitm.menu_id = fm.id
WHERE fitm.item_id = $1
ORDER BY fm.name;

-- name: ListItemsInMenuForAdmin :many
SELECT fi.id, fi.name, fi.image_url, fi.price, fi.is_available, fi.deleted_at
FROM food_items fi
         JOIN food_item_to_menu fitm ON fitm.item_id = fi.id
WHERE fitm.menu_id = $1
ORDER BY fi.name;