package request

import (
	"context"
	"fmt"

	"github.com/Hargeon/videocmprs/pkg/repository/video"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

// Retrieve request from db
func (repo *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)

	defer cancel()

	request := new(Resource)
	origin := new(video.DTO)
	converted := new(video.DTO)

	err := sq.
		Select(fmt.Sprintf("%s.id", TableName),
			fmt.Sprintf("%s.user_id", TableName),
			fmt.Sprintf("%s.status", TableName),
			fmt.Sprintf("%s.details", TableName),
			fmt.Sprintf("%s.bitrate", TableName),
			fmt.Sprintf("%s.resolution_x", TableName),
			fmt.Sprintf("%s.resolution_y", TableName),
			fmt.Sprintf("%s.ratio_x", TableName),
			fmt.Sprintf("%s.ratio_y", TableName),
			fmt.Sprintf("%s.video_name", TableName),
			"origin_video.id",
			"origin_video.name",
			"origin_video.size",
			"origin_video.bitrate",
			"origin_video.resolution_x",
			"origin_video.resolution_y",
			"origin_video.ratio_x",
			"origin_video.ratio_y",
			"origin_video.service_id",
			"converted_video.id",
			"converted_video.name",
			"converted_video.size",
			"converted_video.bitrate",
			"converted_video.resolution_x",
			"converted_video.resolution_y",
			"converted_video.ratio_x",
			"converted_video.ratio_y",
			"converted_video.service_id").
		From(TableName).
		LeftJoin(fmt.Sprintf("%s AS origin_video ON %s.original_file_id = origin_video.id",
			video.TableName, TableName)).
		LeftJoin(fmt.Sprintf("%s AS converted_video ON %s.converted_file_id = converted_video.id",
			video.TableName, TableName)).
		Where(sq.Eq{fmt.Sprintf("%s.id", TableName): id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&request.ID, &request.UserID, &request.Status, &request.DetailsDB, &request.Bitrate,
			&request.ResolutionX, &request.ResolutionY, &request.RatioX, &request.RatioY,
			&request.VideoName, &origin.ID, &origin.Name, &origin.Size, &origin.Bitrate,
			&origin.ResolutionX, &origin.ResolutionY, &origin.RatioX, &origin.RatioY,
			&origin.ServiceID, &converted.ID, &converted.Name, &converted.Size,
			&converted.Bitrate, &converted.ResolutionX, &converted.ResolutionY,
			&converted.RatioX, &converted.RatioY, &converted.ServiceID)

	if err != nil {
		return nil, err
	}

	request.Details = request.DetailsDB.String
	// check if videos exists in db
	if origin.ID.Valid {
		request.OriginalVideo = origin.BuildResource()
	}

	if converted.ID.Valid {
		request.ConvertedVideo = converted.BuildResource()
	}

	return request, err
}
