-- name: GetFoodItemForExternalOrder :one
SELECT id, name, price, vat, is_available
FROM food_items
WHERE id = $1 AND restaurant_id = $2 AND deleted_at IS NULL;

-- name: CreateUnapprovedOrder :one
INSERT INTO unapproved_orders (restaurant_id, api_key_id, customer_name, customer_phone, notes, idempotency_key)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateUnapprovedOrderItem :exec
INSERT INTO unapproved_order_items (unapproved_order_id, item_id, quantity, unit_price_snapshot, vat_snapshot, item_name_snapshot, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetUnapprovedOrderByIdempotencyKey :one
SELECT * FROM unapproved_orders
WHERE api_key_id = $1 AND idempotency_key = $2;

-- name: CountRecentUnapprovedOrdersByApiKey :one
SELECT count(*)
FROM unapproved_orders
WHERE api_key_id = $1
  AND created_at > now() - make_interval(secs => sqlc.arg(window_seconds)::int);

-- name: CountPendingUnapprovedOrdersByApiKey :one
SELECT count(*)
FROM unapproved_orders
WHERE api_key_id = $1 AND status = 'pending';

-- name: CountPendingUnapprovedOrdersByPhone :one
SELECT count(*)
FROM unapproved_orders
WHERE restaurant_id = $1 AND customer_phone = $2 AND status = 'pending';

-- name: ListPendingUnapprovedOrders :many
SELECT * FROM unapproved_orders
WHERE restaurant_id = $1 AND status = 'pending'
ORDER BY created_at;

-- name: ListUnapprovedOrderItems :many
SELECT * FROM unapproved_order_items
WHERE unapproved_order_id = $1
ORDER BY id;

-- name: GetUnapprovedOrderForUpdate :one
SELECT * FROM unapproved_orders
WHERE id = $1 AND restaurant_id = $2
FOR UPDATE;

-- name: MarkUnapprovedOrderApproved :execrows
UPDATE unapproved_orders
SET status = 'approved', approved_order_id = $2, decided_by = $3, decided_at = now()
WHERE id = $1 AND status = 'pending';

-- name: MarkUnapprovedOrderRejected :execrows
UPDATE unapproved_orders
SET status = 'rejected', rejection_reason = $2, decided_by = $3, decided_at = now()
WHERE id = $1 AND status = 'pending';
