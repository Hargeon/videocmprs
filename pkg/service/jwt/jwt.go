// Package jwt uses for generation and parsing jwt tokens for user
package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

const tokenTD = 12 * time.Hour

type authClaims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

// SignedString function creates jwt token
func SignedString(id int64) (string, error) {
	claims := authClaims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTD).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseToken(tokenStr string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*authClaims); ok && token.Valid {
		return claims.Id, nil
	}
	return 0, errors.New("invalid token")
}
