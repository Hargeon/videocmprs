package repository

import (
	"github.com/Hargeon/videocmprs/db/model"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user *model.User) (int64, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
	}
}
