-- name: CreateContact :one
INSERT INTO contacts (customer_id, name, role, email, phone, avatar, is_primary, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteContact :one
UPDATE contacts
SET deleted_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: DeleteContactsByCustomer :exec
UPDATE contacts SET deleted_at = datetime('now') WHERE customer_id = ?;

-- name: ListContactsByCustomer :many
SELECT * FROM contacts WHERE customer_id = ? AND deleted_at IS NULL ORDER BY created_at DESC;
