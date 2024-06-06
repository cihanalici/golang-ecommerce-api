-- name: CreatePasswordReset :one
INSERT INTO password_resets (user_id, reset_token, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, reset_token, expires_at, created_at;

-- name: GetPasswordResetByToken :one
SELECT id, user_id, reset_token, expires_at, created_at
FROM password_resets
WHERE reset_token = $1;

-- name: DeletePasswordReset :exec
DELETE FROM password_resets
WHERE reset_token = $1;

-- name: DeleteExpiredPasswordResets :exec
DELETE FROM password_resets
WHERE expires_at < CURRENT_TIMESTAMP;

-- name: GetPasswordResetByUserId :one
SELECT id, user_id, reset_token, expires_at, created_at
FROM password_resets
WHERE user_id = $1;

-- name: GetPasswordResetByUserIdAndToken :one
SELECT id, user_id, reset_token, expires_at, created_at
FROM password_resets
WHERE user_id = $1 AND reset_token = $2;

-- name: GetPasswordResetByID :one
SELECT id, user_id, reset_token, created_at, expires_at
FROM password_resets
WHERE id = $1;
