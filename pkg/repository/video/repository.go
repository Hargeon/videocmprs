// Package video represent db connection to creating and retrieving video
package video

import (
	"database/sql"
	"time"
)

const queryTimeOut = 5 * time.Second

// Repository represent db connection for video table
type Repository struct {
	db *sql.DB
}

// NewRepository initialize repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
