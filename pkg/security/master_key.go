package security

import (
	"errors"
	"os"
)

// getMasterKey retrieves the master key from the environment variable `MASTER_KEY`.
// The key must be exactly 32 bytes long to ensure compatibility with cryptographic algorithms.
// If the key is missing or has an incorrect length, an error is returned.
func getMasterKey() ([]byte, error) {
	key := os.Getenv("MASTER_KEY")
	if len(key) != 32 {
		return nil, errors.New("MASTER_KEY must be 32 bytes long")
	}
	return []byte(key), nil
}
