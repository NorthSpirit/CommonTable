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
SELECT cm.id,
       CASE WHEN cm.deleted_at IS NULL THEN cm.message ELSE ''::text END::text AS message,
       cm.sent_at,
       e.name AS sender_name,
       cm.deleted_at,
       cm.deleted_by,
       cm.delete_reason
FROM chat_messages cm
         JOIN employee e ON e.id = cm.employee_id
WHERE cm.restaurant_id = $1
ORDER BY cm.sent_at, cm.id;

-- name: ListMessagesByRestaurantAfterID :many
SELECT cm.id,
       CASE WHEN cm.deleted_at IS NULL THEN cm.message ELSE ''::text END::text AS message,
       cm.sent_at,
       e.name AS sender_name,
       cm.deleted_at,
       cm.deleted_by,
       cm.delete_reason
FROM chat_messages cm
         JOIN employee e ON e.id = cm.employee_id
WHERE cm.restaurant_id = $1 AND cm.id > $2
ORDER BY cm.id;

-- name: SoftDeleteChatMessage :one
UPDATE chat_messages
SET deleted_at    = now(),
    deleted_by    = $3,
    delete_reason = $4
WHERE id = $1 AND restaurant_id = $2 AND deleted_at IS NULL
RETURNING id, restaurant_id, employee_id, sent_at, deleted_at, deleted_by, delete_reason;

-- name: ListDeletedMessagesByRestaurant :many
SELECT cm.id, cm.deleted_by, cm.delete_reason
FROM chat_messages cm
WHERE cm.restaurant_id = $1 AND cm.deleted_at IS NOT NULL;
