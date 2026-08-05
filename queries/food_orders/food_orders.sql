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

-- name: ListActiveOrdersByRestaurantForKitchen :many
SELECT fo.id, fo.order_number, fo.order_status, fo.order_creation_date, fo.notes, t.number AS table_number
FROM food_orders fo
         LEFT JOIN tables t ON t.id = fo.table_id
WHERE fo.restaurant_id = $1 AND fo.order_status NOT IN ('delivered', 'cancelled')
ORDER BY fo.order_creation_date;

-- name: ListOrdersByRestaurantForAdmin :many
SELECT fo.id, fo.order_number, fo.order_status, fo.order_creation_date, fo.notes, t.number AS table_number, e.name AS waiter_name
FROM food_orders fo
         LEFT JOIN tables t ON t.id = fo.table_id
         JOIN employee e ON e.id = fo.waiter_id
WHERE fo.restaurant_id = $1
ORDER BY fo.order_creation_date DESC;

-- name: GetOrderForAdmin :one
SELECT fo.id, fo.restaurant_id, fo.order_number, fo.order_status, fo.order_creation_date, fo.notes, t.number AS table_number, e.name AS waiter_name
FROM food_orders fo
         LEFT JOIN tables t ON t.id = fo.table_id
         JOIN employee e ON e.id = fo.waiter_id
WHERE fo.id = $1;

-- name: GetOrderByIdempotencyKey :one
SELECT * FROM food_orders
WHERE idempotency_key = $1;

-- name: CreateFoodOrderWithIdempotencyKey :one
INSERT INTO food_orders (restaurant_id, table_id, waiter_id, order_creation_date, order_status, order_number, notes, idempotency_key)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetNextOrderNumberForToday :one
SELECT COALESCE(MAX(order_number), 0) + 1
FROM food_orders
WHERE restaurant_id = $1
  AND order_creation_date >= CURRENT_DATE
  AND order_creation_date < CURRENT_DATE + INTERVAL '1 day';

-- name: GetActiveOrderForTable :one
SELECT * FROM food_orders
WHERE table_id = $1 AND restaurant_id = $2 AND order_status NOT IN ('delivered', 'cancelled')
ORDER BY order_creation_date DESC
LIMIT 1;

-- name: ListActiveOrdersByWaiter :many
SELECT fo.id, fo.order_number, fo.order_status, fo.order_creation_date, fo.notes, t.number AS table_number
FROM food_orders fo
         LEFT JOIN tables t ON t.id = fo.table_id
WHERE fo.waiter_id = $1 AND fo.order_status NOT IN ('delivered', 'cancelled')
ORDER BY fo.order_creation_date;