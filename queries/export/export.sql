-- name: GetRestaurantForExport :one
SELECT id, name, description, description_alt
FROM restaurant
WHERE id = $1 AND deleted_at IS NULL;