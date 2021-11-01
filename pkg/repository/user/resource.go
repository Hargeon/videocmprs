package user

import (
	"github.com/google/jsonapi"
	"time"
)

// TableName is name of users table in db
const TableName = "users"

// Resource represent users table in db
type Resource struct {
	ID                   int64  `jsonapi:"primary,users" db:"id"`
	Email                string `jsonapi:"attr,email" db:"email" validate:"required,email,min=6,max=32"`
	Password             string `jsonapi:"attr,password,omitempty" validate:"required,min=6,max=250"`
	PasswordConfirmation string `jsonapi:"attr,password_confirmation,omitempty" validate:"required,min=6,max=250,eqfield=Password"`
	Token                string `jsonapi:"attr,token,omitempty"`
	CreatedAt            time.Time
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
