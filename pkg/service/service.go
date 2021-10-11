// Package service represent business logic
package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
)

// Authorization is abstraction for authorization logic
type Authorization interface {
	GenerateToken(email, password string) (string, error)
}

type Retriever interface {
	Retrieve(c *fiber.Ctx) (jsonapi.Metable, error)
}

type Creator interface {
	Create(c *fiber.Ctx, resource jsonapi.Linkable) (jsonapi.Metable, error)
}

type SessionService interface {
	GenerateToken(c *fiber.Ctx) (jsonapi.Metable, error)
}

type UserService interface {
	Retriever
	Creator
}

type VideoService interface {
	Retriever
	Creator
}
