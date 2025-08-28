package connected_user

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type AuthCredentials struct {
	AuthUsername string `json:"AuthUsername"`
	Password     string `json:"password"`
}

var ErrWrongPassword = errors.New("Wrong password")

type AuthRepository interface {
	CreateUserCredentials(username, hashedPassword string) error
	GetHashedPassword(username string) (string, error)
	DeleteCredentials(username string) error
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
