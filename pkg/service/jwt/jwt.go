package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

const tokenTD = 12 * time.Hour

type authClaims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

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
