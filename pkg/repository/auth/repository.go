// Package auth represent db connection to check if user exists
package auth

import (
	"context"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"time"
)

const queryTimeOut = 5 * time.Second

// Repository ...
type Repository struct {
	db *sqlx.DB
}

// NewRepository ...
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Exists function return id if user exists
func (repo *Repository) Exists(ctx context.Context, email, password string) (int64, error) {
	query, args, err := sq.Select("id").From(user.UserTableName).
		Where(sq.And{sq.Eq{"email": email}, sq.Eq{"password_hash": password}}).
		Limit(1).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err = repo.db.QueryRowxContext(c, query, args...).Scan(&id)
	return id, err
}
