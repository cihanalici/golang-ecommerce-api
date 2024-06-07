package token

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Maker interface {
	// CreateToken creates a new token for a specific userID and duration
	CreateToken(userID int32, duration time.Duration) (string, error)
	// VerifyToken verifies the token and returns the payload
	VerifyToken(token string) (*Payload, error)
}

func GenerateResetToken() (string, error) {
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
