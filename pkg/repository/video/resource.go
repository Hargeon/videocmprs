package video

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/google/jsonapi"
)

// TableName is name of table in db
const TableName = "videos"

var _ jsonapi.Linkable = (*Resource)(nil)

// Resource represent video in db
type Resource struct {
	ID          int64  `jsonapi:"primary,videos" json:"id,omitempty"`
	UserID      int64  `json:"user_id"`
	Name        string `jsonapi:"attr,name" json:"name"`
	Size        int64  `jsonapi:"attr,size" json:"size"`
	Bitrate     int64  `jsonapi:"attr,bitrate,omitempty" json:"bitrate"`
	ResolutionX int    `jsonapi:"attr,resolution_x,omitempty" json:"resolution_x"`
	ResolutionY int    `jsonapi:"attr,resolution_y,omitempty" json:"resolution_y"`
	RatioX      int    `jsonapi:"attr,ratio_x,omitempty" json:"ratio_x"`
	RatioY      int    `jsonapi:"attr,ratio_y,omitempty" json:"ratio_y"`
	ServiceID   string `json:"service_id,omitempty"`
}

type DTO struct {
	ID          sql.NullInt64
	UserID      sql.NullInt64
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
		UserID:      dto.UserID.Int64,
		Size:        dto.Size.Int64,
		Bitrate:     dto.Bitrate.Int64,
		ResolutionX: int(dto.ResolutionX.Int32),
		ResolutionY: int(dto.ResolutionY.Int32),
		RatioX:      int(dto.RatioX.Int32),
		RatioY:      int(dto.RatioY.Int32),
		ServiceID:   dto.ServiceID.String,
	}
}

// BuildFields function create map with fields and values for INSERT DB query
func (r *Resource) BuildFields() map[string]interface{} {
	fields := make(map[string]interface{})

	if r.UserID != 0 {
		fields["user_id"] = r.UserID
	}

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
		"self": fmt.Sprintf("%s/api/v1/videos/%d", os.Getenv("BASE_URL"), r.ID),
	}
}
