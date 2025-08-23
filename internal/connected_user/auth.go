package connected_user

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

type AuthCredentials struct {
	Username string
	Password string
}

var ErrWrongPassword = errors.New("Wrong password")

type AuthRepository interface {
	CreateUserCredentials(username, hashedPassword string) error
	GetHashedPassword(username string) (string, error)
	DeleteCredentials(username string) error
}

type UserRepository interface {
	CreateUserProfile(username string) (User, error)
	GetUserByID(id int) (User, error)
	GetUserByUsername(username string) (User, error)
	UpdateUserProfile(user User) error
	DeleteUser(id int) error
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func Login(credentials AuthCredentials, authRepo AuthRepository, userRepo UserRepository) (User, error) {
	hashedPassword := HashPassword(credentials.Password)
	storedHash, err := authRepo.GetHashedPassword(credentials.Username)
	if err != nil {
		return User{}, fmt.Errorf("Failed to login, could't retrieve hash: %w", err)
	}
	if hashedPassword != storedHash {
		return User{}, ErrWrongPassword
	}

	user, err := userRepo.GetUserByUsername(credentials.Username)
	if err != nil {
		return User{}, fmt.Errorf("Failed to retrieve user by username %v got: %w", credentials.Username, err)
	}

	return user, nil
}
