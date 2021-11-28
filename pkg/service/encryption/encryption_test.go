package encryption

import (
	"crypto/sha1"
	"os"
	"testing"
)

func TestEncryption(t *testing.T) {
	text := []byte("oiHIUNPbobnp;inhobLNNOIN LK:Nojion p;mnpinbOJNINHOBON")

	secret := os.Getenv("SECRET")
	hash := sha1.New()
	hash.Write(text)
	expectedText := hash.Sum([]byte(secret))

	if gotText := GenerateHash(text); string(expectedText) != string(gotText) {
		t.Errorf("Invalid encryption, expected: %s, got: %s\n", expectedText, gotText)
	}
}
