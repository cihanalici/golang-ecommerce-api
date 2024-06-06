-- name: CreateOrder :one
INSERT INTO orders (user_id, total_amount, status)
VALUES ($1, $2, $3)
RETURNING id, user_id, total_amount, status, created_at, updated_at;

-- name: GetOrderById :one
SELECT id, user_id, total_amount, status, created_at, updated_at
FROM orders
WHERE id = $1;

-- name: ListOrders :many
SELECT id, user_id, total_amount, status, created_at, updated_at
FROM orders
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateOrder :one
UPDATE orders
SET user_id = $2, total_amount = $3, status = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, total_amount, status, created_at, updated_at;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: GetOrdersByUserId :many
SELECT id, user_id, total_amount, status, created_at, updated_at
FROM orders
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;