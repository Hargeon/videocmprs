// Package user represent entity from users db table
package user

import (
	"github.com/google/jsonapi"
	"time"
)

// TableName is name of users table in db
const TableName = "users"

// Resource represent users table in db
type Resource struct {
	Id                   int    `jsonapi:"primary,users" db:"id"`
	Email                string `json:"email" validate:"required,email,min=6,max=32" jsonapi:"attr,email" db:"email"`
	Password             string `json:"password" validate:"required,min=6,max=250"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,max=250,eqfield=Password"`
	CreatedAt            time.Time
}

func (r *Resource) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"details": "users meta information",
	}
}
