// Package service represent business logic
package service

import (
	"github.com/Hargeon/videocmprs/db/model"
	"github.com/Hargeon/videocmprs/pkg/repository"
)

// Authorization is abstraction for authorization logic
type Authorization interface {
	CreateUser(user *model.User) (int64, error)
	GenerateToken(email, password string) (string, error)
}

// Service is abstraction for business logic
type Service struct {
	Authorization
}

// NewService returns new Service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
	}
}
