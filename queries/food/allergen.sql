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
INSERT INTO allergen (name, name_alt, description, description_alt, code)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: UpdateAllergen :exec
UPDATE allergen
SET name = $2, name_alt = $3, description = $4, description_alt = $5, code = $6
WHERE id = $1;

-- name: DeleteAllergen :exec
DELETE FROM allergen
WHERE id = $1;