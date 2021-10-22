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

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	video, ok := resource.(*Resource)
	if !ok {
		return nil, errors.New("invalid type assertion for *video.Resource in repository")
	}

	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.Insert(TableName).
		Columns("name", "size", "bitrate").
		Values(video.Name, video.Size, video.Bitrate).
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

func (r *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	query, args, err := sq.Select("id", "name", "size", "bitrate", "resolution", "ratio", "service_id").
		From(TableName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	video := new(Resource)
	err = r.db.QueryRowxContext(c, query, args...).StructScan(video)
	return video, err
}
