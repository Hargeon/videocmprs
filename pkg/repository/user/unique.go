package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

// Unique function check if user with email exists
func (repo *Repository) Unique(ctx context.Context, email string) bool {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.Select("id").
		From(TableName).
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&id)

	return err != nil
}
