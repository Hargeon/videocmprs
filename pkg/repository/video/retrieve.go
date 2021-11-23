package video

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

// Retrieve video from db
func (r *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	video := new(ResourceDTO)

	err := sq.
		Select("id", "name", "size", "bitrate", "resolution_x",
			"resolution_y", "ratio_x", "ratio_y", "service_id").
		From(TableName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRowContext(c).
		Scan(&video.ID, &video.Name, &video.Size, &video.Bitrate, &video.ResolutionX,
			&video.ResolutionY, &video.RatioX, &video.RatioY, &video.ServiceID)

	if err != nil {
		return nil, err
	}

	return video.BuildResource(), err
}
