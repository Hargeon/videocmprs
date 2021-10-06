package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/Hargeon/videocmprs/db/model"
	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

const tokenTD = 12 * time.Hour

// AuthService implement authorization
type AuthService struct {
	repo *repository.Repository
}

// authClaims is custom claim for jwt token
type authClaims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

// NewAuthService returns new AuthService
func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser hashes user password and try to create user
func (auth *AuthService) CreateUser(user *model.User) (int64, error) {
	secret := os.Getenv("DB_SECRET")
	passHash := generateHash([]byte(user.Password), []byte(secret))
	user.Password = fmt.Sprintf("%x", passHash)
	return auth.repo.Authorization.CreateUser(user)
}

// GenerateToken try to find user and returns jwt token
func (auth *AuthService) GenerateToken(email, password string) (string, error) {
	id, err := auth.repo.Authorization.GetUser(email, password)
	if err != nil {
		return "", err
	}
	claims := authClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTD).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id: id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

// generateHash for password
func generateHash(password, salt []byte) []byte {
	hash := sha1.New()
	hash.Write(password)
	return hash.Sum(salt)
}
