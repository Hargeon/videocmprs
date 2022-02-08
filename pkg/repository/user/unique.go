package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

// Unique function check if user with email exists
func (repo *Repository) Unique(ctx context.Context, email string) (bool, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var total int64
	err := sq.Select("count(*) as total").
		From(TableName).
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&total)

	return total == 0, err
}
