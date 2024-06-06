-- name: CreateWishlistItem :one
INSERT INTO wishlist (user_id, product_id)
VALUES ($1, $2)
RETURNING id, user_id, product_id, created_at, updated_at;

-- name: GetWishlistItemById :one
SELECT id, user_id, product_id, created_at, updated_at
FROM wishlist
WHERE id = $1;

-- name: ListWishlistItems :many
SELECT id, user_id, product_id, created_at, updated_at
FROM wishlist
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateWishlistItem :one
UPDATE wishlist
SET user_id = $2, product_id = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, product_id, created_at, updated_at;

-- name: DeleteWishlistItem :exec
DELETE FROM wishlist
WHERE id = $1;

-- name: GetWishlistItemsByUserId :many
SELECT id, user_id, product_id, created_at, updated_at
FROM wishlist
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
