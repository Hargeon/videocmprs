package encryption

import (
	"crypto/sha1"
	"os"
)

func GenerateHash(text []byte) []byte {
	secret := os.Getenv("SECRET")
	hash := sha1.New()
	hash.Write(text)

	return hash.Sum([]byte(secret))
}
