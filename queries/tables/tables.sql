
-- name: GetTable :one
SELECT * FROM tables
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListTables :many
SELECT * FROM tables
WHERE deleted_at IS NULL
ORDER BY number;

-- name: ListTablesForAdmin :many
SELECT * FROM tables
ORDER BY deleted_at IS NULL DESC, number;

-- name: ListTablesByStatus :many
SELECT * FROM tables
WHERE status = $1 AND deleted_at IS NULL
ORDER BY number;

-- name: CreateTable :one
INSERT INTO tables (number, capacity, section, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTable :exec
UPDATE tables
SET number = $2, capacity = $3, section = $4
WHERE id = $1 AND deleted_at IS NULL;

-- name: SetTableStatus :exec
UPDATE tables
SET status = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteTable :exec
UPDATE tables
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreTable :exec
UPDATE tables
SET deleted_at = NULL
WHERE id = $1;

-- name: ListTablesWithActiveOrderFlag :many
SELECT t.id, t.number, t.capacity, t.status,
       EXISTS (
           SELECT 1 FROM food_orders fo
           WHERE fo.table_id = t.id AND fo.restaurant_id = $1 AND fo.order_status NOT IN ('delivered', 'cancelled')
       ) AS has_active_order
FROM tables t
WHERE t.deleted_at IS NULL
ORDER BY t.number;

-- name: ListTablesWithOrderStatusForRestaurant :many
SELECT t.id, t.number, t.capacity, t.status,
       EXISTS (
           SELECT 1 FROM food_orders fo
           WHERE fo.table_id = t.id AND fo.restaurant_id = $1 AND fo.order_status NOT IN ('delivered', 'cancelled')
       ) AS occupied_by_restaurant
FROM tables t
WHERE t.deleted_at IS NULL
ORDER BY t.number;

-- name: ListTablesWithGlobalOrderStatus :many
SELECT t.id, t.number, t.capacity, t.status,
       EXISTS (
           SELECT 1 FROM food_orders fo
           WHERE fo.table_id = t.id AND fo.order_status NOT IN ('delivered', 'cancelled')
       ) AS occupied_by_any_restaurant
FROM tables t
WHERE t.deleted_at IS NULL
ORDER BY t.number;

-- name: TableHasActiveOrders :one
SELECT EXISTS (
    SELECT 1 FROM food_orders
    WHERE table_id = $1 AND order_status NOT IN ('delivered', 'cancelled')
) AS has_active_orders;