package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"testing"
	"time"
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

func TestParseToken(t *testing.T) {
	cases := []struct {
		name     string
		id       int64
		timeFrom time.Time
		secret   []byte

		expectedId   int64
		errorPresent bool
	}{
		{
			name:     "With invalid secret method",
			id:       65,
			secret:   []byte("invalid secret"),
			timeFrom: time.Now(),

			expectedId:   0,
			errorPresent: true,
		},
		{
			name:     "With expired token",
			id:       65,
			secret:   []byte(os.Getenv("TOKEN_SECRET")),
			timeFrom: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),

			expectedId:   0,
			errorPresent: true,
		},
		{
			name:     "With valid paras",
			id:       65,
			secret:   []byte(os.Getenv("TOKEN_SECRET")),
			timeFrom: time.Now(),

			expectedId:   65,
			errorPresent: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			claims := authClaims{
				Id: testCase.id,
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  testCase.timeFrom.Unix(),
					ExpiresAt: testCase.timeFrom.Add(tokenTD).Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, err := token.SignedString(testCase.secret)
			if err != nil {
				t.Fatalf("Unexpected error: %s\n", err.Error())
			}

			id, err := ParseToken(tokenStr)

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if id != testCase.expectedId {
				t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedId, id)
			}
		})
	}
}
