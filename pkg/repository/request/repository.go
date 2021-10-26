// Package request represent db connection to creating and retrieving request
package request

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

// NewRepository ...
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Create request in db
func (repo *Repository) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	request, ok := resource.(*Resource)
	if !ok {
		return nil, errors.New("invalid type assertion *request.Resource in request repository")
	}

	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var id int64
	err := sq.Insert(TableName).
		Columns("bitrate", "resolution", "ration", "user_id").
		Values(request.Bitrate, request.Resolution, request.Ratio, request.UserID).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&id)
	if err != nil {
		return nil, err
	}

	return repo.Retrieve(ctx, id)
}

func (repo *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	request := new(Resource)
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	query, args, err := sq.
		Select("id", "status", "details", "bitrate", "resolution", "ratio").
		From(TableName).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = repo.db.GetContext(c, request, query, args...)
	return request, err
}

func (repo *Repository) Update(ctx context.Context, id int64, fields map[string]interface{}) (jsonapi.Linkable, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var reqId int64
	err := sq.
		Update(TableName).
		SetMap(fields).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id").
		RunWith(repo.db).
		PlaceholderFormat(sq.Dollar).
		QueryRowContext(c).
		Scan(&reqId)

	if err != nil {
		return nil, err
	}

	return repo.Retrieve(ctx, reqId)
}
