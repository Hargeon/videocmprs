package request

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"

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
	DetailsDB sql.NullString

	Bitrate     int64 `jsonapi:"attr,bitrate" validate:"required_if=ResolutionX 0 ResolutionY 0 RatioX 0 RatioY 0"`
	ResolutionX int   `jsonapi:"attr,resolution_x" validate:"required_if=Bitrate 0 RatioX 0 RatioY 0,required_with=ResolutionY"` //nolint:lll
	ResolutionY int   `jsonapi:"attr,resolution_y" validate:"required_if=Bitrate 0 RatioX 0 RatioY 0,required_with=ResolutionX"` //nolint:lll
	RatioX      int   `jsonapi:"attr,ratio_x" validate:"required_if=ResolutionX 0 ResolutionY 0 Bitrate 0,required_with=RatioY"` //nolint:lll
	RatioY      int   `jsonapi:"attr,ratio_y" validate:"required_if=ResolutionX 0 ResolutionY 0 Bitrate 0,required_with=RatioX"` //nolint:lll

	OriginalVideo  *video.Resource `jsonapi:"relation,original_video,omitempty"`
	ConvertedVideo *video.Resource `jsonapi:"relation,converted_video,omitempty"`

	VideoRequest *multipart.FileHeader
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("%s/api/v1/requests/%d", os.Getenv("BASE_URL"), r.ID),
	}
}
