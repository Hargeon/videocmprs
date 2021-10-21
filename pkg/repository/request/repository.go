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

	req := new(Resource)
	//err := sq.Insert(TableName).
	//	Columns("bitrate", "resolution", "ration").
	//	Values(request.Bitrate, request.Resolution, request.Ration).
	//	Suffix("RETURNING id, bitrate, resolution, ration").
	//	PlaceholderFormat(sq.Dollar).
	//	RunWith(repo.db).
	//	QueryRowContext(c).
	//	Scan(&req.ID, &req.Bitrate, &req.Resolution, &req.Ration)

	query, args, err := sq.Insert(TableName).
		Columns("bitrate", "resolution", "ration").
		Values(request.Bitrate, request.Resolution, request.Ration).
		Suffix("RETURNING id, bitrate, resolution, ration").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRowxContext(c, query, args...).StructScan(req)
	return req, err
}
