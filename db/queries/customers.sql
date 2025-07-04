-- name: CreateCustomer :one
INSERT INTO customers (name, logo, status, email, phone, address, website, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers WHERE id = ?;

-- name: ListCustomers :many
SELECT * FROM customers ORDER BY created_at DESC;
