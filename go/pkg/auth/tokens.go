package auth

import (
    "crypto/rand"
    "crypto/sha256"
    "crypto/subtle"
    "encoding/hex"
    "errors"
    "io"
)

// GenerateToken returns a URL-safe random token of n bytes, hex-encoded (length 2n).
func GenerateToken(n int) (string, error) {
    if n <= 0 { return "", errors.New("token: n must be > 0") }
    b := make([]byte, n)
    if _, err := io.ReadFull(rand.Reader, b); err != nil { return "", err }
    dst := make([]byte, hex.EncodedLen(len(b)))
    hex.Encode(dst, b)
    return string(dst), nil
}

// HashToken returns a hex-encoded SHA-256 hash of token.
func HashToken(token string) string {
    sum := sha256.Sum256([]byte(token))
    return hex.EncodeToString(sum[:])
}

// CompareTokenHash compares a hex-encoded hash to a token in constant time.
func CompareTokenHash(hashHex, token string) bool {
    if hashHex == "" || token == "" { return false }
    expected := HashToken(token)
    if len(expected) != len(hashHex) { return false }
    return subtle.ConstantTimeCompare([]byte(expected), []byte(hashHex)) == 1
}

