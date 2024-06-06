-- name: CreateReview :one
INSERT INTO reviews (product_id, user_id, rating, comment)
VALUES ($1, $2, $3, $4)
RETURNING id, product_id, user_id, rating, comment, created_at, updated_at;

-- name: GetReviewById :one
SELECT id, product_id, user_id, rating, comment, created_at, updated_at
FROM reviews
WHERE id = $1;

-- name: ListReviews :many
SELECT id, product_id, user_id, rating, comment, created_at, updated_at
FROM reviews
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateReview :one
UPDATE reviews
SET product_id = $2, user_id = $3, rating = $4, comment = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, product_id, user_id, rating, comment, created_at, updated_at;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;

-- name: GetReviewsByProductId :many
SELECT id, product_id, user_id, rating, comment, created_at, updated_at
FROM reviews
WHERE product_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;