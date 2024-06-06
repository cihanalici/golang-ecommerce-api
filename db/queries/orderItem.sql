-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_variant_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING id, order_id, product_variant_id, quantity, price, created_at, updated_at;

-- name: GetOrderItemById :one
SELECT id, order_id, product_variant_id, quantity, price, created_at, updated_at
FROM order_items
WHERE id = $1;

-- name: ListOrderItems :many
SELECT id, order_id, product_variant_id, quantity, price, created_at, updated_at
FROM order_items
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateOrderItem :one
UPDATE order_items
SET order_id = $2, product_variant_id = $3, quantity = $4, price = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, order_id, product_variant_id, quantity, price, created_at, updated_at;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;

-- name: GetOrderItemsByOrderId :many
SELECT id, order_id, product_variant_id, quantity, price, created_at, updated_at
FROM order_items
WHERE order_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
