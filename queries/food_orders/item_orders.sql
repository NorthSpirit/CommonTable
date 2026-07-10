-- name: GetItemOrder :one
SELECT * FROM item_orders
WHERE id = $1;

-- name: ListItemsForOrder :many
SELECT * FROM item_orders
WHERE order_id = $1;

-- name: AddItemToOrder :one
INSERT INTO item_orders (order_id, item_id, quantity, unit_price_snapshot, vat_snapshot, notes, item_name_snapshot, item_description_snapshot)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateItemOrderQuantity :exec
UPDATE item_orders
SET quantity = $2
WHERE id = $1;

-- name: UpdateItemOrderNotes :exec
UPDATE item_orders
SET notes = $2
WHERE id = $1;

-- name: RemoveItemFromOrder :exec
DELETE FROM item_orders
WHERE id = $1;

-- name: ClearOrderItems :exec
DELETE FROM item_orders
WHERE order_id = $1;