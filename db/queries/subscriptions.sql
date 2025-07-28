-- name: ListSubscriptionsByCustomer :many
SELECT s.*, s.start_date as next_billing_date FROM subscriptions s WHERE customer_id = ? AND deleted_at IS NULL ORDER BY created_at DESC;
