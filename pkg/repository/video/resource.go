package video

import (
	"github.com/google/jsonapi"
)

// TableName is name of table in db
const TableName = "videos"

// Resource represent video in db
type Resource struct {
	ID         int64  `jsonapi:"primary,videos"`
	Name       string `jsonapi:"attr,name"`
	Size       int64  `jsonapi:"attr,size"`
	Bitrate    int64  `jsonapi:"attr,bitrate,omitempty"`
	Resolution string `jsonapi:"attr,resolution,omitempty"`
	Ratio      string `jsonapi:"attr,ratio,omitempty"`
	ServiceId  string
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
