-- name: CreateUser :one
INSERT INTO users (name, email, password, role, addresses)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, email, password, role, addresses, created_at, updated_at;

-- name: GetUserById :one
SELECT id, name, email, password, role, addresses, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, password, role, addresses, created_at, updated_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, name, email, role, password, addresses, created_at, updated_at
FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET name = $1, email = $2, role = $3, addresses = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $5
RETURNING id, name, email, role, password, addresses, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
