// Package user represent db connection to creating and retrieving user
package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
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

// Create user in db table users
func (repo *Repository) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	user, ok := resource.(*Resource)
	if !ok {
		return nil, errors.New("invalid type assertion in repository")
	}

	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.
		Insert(TableName).
		Columns("email", "password_hash").
		Values(user.Email, user.Password).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&id)

	if err != nil {
		return nil, err
	}

	return repo.Retrieve(ctx, id)
}

// Retrieve user by id
func (repo *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	user := new(Resource)
	err := sq.
		Select("id", "email").
		From(TableName).
		Where(sq.Eq{"id": id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&user.ID, &user.Email)

	if err != nil {
		return nil, err
	}

	return user, err
}
