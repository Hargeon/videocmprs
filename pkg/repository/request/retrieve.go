package request

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

// Retrieve request from db
func (repo *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	request := new(Resource)
	c, cancel := context.WithTimeout(ctx, queryTimeOut)

	defer cancel()

	err := sq.
		Select("id", "status", "details", "bitrate", "resolution_x", "resolution_y",
			"ratio_x", "ratio_y", "video_name").
		From(TableName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&request.ID, &request.Status, &request.DetailsDB, &request.Bitrate,
			&request.ResolutionX, &request.ResolutionY, &request.RatioX, &request.RatioY,
			&request.VideoName)
	if err != nil {
		return nil, err
	}

	request.Details = request.DetailsDB.String

	return request, err
}
