-- name: GetDashboardStats :one
SELECT
    (SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL) AS total_customers,
    (SELECT COUNT(*) FROM customers WHERE status = 'active' AND deleted_at IS NULL) AS active_customers,
    (SELECT COUNT(*) FROM contacts WHERE deleted_at IS NULL) AS total_contacts,
    7 AS total_projects, -- TODO:
    1247 AS monthly_revenue, -- TODO:
    15 AS revenue_change, -- TODO:
    (SELECT COUNT(*) FROM subscriptions WHERE status = 'active' AND deleted_at IS NULL) AS active_subscriptions,
    12 AS total_invoices, -- TODO:
    2 AS pending_invoices, -- TODO:
    1 AS overdue_invoices, -- TODO:
    348 AS total_invoice_amount, -- TODO:
    9 AS paid_invoices -- TODO:
;

-- name: GetRecentActivity :many
SELECT
    al.id,
    al.customer_id,
    al.activity_type,
    al.action,
    al.description,
    al.created_at,
    c.name AS customer_name
FROM
    activity_log al
JOIN
    customers c ON al.customer_id = c.id
WHERE c.deleted_at IS NULL
ORDER BY
    al.created_at DESC
LIMIT 10;
