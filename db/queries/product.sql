-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock, category_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, price, stock, category_id, created_at, updated_at;

-- name: GetProductById :one
SELECT id, name, description, price, stock, category_id, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, description, price, stock, category_id, created_at, updated_at
FROM products
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, stock = $5, category_id = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, price, stock, category_id, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;