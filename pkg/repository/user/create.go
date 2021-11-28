package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

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
