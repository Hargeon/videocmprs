package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

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
