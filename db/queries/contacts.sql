-- name: CreateContact :one
INSERT INTO contacts (customer_id, name, role, email, phone, avatar, is_primary, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListContactsByCustomer :many
SELECT * FROM contacts WHERE customer_id = ? ORDER BY created_at DESC;
