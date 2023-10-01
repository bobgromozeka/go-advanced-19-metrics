package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hasher provides sha256 hashing functionality
type Hasher struct {
	key string
}

// New Create Hasher
func New(key string) Hasher {
	return Hasher{
		key: key,
	}
}

// Sha256 Makes hash string from hash parameter
func (h Hasher) Sha256(value string) string {
	sum := sha256.Sum256([]byte(value + h.key))
	return hex.EncodeToString(sum[:])
}

// IsValidSum checks if provided string is equal to sum after sha256 hashing
func (h Hasher) IsValidSum(sum string, value string) bool {
	return sum == h.Sha256(value)
}

// Sign Facade function for Hasher.Sha256
func Sign(hashKey string, body []byte) string {
	if hashKey == "" {
		return ""
	}

	hasher := New(hashKey)

	return hasher.Sha256(string(body))
}

// IsValidSum Facade function for Hasher.IsValidSum
func IsValidSum(sum string, value string, key string) bool {
	hasher := New(key)

	return hasher.IsValidSum(sum, value)
}
