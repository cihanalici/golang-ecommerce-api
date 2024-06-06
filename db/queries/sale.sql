-- name: CreateSale :one
INSERT INTO sales (month, year, total_sales)
VALUES ($1, $2, $3)
RETURNING id, month, year, total_sales, created_at, updated_at;

-- name: GetSaleById :one
SELECT id, month, year, total_sales, created_at, updated_at
FROM sales
WHERE id = $1;

-- name: ListSales :many
SELECT id, month, year, total_sales, created_at, updated_at
FROM sales
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateSale :one
UPDATE sales
SET month = $2, year = $3, total_sales = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, month, year, total_sales, created_at, updated_at;

-- name: DeleteSale :exec
DELETE FROM sales
WHERE id = $1;
