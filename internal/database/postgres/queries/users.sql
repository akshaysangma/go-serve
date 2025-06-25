-- name: CreateUser :one
INSERT INTO users (username, email)
VALUES ($1,$2)
RETURNING id, username, email, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, username, email, created_at, updated_at FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users SET username = $2, email = $3, updated_at = NOW() WHERE id = $1 RETURNING id, username, email, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users where id = $1;
