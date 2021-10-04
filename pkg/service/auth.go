package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/Hargeon/videocmprs/db/model"
	"github.com/Hargeon/videocmprs/pkg/repository"
	"os"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (auth *AuthService) CreateUser(user *model.User) (int64, error) {
	secret := os.Getenv("DB_SECRET")
	passHash := generateHash([]byte(user.Password), []byte(secret))
	user.Password = fmt.Sprintf("%x", passHash)
	return auth.repo.CreateUser(user)
}

func generateHash(password, salt []byte) []byte {
	hash := sha1.New()
	hash.Write(password)
	return hash.Sum(salt)
}
