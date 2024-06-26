// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: passwordReset.sql

package sqlc

import (
	"context"
	"time"
)

const createPasswordReset = `-- name: CreatePasswordReset :one
INSERT INTO password_resets (user_id, reset_token, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, reset_token, expires_at, created_at
`

type CreatePasswordResetParams struct {
	UserID     int32     `json:"user_id"`
	ResetToken string    `json:"reset_token"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type CreatePasswordResetRow struct {
	ID         int32     `json:"id"`
	UserID     int32     `json:"user_id"`
	ResetToken string    `json:"reset_token"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (q *Queries) CreatePasswordReset(ctx context.Context, arg CreatePasswordResetParams) (CreatePasswordResetRow, error) {
	row := q.db.QueryRowContext(ctx, createPasswordReset, arg.UserID, arg.ResetToken, arg.ExpiresAt)
	var i CreatePasswordResetRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ResetToken,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteExpiredPasswordResets = `-- name: DeleteExpiredPasswordResets :exec
DELETE FROM password_resets
WHERE expires_at < CURRENT_TIMESTAMP
`

func (q *Queries) DeleteExpiredPasswordResets(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteExpiredPasswordResets)
	return err
}

const deletePasswordReset = `-- name: DeletePasswordReset :exec
DELETE FROM password_resets
WHERE reset_token = $1
`

func (q *Queries) DeletePasswordReset(ctx context.Context, resetToken string) error {
	_, err := q.db.ExecContext(ctx, deletePasswordReset, resetToken)
	return err
}

const getPasswordResetByID = `-- name: GetPasswordResetByID :one
SELECT id, user_id, reset_token, created_at, expires_at
FROM password_resets
WHERE id = $1
`

func (q *Queries) GetPasswordResetByID(ctx context.Context, id int32) (PasswordReset, error) {
	row := q.db.QueryRowContext(ctx, getPasswordResetByID, id)
	var i PasswordReset
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ResetToken,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getPasswordResetByToken = `-- name: GetPasswordResetByToken :one
SELECT id, user_id, reset_token, expires_at, created_at
FROM password_resets
WHERE reset_token = $1
`

type GetPasswordResetByTokenRow struct {
	ID         int32     `json:"id"`
	UserID     int32     `json:"user_id"`
	ResetToken string    `json:"reset_token"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (q *Queries) GetPasswordResetByToken(ctx context.Context, resetToken string) (GetPasswordResetByTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getPasswordResetByToken, resetToken)
	var i GetPasswordResetByTokenRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ResetToken,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getPasswordResetByUserId = `-- name: GetPasswordResetByUserId :one
SELECT id, user_id, reset_token, expires_at, created_at
FROM password_resets
WHERE user_id = $1
`

type GetPasswordResetByUserIdRow struct {
	ID         int32     `json:"id"`
	UserID     int32     `json:"user_id"`
	ResetToken string    `json:"reset_token"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (q *Queries) GetPasswordResetByUserId(ctx context.Context, userID int32) (GetPasswordResetByUserIdRow, error) {
	row := q.db.QueryRowContext(ctx, getPasswordResetByUserId, userID)
	var i GetPasswordResetByUserIdRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ResetToken,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getPasswordResetByUserIdAndToken = `-- name: GetPasswordResetByUserIdAndToken :one
SELECT id, user_id, reset_token, expires_at, created_at
FROM password_resets
WHERE user_id = $1 AND reset_token = $2
`

type GetPasswordResetByUserIdAndTokenParams struct {
	UserID     int32  `json:"user_id"`
	ResetToken string `json:"reset_token"`
}

type GetPasswordResetByUserIdAndTokenRow struct {
	ID         int32     `json:"id"`
	UserID     int32     `json:"user_id"`
	ResetToken string    `json:"reset_token"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (q *Queries) GetPasswordResetByUserIdAndToken(ctx context.Context, arg GetPasswordResetByUserIdAndTokenParams) (GetPasswordResetByUserIdAndTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getPasswordResetByUserIdAndToken, arg.UserID, arg.ResetToken)
	var i GetPasswordResetByUserIdAndTokenRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ResetToken,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}
