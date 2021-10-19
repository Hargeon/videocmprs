package video

import (
	"context"
	"github.com/google/jsonapi"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	return nil, nil
}

func (r *Repository) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	return nil, nil
}
