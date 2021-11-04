package request

import (
	"mime/multipart"

	"github.com/google/jsonapi"

	"github.com/Hargeon/videocmprs/pkg/repository/video"
)

// TableName is table name in db
const TableName = "requests"

// Resource represent requests in db
type Resource struct {
	ID      int64 `jsonapi:"primary,requests"`
	UserID  int64
	Status  string `jsonapi:"attr,status,omitempty"`
	Details string `jsonapi:"attr,details,omitempty"`

	Bitrate     int64 `jsonapi:"attr,bitrate" validate:"required"`
	ResolutionX int   `jsonapi:"attr,resolution_x,omitempty"`
	ResolutionY int   `jsonapi:"attr,resolution_y,omitempty"`
	RatioX      int   `jsonapi:"attr,ratio_x,omitempty"`
	RatioY      int   `jsonapi:"attr,ratio_y,omitempty"`

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
