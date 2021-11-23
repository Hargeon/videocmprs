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

	video := new(Resource)

	err := sq.
		Select("id", "name", "size", "bitrate", "resolution_x",
			"resolution_y", "ratio_x", "ratio_y", "service_id").
		From(TableName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRowContext(c).
		Scan(&video.ID, &video.Name, &video.Size, &video.BitrateDB, &video.ResolutionXDB,
			&video.ResolutionYDB, &video.RatioXDB, &video.RatioYDB, &video.ServiceID)

	if err != nil {
		return nil, err
	}

	if video.BitrateDB.Valid {
		video.Bitrate = video.BitrateDB.Int64
	}

	if video.ResolutionXDB.Valid {
		video.ResolutionX = int(video.ResolutionXDB.Int32)
	}

	if video.ResolutionYDB.Valid {
		video.ResolutionY = int(video.ResolutionYDB.Int32)
	}

	if video.RatioXDB.Valid {
		video.RatioX = int(video.RatioXDB.Int32)
	}

	if video.RatioYDB.Valid {
		video.RatioY = int(video.RatioYDB.Int32)
	}

	return video, err
}
