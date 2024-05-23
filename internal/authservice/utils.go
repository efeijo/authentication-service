package authservice

import (
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(encrypted), err
}

// validatePassword compares the hashed password with the password given, hashing the password given is done by the library
func validatePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
