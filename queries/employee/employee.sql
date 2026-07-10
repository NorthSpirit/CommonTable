-- name: GetEmployee :one
SELECT * FROM employee
WHERE id = $1;

-- name: ListEmployeesByRestaurant :many
SELECT * FROM employee
WHERE restaurant_id = $1
ORDER BY name;

-- name: ListActiveEmployeesByRestaurant :many
SELECT * FROM employee
WHERE restaurant_id = $1 AND is_active = true
ORDER BY name;

-- name: CreateEmployee :one
INSERT INTO employee (restaurant_id, name, role, pin_hash, hire_date, email, phone_number)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateEmployee :exec
UPDATE employee
SET name = $2, role = $3, email = $4, phone_number = $5
WHERE id = $1;

-- name: UpdateEmployeePinHash :exec
UPDATE employee
SET pin_hash = $2
WHERE id = $1;

-- name: SetEmployeeActive :exec
UPDATE employee
SET is_active = $2
WHERE id = $1;

-- name: UpdateEmployeeLastLogin :exec
UPDATE employee
SET last_login = $2
WHERE id = $1;

-- name: DeleteEmployee :exec
DELETE FROM employee
WHERE id = $1;