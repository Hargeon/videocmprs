package user

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

// NewRepository ...
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Create user in db table users
func (repo *Repository) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	user, ok := resource.(*Resource)
	if !ok {
		return nil, errors.New("invalid type assertion in repository")
	}

	query, args, err := sq.Insert(UserTableName).Columns("email", "password_hash").
		Values(user.Email, user.Password).Suffix("RETURNING id").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var id int64
	err = repo.db.QueryRowxContext(c, query, args...).Scan(&id)
	if err != nil {
		return nil, err
	}

	return repo.Retrieve(ctx, id)
}

// Retrieve user by id
func (repo *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	query, args, err := sq.Select("id", "email").From(UserTableName).
		Where(sq.Eq{"id": id}).Limit(1).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := new(Resource)
	err = repo.db.QueryRowxContext(c, query, args...).StructScan(user)
	return user, err
}
