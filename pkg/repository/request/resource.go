package request

import (
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"mime/multipart"
	"regexp"
)

// TableName is table name in db
const TableName = "requests"

// Resource represent requests in db
type Resource struct {
	ID      int64 `jsonapi:"primary,requests"`
	UserID  int64
	Status  string `jsonapi:"attr,status,omitempty"`
	Details string `jsonapi:"attr,details,omitempty"`

	Bitrate    int64  `jsonapi:"attr,bitrate" validate:"required"`
	Resolution string `jsonapi:"attr,resolution" validate:"required,resolution"`
	Ratio      string `jsonapi:"attr,ratio" validate:"required,resolution"`

	OriginalVideo  *video.Resource
	ConvertedVideo *video.Resource

	VideoRequest *multipart.FileHeader
}

// ValidateResolution function validate resolution and ratio
func ValidateResolution(fl validator.FieldLevel) bool {
	res := fl.Field().String()

	re, err := regexp.Compile(`\A([1-9]+[0-9]*):([1-9]+[0-9]*)\z`)
	if err != nil {
		return false
	}
	return re.MatchString(res)
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
