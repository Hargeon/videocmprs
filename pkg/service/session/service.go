package session

import (
	"errors"
	"fmt"
	"github.com/Hargeon/videocmprs/db/model/user"
	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"os"
	"time"
)

const tokenTD = 12 * time.Hour

type Service struct {
	repo repository.SessionRepository
}

// authClaims is custom claim for jwt token
type authClaims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

func NewService(repo repository.SessionRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GenerateToken(c *fiber.Ctx) (jsonapi.Metable, error) {
	usr, ok := c.Locals("user").(*user.Resource)
	if !ok {
		return nil, errors.New("can't make a type assertion")
	}

	hashPassword := encryption.GenerateHash([]byte(usr.Password))
	usr.Password = fmt.Sprintf("%x", hashPassword)
	id, err := s.repo.Retrieve(usr)
	if err != nil {
		return nil, err
	}

	claims := authClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTD).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id: id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	resource := new(Resource)
	resource.Id = id
	resource.Token = signedToken
	return resource, nil
}
