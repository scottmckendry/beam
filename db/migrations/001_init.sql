CREATE TABLE IF NOT EXISTS migrations (
    id UUID PRIMARY KEY DEFAULT (uuid()),
    name TEXT NOT NULL UNIQUE,
    applied DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT (uuid()),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    github_id TEXT NOT NULL UNIQUE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT (uuid()),
    name TEXT NOT NULL,
    logo TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    email TEXT,
    phone TEXT,
    address TEXT,
    website TEXT,
    notes TEXT,
    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now')),
    deleted_at DATETIME DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS contacts (
    id UUID PRIMARY KEY DEFAULT (uuid()),
    customer_id UUID NOT NULL,
    name TEXT NOT NULL,
    role TEXT,
    email TEXT,
    phone TEXT,
    avatar TEXT,
    is_primary BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now')),
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);


CREATE TABLE IF NOT EXISTS activity_log (
    id UUID PRIMARY KEY DEFAULT (uuid()),
    customer_id UUID NOT NULL,
    activity_type TEXT NOT NULL,
    action TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT (datetime('now')),
    -- TODO: deleting a customer also deletes its activity log
    -- implement soft delete for customers and anywhere else it makes sense
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT (uuid()),
    customer_id UUID NOT NULL,
    description TEXT NOT NULL,
    amount NUMERIC NOT NULL,
    term TEXT NOT NULL, -- e.g. 'yearly', 'monthly', 'quarterly'
    billing_cadence TEXT NOT NULL, -- e.g. 'monthly', 'yearly', 'quarterly'
    start_date DATETIME NOT NULL,
    end_date DATETIME DEFAULT NULL,
    status TEXT NOT NULL DEFAULT 'active', -- 'active', 'paused', 'cancelled'
    notes TEXT,
    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now')),
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);
