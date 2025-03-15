package security

import (
	"errors"
	"log"
	"os"
)

func getMasterKey() ([]byte, error) {
	key := os.Getenv("MASTER_KEY")
	if len(key) != 32 {
		log.Println(len(key))
		return nil, errors.New("MASTER_KEY must be 32 bytes long")
	}
	return []byte(key), nil
}
