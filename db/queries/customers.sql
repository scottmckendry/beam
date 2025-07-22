-- name: CreateCustomer :one
INSERT INTO customers (name, logo, status, email, phone, address, website, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateCustomer :one
UPDATE customers
SET name = ?, logo = ?, status = ?, email = ?, phone = ?, address = ?, website = ?, notes = ?
WHERE id = ?
RETURNING *;

-- name: DeleteCustomer :one
UPDATE customers
SET deleted_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: GetCustomer :one
SELECT
    c.*,
    (SELECT COUNT(*) FROM contacts WHERE customer_id = c.id AND deleted_at IS NULL) AS contact_count,
    -- TODO: Replace these with actual counts from the respective tables
    3 AS subscription_count,
    8 AS project_count,
    238 AS subscription_revenue,
    267 AS monthly_revenue,
    15 AS revenue_change
FROM customers c
WHERE c.id = ? AND c.deleted_at IS NULL;

-- name: ListCustomers :many
SELECT * FROM customers WHERE deleted_at IS NULL ORDER BY created_at DESC;

-- name: UpdateCustomerLogo :exec
UPDATE customers
SET logo = ?
WHERE id = ?;
