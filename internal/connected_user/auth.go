package connected_user

import (
	"fmt"
)

func Login(credentials AuthCredentials, authRepo AuthRepository, userRepo UserRepository) (User, error) {
	hashedPassword := HashPassword(credentials.Password)
	storedHash, err := authRepo.GetHashedPassword(credentials.AuthUsername)
	if err != nil {
		return User{}, fmt.Errorf("Failed to login, could't retrieve hash: %w", err)
	}
	if hashedPassword != storedHash {
		return User{}, ErrWrongPassword
	}

	user, err := userRepo.GetUserByUsername(credentials.AuthUsername)
	if err != nil {
		return User{}, fmt.Errorf("Failed to retrieve user by username %v got: %w", credentials.AuthUsername, err)
	}

	return user, nil
}
