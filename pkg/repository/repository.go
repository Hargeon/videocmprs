// Package repository represent database connection
package repository

import (
	"github.com/Hargeon/videocmprs/db/model"
	"github.com/jmoiron/sqlx"
)

// Authorization is abstraction for users manipulation
type Authorization interface {
	CreateUser(user *model.User) (int64, error)
	GetUser(email, password string) (int64, error)
}

// Repository represent abstraction for database connection
type Repository struct {
	Authorization
}

// NewRepository return new Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
	}
}
