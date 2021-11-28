package request

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

// Update request in db
func (repo *Repository) Update(ctx context.Context, id int64, fields map[string]interface{}) (jsonapi.Linkable, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var reqID int64
	err := sq.
		Update(TableName).
		SetMap(fields).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id").
		RunWith(repo.db).
		PlaceholderFormat(sq.Dollar).
		QueryRowContext(c).
		Scan(&reqID)

	if err != nil {
		return nil, err
	}

	return repo.Retrieve(ctx, reqID)
}
