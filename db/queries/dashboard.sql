-- name: GetDashboardStats :one
SELECT
    (SELECT COUNT(*) FROM customers) AS total_customers,
    (SELECT COUNT(*) FROM customers WHERE status = 'active') AS active_customers,
    (SELECT COUNT(*) FROM contacts) AS total_contacts,
    7 AS total_projects, -- hardcoded
    1247 AS monthly_revenue, -- hardcoded
    15 AS revenue_change, -- hardcoded
    3 AS active_subscriptions, -- hardcoded
    12 AS total_invoices, -- hardcoded
    2 AS pending_invoices, -- hardcoded
    1 AS overdue_invoices, -- hardcoded
    348 AS total_invoice_amount, -- hardcoded
    9 AS paid_invoices -- hardcoded
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
ORDER BY
    al.created_at DESC
LIMIT 10;
