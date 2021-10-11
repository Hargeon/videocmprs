package middleware

import (
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type Session struct {
	service service.SessionService
}

func NewSessionMiddleware(s service.SessionService) func(c *fiber.Ctx) error {
	middleware := &Session{service: s}
	return middleware.ParseJWT
}

func (s *Session) ParseJWT(c *fiber.Ctx) error {
	fmt.Println("Parse jwt")
	return nil
}
