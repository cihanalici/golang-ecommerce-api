package token

import "time"

type Maker interface {
	// CreateToken creates a new token for a specific userID and duration
	CreateToken(userID int32, duration time.Duration) (string, error)
	// VerifyToken verifies the token and returns the payload
	VerifyToken(token string) (*Payload, error)
}
