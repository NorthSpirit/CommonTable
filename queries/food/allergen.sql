-- name: GetAllergen :one
SELECT * FROM allergen
WHERE id = $1;

-- name: GetAllergenByCode :one
SELECT * FROM allergen
WHERE code = $1;

-- name: ListAllergens :many
SELECT * FROM allergen
ORDER BY name;

-- name: CreateAllergen :one
INSERT INTO allergen (name, description, code)
VALUES ($1, $2, $3)
    RETURNING *;

-- name: UpdateAllergen :exec
UPDATE allergen
SET name = $2, description = $3, code = $4
WHERE id = $1;

-- name: DeleteAllergen :exec
DELETE FROM allergen
WHERE id = $1;