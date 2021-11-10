package request

import (
	"mime/multipart"

	"github.com/Hargeon/videocmprs/pkg/repository/video"

	"github.com/google/jsonapi"
)

// TableName is table name in db
const TableName = "requests"

var _ jsonapi.Linkable = (*Resource)(nil)

// Resource represent requests in db
type Resource struct {
	ID        int64 `jsonapi:"primary,requests"`
	UserID    int64
	VideoName string `jsonapi:"attr,video_name"`
	Status    string `jsonapi:"attr,status,omitempty"`
	Details   string `jsonapi:"attr,details,omitempty"`

	Bitrate     int64 `jsonapi:"attr,bitrate"`
	ResolutionX int   `jsonapi:"attr,resolution_x"`
	ResolutionY int   `jsonapi:"attr,resolution_y"`
	RatioX      int   `jsonapi:"attr,ratio_x"`
	RatioY      int   `jsonapi:"attr,ratio_y"`

	OriginalVideo  *video.Resource `jsonapi:"relation,original_video,omitempty"`
	ConvertedVideo *video.Resource `jsonapi:"relation,converted_video,omitempty"`

	VideoRequest *multipart.FileHeader
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
