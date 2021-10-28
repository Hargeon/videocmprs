// Package video represent db connection to creating and retrieving video
package video

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
	"github.com/jmoiron/sqlx"
	"time"
)

const queryTimeOut = 5 * time.Second

// Repository ...
type Repository struct {
	db *sqlx.DB
}

// NewRepository initialize repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Create video in db
func (r *Repository) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	video, ok := resource.(*Resource)
	if !ok {
		return nil, errors.New("invalid type assertion for *video.Resource in repository")
	}

	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.Insert(TableName).
		Columns("name", "size", "service_id").
		Values(video.Name, video.Size, video.ServiceId).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRowContext(c).
		Scan(&id)

	if err != nil {
		return nil, err
	}

	return r.Retrieve(ctx, id)
}

// Retrieve video from db
func (r *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	video := new(Resource)

	err := sq.
		Select("id", "name", "size", "bitrate", "resolution", "ratio", "service_id").
		From(TableName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRowContext(c).
		Scan(&video.ID, &video.Name, &video.Size, &video.Bitrate, &video.Resolution,
			&video.Ratio, &video.ServiceId)

	return video, err
}

// Update video TODO implement in future
func (r *Repository) Update(ctx context.Context, id int64, fields map[string]interface{}) (jsonapi.Linkable, error) {
	return nil, nil
}
