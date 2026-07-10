-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1;

-- name: ListPaymentsByOrder :many
SELECT * FROM payments
WHERE order_id = $1
ORDER BY payment_time;

-- name: CreatePayment :one
INSERT INTO payments (order_id, amount, method, transaction_id, payment_time, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdatePaymentStatus :exec
UPDATE payments
SET status = $2
WHERE id = $1;

-- name: SetPaymentTransactionID :exec
UPDATE payments
SET transaction_id = $2
WHERE id = $1;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;