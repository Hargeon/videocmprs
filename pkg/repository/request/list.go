package request

import (
	"context"
	"fmt"

	"github.com/Hargeon/videocmprs/api/query"
	"github.com/Hargeon/videocmprs/pkg/repository/video"

	sq "github.com/Masterminds/squirrel"
)

// List returns requests
func (repo *Repository) List(ctx context.Context, params *query.Params) ([]interface{}, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	requests := make([]interface{}, 0, params.PageSize)

	rows, err := sq.
		Select(fmt.Sprintf("%s.id", TableName),
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
		Where(sq.Eq{fmt.Sprintf("%s.user_id", TableName): params.RelationID}).
		OrderBy(fmt.Sprintf("%s.created_at DESC", TableName)).
		Limit(params.PageSize).
		Offset(params.PageNumber).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryContext(c)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		request := new(Resource)
		origin := new(video.Resource)
		converted := new(video.Resource)

		err = rows.Scan(&request.ID, &request.Status, &request.Details, &request.Bitrate,
			&request.ResolutionX, &request.ResolutionY, &request.RatioX, &request.RatioY,
			&request.VideoName, &origin.IDDB, &origin.NameDB, &origin.SizeDB, &origin.BitrateDB,
			&origin.ResolutionXDB, &origin.ResolutionYDB, &origin.RatioXDB, &origin.RatioYDB,
			&origin.ServiceIDDB, &converted.IDDB, &converted.NameDB, &converted.SizeDB,
			&converted.BitrateDB, &converted.ResolutionXDB, &converted.ResolutionYDB,
			&converted.RatioXDB, &converted.RatioYDB, &converted.ServiceIDDB)

		if err != nil {
			return nil, err
		}

		// unmarshal origin video
		repo.unmarshalDBVideo(origin)
		// unmarshal converted video
		repo.unmarshalDBVideo(converted)

		// check if videos exists in db
		if origin.ID > 0 {
			request.OriginalVideo = origin
		}

		if converted.ID > 0 {
			request.ConvertedVideo = converted
		}

		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

// unmarshalDBVideo replace values from sql.Null fields with primitive types
func (repo *Repository) unmarshalDBVideo(resource *video.Resource) {
	if resource.IDDB.Valid {
		resource.ID = resource.IDDB.Int64
	}

	if resource.NameDB.Valid {
		resource.Name = resource.NameDB.String
	}

	if resource.SizeDB.Valid {
		resource.Size = resource.SizeDB.Int64
	}

	if resource.BitrateDB.Valid {
		resource.Bitrate = resource.BitrateDB.Int64
	}

	if resource.ResolutionXDB.Valid {
		resource.ResolutionX = int(resource.ResolutionXDB.Int32)
	}

	if resource.ResolutionYDB.Valid {
		resource.ResolutionY = int(resource.ResolutionYDB.Int32)
	}

	if resource.RatioXDB.Valid {
		resource.RatioX = int(resource.RatioXDB.Int32)
	}

	if resource.RatioYDB.Valid {
		resource.RatioY = int(resource.RatioYDB.Int32)
	}

	if resource.ServiceIDDB.Valid {
		resource.ServiceID = resource.ServiceIDDB.String
	}
}
