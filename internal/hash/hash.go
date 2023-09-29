package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hasher struct {
	key string
}

func New(key string) Hasher {
	return Hasher{
		key: key,
	}
}

func (h Hasher) Sha256(value string) string {
	sum := sha256.Sum256([]byte(value + h.key))
	return hex.EncodeToString(sum[:])
}

func (h Hasher) IsValidSum(sum string, value string) bool {
	return sum == h.Sha256(value)
}

func Sign(hashKey string, body []byte) string {
	if hashKey == "" {
		return ""
	}

	hasher := New(hashKey)

	return hasher.Sha256(string(body))
}

func IsValidSum(sum string, value string, key string) bool {
	hasher := New(key)

	return hasher.IsValidSum(sum, value)
}
