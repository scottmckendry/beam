-- name: ListSubscriptionsByCustomer :many
SELECT s.*, s.start_date as next_billing_date FROM subscriptions s WHERE customer_id = ? AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: CreateSubscription :one
INSERT INTO subscriptions ( customer_id, description, amount, term, billing_cadence, status, start_date, notes)
VALUES ( ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: GetSubscription :one
SELECT * FROM subscriptions WHERE id = ? AND deleted_at IS NULL;

-- name: UpdateSubscription :one
UPDATE subscriptions SET description = ?, amount = ?, term = ?, billing_cadence = ?, status = ?, start_date = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSubscription :one
UPDATE subscriptions SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL RETURNING *;
