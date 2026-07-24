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