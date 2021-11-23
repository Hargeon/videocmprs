package video

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

// Create video in db
func (r *Repository) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	video, ok := resource.(*Resource)
	if !ok {
		return nil, errors.New("invalid type assertion for *video.Resource in repository")
	}

	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.Insert(TableName).
		Columns("name", "size", "service_id").
		Values(video.Name, video.Size, video.ServiceID).
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
