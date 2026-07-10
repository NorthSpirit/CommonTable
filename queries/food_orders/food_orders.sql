-- name: GetFoodOrder :one
SELECT * FROM food_orders
WHERE id = $1;

-- name: ListOrdersByRestaurant :many
SELECT * FROM food_orders
WHERE restaurant_id = $1
ORDER BY order_creation_date DESC;

-- name: ListOrdersByTable :many
SELECT * FROM food_orders
WHERE table_id = $1
ORDER BY order_creation_date DESC;

-- name: ListOrdersByStatus :many
SELECT * FROM food_orders
WHERE restaurant_id = $1 AND order_status = $2
ORDER BY order_creation_date DESC;

-- name: CreateFoodOrder :one
INSERT INTO food_orders (restaurant_id, table_id, waiter_id, order_creation_date, order_status, order_number, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateOrderStatus :exec
UPDATE food_orders
SET order_status = $2
WHERE id = $1;

-- name: UpdateFoodOrderNotes :exec
UPDATE food_orders
SET notes = $2
WHERE id = $1;

-- name: DeleteFoodOrder :exec
DELETE FROM food_orders
WHERE id = $1;