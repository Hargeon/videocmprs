package user

import (
	"context"
	"fmt"
	"github.com/Hargeon/videocmprs/db/model/user"
	"github.com/jmoiron/sqlx"
	"time"
)

const timeOutQuery = 5 * time.Second

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(u *user.Resource) (*user.Resource, error) {
	nUser := new(user.Resource)

	u.CreatedAt = time.Now()
	query := fmt.Sprintf("INSERT INTO %s (email, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id, email", user.TableName)
	ctx, cancel := context.WithTimeout(context.Background(), timeOutQuery)
	defer cancel()

	row := r.db.QueryRowxContext(ctx, query, u.Email, u.Password, u.CreatedAt)
	err := row.StructScan(nUser)

	return nUser, err
}
