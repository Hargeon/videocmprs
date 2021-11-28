// Package user represent db connection to creating and retrieving user
package user

import (
	"database/sql"
	"time"
)

const queryTimeOut = 5 * time.Second

// Repository represent db connection for user table
type Repository struct {
	db *sql.DB
}

// NewRepository initialize Repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
