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

type DTO struct {
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

func (dto *DTO) BuildResource() *Resource {
	return &Resource{
		ID:          dto.ID.Int64,
		Name:        dto.Name.String,
		Size:        dto.Size.Int64,
		Bitrate:     dto.Bitrate.Int64,
		ResolutionX: int(dto.ResolutionX.Int32),
		ResolutionY: int(dto.ResolutionY.Int32),
		RatioX:      int(dto.RatioX.Int32),
		RatioY:      int(dto.RatioY.Int32),
		ServiceID:   dto.ServiceID.String,
	}
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
