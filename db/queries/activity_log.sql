-- name: LogActivity :one
INSERT INTO activity_log (customer_id, activity_type, action, description)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: ListRecentActivity :many
SELECT * FROM activity_log ORDER BY created_at DESC LIMIT 50;
