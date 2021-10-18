package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"testing"
)

func TestSignedString(t *testing.T) {
	var expectedId int64 = 65

	tokenString, err := SignedString(expectedId)
	if err != nil {
		t.Fatalf("Unexpected error when signing string, error: %s\n", err.Error())
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("TOKEN_SECRET")), err
	})

	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err.Error())
	}

	claims, ok := parsedToken.Claims.(*authClaims)
	if !ok {
		t.Fatalf("Invalid type assertion for authClaims\n")
	}

	id := claims.Id
	if id != expectedId {
		t.Errorf("Invalid id, expected: %d, got: %d\n", expectedId, id)
	}
}
