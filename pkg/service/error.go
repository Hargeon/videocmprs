package service

import "errors"

var (
	ErrUserNotExists = errors.New("user is not exists")
	ErrAlreadyExists = errors.New("the user is already exists")

	ErrInvalidTypeAssertion = errors.New("invalid type assertion in service")
)
