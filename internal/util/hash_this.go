package util

import (
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// HashThis returns the bcrypt hash of the code
func HashThis(code int) (string, error) {
	codeString := strconv.Itoa(code)

	hashedCode, err := bcrypt.GenerateFromPassword([]byte(codeString), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedCode), nil
}

// HashVerify checks if the provided code is correct or not
func HashVerify(code string, hashedCode string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(code))
}
