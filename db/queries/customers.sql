-- name: CreateCustomer :one
INSERT INTO customers (name, logo, status, email, phone, address, website, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetCustomer :one
SELECT
    c.*,
    (SELECT COUNT(*) FROM contacts WHERE customer_id = c.id) AS contact_count,
    3 AS subscription_count, -- TODO: Replace with actual count from subscriptions table
    8 AS project_count -- TODO: Replace with actual count from projects table
FROM customers c
WHERE c.id = ?;

-- name: ListCustomers :many
SELECT * FROM customers ORDER BY created_at DESC;
