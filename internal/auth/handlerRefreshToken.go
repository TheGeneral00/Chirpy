package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
        key := make([]byte, 32)
        _, err := rand.Read(key)
        if err != nil {
                return "", fmt.Errorf("Failed to generate random key")
        }
        return hex.EncodeToString(key), nil
}
