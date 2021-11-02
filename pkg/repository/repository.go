// Package repository represent database connection
package repository

import (
	"context"
	"github.com/google/jsonapi"
)

type Creator interface {
	Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error)
}

type Retriever interface {
	Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error)
}

type Existable interface {
	Retriever
	Exists(ctx context.Context, email, password string) (int64, error)
}

type Updater interface {
	Update(ctx context.Context, id int64, fields map[string]interface{}) (jsonapi.Linkable, error)
}

type Repository interface {
	Creator
	Retriever
}

type UpdaterRepository interface {
	Repository
	Updater
}

// change this
type RequestRepository interface {
	UpdaterRepository
	RetrieveList(ctx context.Context, relationId int64, page int64) (jsonapi.Linkable, error)
}
