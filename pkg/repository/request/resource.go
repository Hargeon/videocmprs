package request

import (
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/google/jsonapi"
	"mime/multipart"
)

// TableName is table name in db
const TableName = "requests"

// Resource represent requests in db
type Resource struct {
	ID      int64  `db:"id"`
	UserID  int64  `db:"user_id"`
	Status  string `db:"status"`
	Details string `db:"details"`

	Bitrate    int64  `db:"bitrate"`
	Resolution string `db:"resolution"`
	Ratio      string `db:"ratio"`

	OriginalVideo  *video.Resource
	ConvertedVideo *video.Resource

	VideoRequest *multipart.FileHeader
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
