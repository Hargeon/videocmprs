package request

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

// Create request in db
func (repo *Repository) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	request, ok := resource.(*Resource)
	if !ok {
		return nil, errors.New("invalid type assertion *request.Resource in request repository")
	}

	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.Insert(TableName).
		Columns("bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "user_id", "video_name").
		Values(request.Bitrate, request.ResolutionX, request.ResolutionY, request.RatioX,
			request.RatioY, request.UserID, request.VideoName).
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
