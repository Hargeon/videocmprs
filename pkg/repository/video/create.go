package video

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

// Create video in db
func (r *Repository) Create(ctx context.Context, fields map[string]interface{}) (jsonapi.Linkable, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.Insert(TableName).
		SetMap(fields).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRowContext(c).
		Scan(&id)

	if err != nil {
		return nil, err
	}

	return r.Retrieve(ctx, id)
}
