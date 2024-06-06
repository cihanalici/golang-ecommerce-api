package util

import "errors"

var (
	ErrUnauthorized          = errors.New("unauthorized")
	ErrNoAuthorizationHeader = errors.New("no authorization header provided")
	ErrInvalidAuthFormat     = errors.New("invalid authorization header format")
)
