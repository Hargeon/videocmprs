package video

import "github.com/google/jsonapi"

const TableName = "videos"

type Resource struct {
	ID         int64  `jsonapi:"primary,videos" db:"id"`
	Name       string `jsonapi:"attr,name" db:"name"`
	Size       int64  `jsonapi:"attr,size" db:"size"`
	Bitrate    int64  `jsonapi:"attr,bitrate,omitempty" db:"bitrate"`
	Resolution string `jsonapi:"attr,resolution,omitempty" db:"resolution"`
	Ratio      string `jsonapi:"attr,ratio,omitempty" db:"ratio"`
	ServiceId  string `db:"service_id"`
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
