package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

// Exists function return id if user exists
func (repo *Repository) Exists(ctx context.Context, email, password string) (int64, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.
		Select("id").
		From(TableName).
		Where(sq.And{sq.Eq{"email": email}, sq.Eq{"password_hash": password}}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&id)

	return id, err
}
