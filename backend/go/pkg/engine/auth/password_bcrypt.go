//go:build secure
// +build secure

package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a password using bcrypt with default cost.
func HashPassword(password string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil { return "", err }
    return string(b), nil
}

// CheckPassword compares a bcrypt hash and plaintext password.
func CheckPassword(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

