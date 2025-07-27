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
SELECT * FROM contacts WHERE customer_id = ? AND deleted_at IS NULL ORDER BY is_primary DESC, created_at DESC;

-- name: GetContact :one
SELECT * FROM contacts WHERE id = ? AND deleted_at IS NULL;

-- name: UpdateContact :exec
UPDATE contacts
SET name = ?, role = ?, email = ?, phone = ?, is_primary = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UnsetOtherPrimaryContacts :exec
UPDATE contacts
SET is_primary = 0
WHERE customer_id = ? AND id != ? AND deleted_at IS NULL;

-- name: UpdateContactAvatar :exec
UPDATE contacts
SET avatar = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteContactAvatar :exec
UPDATE contacts
SET avatar = NULL, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
