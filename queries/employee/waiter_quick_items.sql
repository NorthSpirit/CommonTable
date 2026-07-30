-- name: CountWaiterQuickItems :one
SELECT COUNT(*) FROM waiter_quick_items
WHERE employee_id = $1;

-- name: AddWaiterQuickItem :exec
INSERT INTO waiter_quick_items (employee_id, item_id, sort_order)
VALUES ($1, $2, $3);

-- name: RemoveWaiterQuickItem :exec
DELETE FROM waiter_quick_items
WHERE employee_id = $1 AND item_id = $2;

-- name: ClearWaiterQuickItems :exec
DELETE FROM waiter_quick_items
WHERE employee_id = $1;

-- name: UpdateWaiterQuickItemSortOrder :exec
UPDATE waiter_quick_items
SET sort_order = $3
WHERE employee_id = $1 AND item_id = $2;

-- name: ListWaiterQuickItems :many
SELECT fi.id, fi.name, fi.image_url, fi.price, wqi.sort_order
FROM waiter_quick_items wqi
         JOIN food_items fi ON fi.id = wqi.item_id
WHERE wqi.employee_id = $1
  AND fi.is_available = true
  AND fi.deleted_at IS NULL
ORDER BY wqi.sort_order, fi.name;

