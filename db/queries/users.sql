-- name: InsertUser :exec
INSERT INTO users (name, email, github_id, is_admin) VALUES (?, ?, ?, 0)
ON CONFLICT(github_id) DO NOTHING;

-- name: GetUserByGithubID :one
SELECT * FROM users WHERE github_id = ? LIMIT 1;

-- name: IsUserAdmin :one
SELECT is_admin FROM users WHERE github_id = ? LIMIT 1;
