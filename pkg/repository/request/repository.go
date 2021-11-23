// Package request represent db connection to creating and retrieving request
package request

import (
	"database/sql"
	"time"
)

const queryTimeOut = 5 * time.Second

// Repository represent db connection for request table
type Repository struct {
	db *sql.DB
}

// NewRepository initialize Repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
