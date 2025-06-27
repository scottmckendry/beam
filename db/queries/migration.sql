-- name: CreateMigrationsTable :exec
CREATE TABLE IF NOT EXISTS migrations (
    id UUID PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL UNIQUE,
    applied DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- name: GetMigration :one
SELECT * FROM migrations
WHERE name = ? LIMIT 1;

-- name: ListMigrations :many
SELECT * FROM migrations;

-- name: ApplyMigration :exec
INSERT INTO migrations (name, applied)
VALUES (?, datetime('now'));
