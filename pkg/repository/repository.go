// Package repository represent database connection
package repository

import (
	"github.com/Hargeon/videocmprs/db/model/user"
)

type UserRepository interface {
	Create(u *user.Resource) (*user.Resource, error)
}

type SessionRepository interface {
	Retrieve(u *user.Resource) (int64, error)
}
