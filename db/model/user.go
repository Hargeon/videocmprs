// Package model consists models for database
package model

import "time"

// UserTableName is name of users table in db
const UserTableName = "users"

// User represent users table in db
type User struct {
	Id                   int
	Email                string `json:"email" validate:"required,email,min=6,max=32"`
	Password             string `json:"password" validate:"required,min=6,max=250"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,max=250,eqfield=Password"`
	CreatedAt            time.Time
}
