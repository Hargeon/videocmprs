package request

import (
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/google/jsonapi"
)

// TableName is table name in db
const TableName = "requests"

// Resource represent requests in db
type Resource struct {
	ID      int64 `db:"id"`
	Status  string
	Details string

	Bitrate    int64  `db:"bitrate"`
	Resolution string `db:"resolution"`
	Ration     string `db:"ratio"`

	OriginalVideo  video.Resource
	ConvertedVideo video.Resource
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
