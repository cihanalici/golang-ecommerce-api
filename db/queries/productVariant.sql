-- name: CreateProductVariant :one
INSERT INTO product_variants (product_id, color, size, stock, price)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, product_id, color, size, stock, price, created_at, updated_at;

-- name: GetProductVariantById :one
SELECT id, product_id, color, size, stock, price, created_at, updated_at
FROM product_variants
WHERE id = $1;

-- name: ListProductVariants :many
SELECT id, product_id, color, size, stock, price, created_at, updated_at
FROM product_variants
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateProductVariant :one
UPDATE product_variants
SET product_id = $2, color = $3, size = $4, stock = $5, price = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, product_id, color, size, stock, price, created_at, updated_at;

-- name: DeleteProductVariant :exec
DELETE FROM product_variants
WHERE id = $1;