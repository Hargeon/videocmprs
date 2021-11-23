package video

import (
	"database/sql"

	"github.com/google/jsonapi"
)

// TableName is name of table in db
const TableName = "videos"

var _ jsonapi.Linkable = (*Resource)(nil)

// Resource represent video in db
type Resource struct {
	ID          int64  `jsonapi:"primary,videos"`
	Name        string `jsonapi:"attr,name"`
	Size        int64  `jsonapi:"attr,size"`
	Bitrate     int64  `jsonapi:"attr,bitrate,omitempty"`
	ResolutionX int    `jsonapi:"attr,resolution_x,omitempty"`
	ResolutionY int    `jsonapi:"attr,resolution_y,omitempty"`
	RatioX      int    `jsonapi:"attr,ratio_x,omitempty"`
	RatioY      int    `jsonapi:"attr,ratio_y,omitempty"`
	ServiceID   string
}

type ResourceDTO struct {
	ID          sql.NullInt64
	Name        sql.NullString
	Size        sql.NullInt64
	Bitrate     sql.NullInt64
	ResolutionX sql.NullInt32
	ResolutionY sql.NullInt32
	RatioX      sql.NullInt32
	RatioY      sql.NullInt32
	ServiceID   sql.NullString
}

func (r *ResourceDTO) BuildResource() *Resource {
	res := new(Resource)

	if r.ID.Valid {
		res.ID = r.ID.Int64
	}

	if r.Name.Valid {
		res.Name = r.Name.String
	}

	if r.Size.Valid {
		res.Size = r.Size.Int64
	}

	if r.Bitrate.Valid {
		res.Bitrate = r.Bitrate.Int64
	}

	if r.ResolutionX.Valid {
		res.ResolutionX = int(r.ResolutionX.Int32)
	}

	if r.ResolutionY.Valid {
		res.ResolutionY = int(r.ResolutionY.Int32)
	}

	if r.RatioX.Valid {
		res.RatioX = int(r.RatioX.Int32)
	}

	if r.RatioY.Valid {
		res.RatioY = int(r.RatioY.Int32)
	}

	if r.ServiceID.Valid {
		res.ServiceID = r.ServiceID.String
	}

	return res
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
