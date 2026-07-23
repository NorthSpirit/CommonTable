-- name: GetChatMessage :one
SELECT * FROM chat_messages
WHERE id = $1;

-- name: ListMessagesByRestaurant :many
SELECT * FROM chat_messages
WHERE restaurant_id = $1
ORDER BY sent_at;

-- name: ListMessagesByEmployee :many
SELECT * FROM chat_messages
WHERE employee_id = $1
ORDER BY sent_at DESC;

-- name: CreateChatMessage :one
INSERT INTO chat_messages (restaurant_id, employee_id, message, sent_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteChatMessage :exec
DELETE FROM chat_messages
WHERE id = $1;

-- name: ListMessagesByRestaurantWithSender :many
SELECT cm.id, cm.message, cm.sent_at, e.name AS sender_name
FROM chat_messages cm
         JOIN employee e ON e.id = cm.employee_id
WHERE cm.restaurant_id = $1
ORDER BY cm.sent_at, cm.id;

-- name: ListMessagesByRestaurantAfterID :many
SELECT cm.id, cm.message, cm.sent_at, e.name AS sender_name
FROM chat_messages cm
         JOIN employee e ON e.id = cm.employee_id
WHERE cm.restaurant_id = $1 AND cm.id > $2
ORDER BY cm.id;