package video

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

// Update video TODO implement in future
func (r *Repository) Update(ctx context.Context, id int64, fields map[string]interface{}) (jsonapi.Linkable, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var reqID int64
	err := sq.
		Update(TableName).
		SetMap(fields).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar).
		QueryRowContext(c).
		Scan(&reqID)

	if err != nil {
		return nil, err
	}

	return r.Retrieve(ctx, reqID)
}
