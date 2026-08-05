-- name: GetFoodItem :one
SELECT id, restaurant_id, name, name_alt, description, description_alt, is_available, image_url, prep_time_mins, instructions, instructions_alt, price, deleted_at
FROM food_items
WHERE id = $1;

-- name: GetFoodItemForDisplay :one
SELECT
    id,
    restaurant_id,
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(name_alt, name) ELSE name END::varchar AS name,
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(description_alt, description) ELSE description END::text AS description,
    is_available,
    image_url,
    prep_time_mins,
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(instructions_alt, instructions) ELSE instructions END::text AS instructions
FROM food_items
WHERE id = sqlc.arg(id)::integer AND deleted_at IS NULL;

-- name: ListFoodItemsByRestaurant :many
SELECT id, restaurant_id, name, name_alt, description, description_alt, is_available, image_url, prep_time_mins, instructions, instructions_alt, price
FROM food_items
WHERE restaurant_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: ListFoodItemsByRestaurantForAdmin :many
SELECT id, restaurant_id, name, name_alt, description, description_alt, is_available, image_url, prep_time_mins, instructions, instructions_alt, price, deleted_at
FROM food_items
WHERE restaurant_id = $1
ORDER BY deleted_at IS NULL DESC, name;

-- name: ListAvailableFoodItemsByRestaurant :many
SELECT id, restaurant_id, name, name_alt, description, description_alt, is_available, image_url, prep_time_mins, instructions, instructions_alt, price
FROM food_items
WHERE restaurant_id = $1 AND is_available = true AND deleted_at IS NULL
ORDER BY name;

-- name: CreateFoodItem :one
INSERT INTO food_items (restaurant_id, name, name_alt, description, description_alt, is_available, image_url, prep_time_mins, instructions, instructions_alt, price)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, restaurant_id, name, name_alt, description, description_alt, is_available, image_url, prep_time_mins, instructions, instructions_alt, price;

-- name: UpdateFoodItem :exec
UPDATE food_items
SET name = $2, name_alt = $3, description = $4, description_alt = $5, is_available = $6, image_url = $7, prep_time_mins = $8, instructions = $9, instructions_alt = $10, price = $11
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

-- name: ListMenuItemsForDisplay :many
SELECT
    fi.id,
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(fi.name_alt, fi.name) ELSE fi.name END::varchar AS name,
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(fi.description_alt, fi.description) ELSE fi.description END::text AS description,
    fi.price,
    fi.image_url,
    fi.prep_time_mins,
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(fi.instructions_alt, fi.instructions) ELSE fi.instructions END::text AS instructions
FROM food_items fi
         JOIN food_item_to_menu fitm ON fitm.item_id = fi.id
WHERE fitm.menu_id = sqlc.arg(menu_id)::integer AND fi.is_available = true AND fi.deleted_at IS NULL
ORDER BY fi.name;

-- name: ListMostOrderedItemsByRestaurant :many
SELECT
    fi.id,
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(fi.name_alt, fi.name) ELSE fi.name END::varchar AS name,
    fi.image_url,
    COUNT(io.id) AS order_count
FROM food_items fi
         JOIN item_orders io ON io.item_id = fi.id
WHERE fi.restaurant_id = sqlc.arg(restaurant_id)::integer
  AND fi.is_available = true
  AND fi.deleted_at IS NULL
GROUP BY fi.id, fi.name, fi.name_alt, fi.image_url
ORDER BY order_count DESC
LIMIT sqlc.arg(item_limit)::integer;

-- name: CountFoodItems :one
SELECT COUNT(*) FROM food_items
WHERE deleted_at IS NULL;

-- name: SetFoodItemInstructions :exec
UPDATE food_items
SET instructions = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListMenuItemsForWaiter :many
SELECT fi.id, fi.name, fi.description, fi.price, fi.image_url, fi.is_available
FROM food_items fi
         JOIN food_item_to_menu fitm ON fitm.item_id = fi.id
WHERE fitm.menu_id = $1 AND fi.deleted_at IS NULL
ORDER BY fi.name;

-- name: SearchFoodItemsForWaiter :many
SELECT DISTINCT fi.id, fi.name, fi.description, fi.price, fi.image_url, fi.is_available
FROM food_items fi
         JOIN food_item_to_menu fitm ON fitm.item_id = fi.id
         JOIN food_menu fm ON fm.id = fitm.menu_id
WHERE fm.restaurant_id = $1
  AND fi.deleted_at IS NULL
  AND fi.name ILIKE '%' || $2 || '%'
ORDER BY fi.name;