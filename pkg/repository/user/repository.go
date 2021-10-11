package user

import (
	"context"
	"fmt"
	"github.com/Hargeon/videocmprs/db/model/user"
	"github.com/google/jsonapi"
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

func (r *Repository) Create(u *user.Resource, resource jsonapi.Linkable) (*jsonapi.Linkable, error) {
	nUser := new(user.Resource)

	// TODO use current_timestamp()
	u.CreatedAt = time.Now()
	// TODO use builder https://github.com/Masterminds/squirrel
	query := fmt.Sprintf("INSERT INTO %s (email, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id, email", user.TableName)
	ctx, cancel := context.WithTimeout(context.Background(), timeOutQuery)
	defer cancel()

	row := r.db.QueryRowxContext(ctx, query, u.Email, u.Password, u.CreatedAt)
	err := row.StructScan(nUser)
	// TODO separate
	// insert for insert
	// use return Retrieve(...)
	return nUser, err
}
