// Package request represent db connection to creating and retrieving request
package request

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Hargeon/videocmprs/api/query"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/jsonapi"
)

const queryTimeOut = 5 * time.Second

// Repository represent db connection for request table
type Repository struct {
	db *sql.DB
}

// NewRepository initialize Repository
func NewRepository(db *sql.DB) *Repository {
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
		Columns("bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "user_id").
		Values(request.Bitrate, request.ResolutionX, request.ResolutionY, request.RatioX,
			request.RatioY, request.UserID).
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

// Retrieve request from db
func (repo *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	request := new(Resource)
	c, cancel := context.WithTimeout(ctx, queryTimeOut)

	defer cancel()

	err := sq.
		Select("id", "status", "details", "bitrate", "resolution_x", "resolution_y",
			"ratio_x", "ratio_y").
		From(TableName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRowContext(c).
		Scan(&request.ID, &request.Status, &request.Details, &request.Bitrate,
			&request.ResolutionX, &request.ResolutionY, &request.RatioX, &request.RatioY)
	if err != nil {
		return nil, err
	}

	return request, err
}

// Update request in db
func (repo *Repository) Update(ctx context.Context, id int64, fields map[string]interface{}) (jsonapi.Linkable, error) {
	c, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	var reqID int64
	err := sq.
		Update(TableName).
		SetMap(fields).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id").
		RunWith(repo.db).
		PlaceholderFormat(sq.Dollar).
		QueryRowContext(c).
		Scan(&reqID)

	if err != nil {
		return nil, err
	}

	return repo.Retrieve(ctx, reqID)
}

func (repo *Repository) List(ctx context.Context, params *query.Params) ([]jsonapi.Linkable, error) {
	//c, cancel := context.WithTimeout(ctx, queryTimeOut)
	//defer cancel()
	//
	//requests := make([]*Resource, 0, params.PageSize)
	//
	//rows, err := sq.
	//	Select("id", "status", "details", "bitrate", "resolution_x", "resolution_y",
	//	"ratio_x", "ratio_y").
	//	From(TableName).
	//	Where(sq.Eq{"user_id": params.RelationID}).
	//	OrderBy("created_at DESC").
	//	PlaceholderFormat(sq.Dollar).
	//	Limit(params.PageSize).
	//	Offset(params.PageNumber).
	//	RunWith(repo.db).
	//	QueryContext(c)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//defer rows.Close()
	//
	//for rows.Next() {
	//	request := new(Resource)
	//	err = rows.Scan(&request.ID, &request.Status, &request.Details, &request.Bitrate,
	//		&request.ResolutionX, &request.ResolutionY, &request.RatioX, &request.RatioY)
	//	if err != nil {
	//		return nil, err
	//	}
	//	requests = append(requests, request)
	//}
	//return requests, nil
	return nil, nil
}
