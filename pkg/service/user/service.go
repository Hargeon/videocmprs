package user

import (
	"errors"
	"fmt"
	"github.com/Hargeon/videocmprs/db/model/user"
	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
)

type Service struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Retrieve(c *fiber.Ctx) (jsonapi.Metable, error) {
	fmt.Println("User retrieve")
	return nil, nil
}

func (s *Service) Create(c *fiber.Ctx) (jsonapi.Metable, error) {
	fmt.Println("User create service")
	u, ok := c.Locals("user").(*user.Resource)
	if !ok {
		return nil, errors.New("can't make a type assertion")
	}

	hashedPass := encryption.GenerateHash([]byte(u.Password))
	u.Password = fmt.Sprintf("%x", hashedPass)
	return s.repo.Create(u)
}
