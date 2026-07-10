-- name: GetTable :one
SELECT * FROM tables
WHERE id = $1;

-- name: ListTables :many
SELECT * FROM tables
ORDER BY number;

-- name: ListTablesByStatus :many
SELECT * FROM tables
WHERE status = $1
ORDER BY number;

-- name: CreateTable :one
INSERT INTO tables (number, capacity, section, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTable :exec
UPDATE tables
SET number = $2, capacity = $3, section = $4
WHERE id = $1;

-- name: SetTableStatus :exec
UPDATE tables
SET status = $2
WHERE id = $1;

-- name: DeleteTable :exec
DELETE FROM tables
WHERE id = $1;