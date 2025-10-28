package auth

import "errors"

// HashPassword hashes a password. Default build returns error to avoid pulling x/crypto.
func HashPassword(password string) (string, error) {
    return "", errors.New("password hashing not enabled; build with -tags=secure")
}

// CheckPassword compares a plaintext password to a stored hash.
func CheckPassword(hash, password string) bool { return false }

