// Package auth represent db connection to check if user exists
package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/Hargeon/videocmprs/pkg/repository/user"

	sq "github.com/Masterminds/squirrel"
)

const queryTimeOut = 5 * time.Second

// AuthorizationRepository ...
type AuthorizationRepository struct {
	db *sql.DB
}

// NewRepository ...
func NewRepository(db *sql.DB) *AuthorizationRepository {
	return &AuthorizationRepository{db: db}
}

// Exists function return id if user exists
func (repo *AuthorizationRepository) Exists(ctx context.Context, email, password string) (int64, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.
		Select("id").
		From(user.TableName).
		Where(sq.And{sq.Eq{"email": email}, sq.Eq{"password_hash": password}}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&id)

	return id, err
}
