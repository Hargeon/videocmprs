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
	ID          int64  `jsonapi:"primary,videos" json:"id,omitempty"`
	Name        string `jsonapi:"attr,name" json:"name"`
	Size        int64  `jsonapi:"attr,size" json:"size"`
	Bitrate     int64  `jsonapi:"attr,bitrate,omitempty" json:"bitrate"`
	ResolutionX int    `jsonapi:"attr,resolution_x,omitempty" json:"resolution_x"`
	ResolutionY int    `jsonapi:"attr,resolution_y,omitempty" json:"resolution_y"`
	RatioX      int    `jsonapi:"attr,ratio_x,omitempty" json:"ratio_x"`
	RatioY      int    `jsonapi:"attr,ratio_y,omitempty" json:"ratio_y"`
	ServiceID   string `json:"service_id,omitempty"`

	IDDB          sql.NullInt64
	NameDB        sql.NullString
	SizeDB        sql.NullInt64
	BitrateDB     sql.NullInt64
	ResolutionXDB sql.NullInt32
	ResolutionYDB sql.NullInt32
	RatioXDB      sql.NullInt32
	RatioYDB      sql.NullInt32
	ServiceIDDB   sql.NullString
}

// BuildFields function create map with fields and values for INSERT DB query
func (r *Resource) BuildFields() map[string]interface{} {
	fields := make(map[string]interface{})

	if r.Name != "" {
		fields["name"] = r.Name
	}

	if r.Size != 0 {
		fields["size"] = r.Size
	}

	if r.Bitrate != 0 {
		fields["bitrate"] = r.Bitrate
	}

	if r.ResolutionX != 0 {
		fields["resolution_x"] = r.ResolutionX
	}

	if r.ResolutionY != 0 {
		fields["resolution_y"] = r.ResolutionY
	}

	if r.RatioX != 0 {
		fields["ratio_x"] = r.RatioX
	}

	if r.RatioY != 0 {
		fields["ratio_y"] = r.RatioY
	}

	if r.ServiceID != "" {
		fields["service_id"] = r.ServiceID
	}

	return fields
}

// JSONAPILinks ...
func (r *Resource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "need add", // TODO need add link
	}
}
