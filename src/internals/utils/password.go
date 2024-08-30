package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	// Generating a hashed version of the password with a default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// CheckPassword compares the hashed password with the plaintext password
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
