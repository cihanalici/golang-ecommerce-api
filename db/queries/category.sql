-- name: CreateCategory :one
INSERT INTO categories (name, description)
VALUES ($1, $2)
RETURNING id, name, description, created_at, updated_at;

-- name: GetCategoryById :one
SELECT id, name, description, created_at, updated_at
FROM categories
WHERE id = $1;

-- name: ListCategories :many
SELECT id, name, description, created_at, updated_at
FROM categories
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, created_at, updated_at;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

