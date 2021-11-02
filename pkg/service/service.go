// Package service represent business logic
package service

import (
	"context"
	"mime/multipart"

	"github.com/Hargeon/videocmprs/api/query"

	"github.com/google/jsonapi"
)

type Creator interface {
	Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error)
}

type Retriever interface {
	Retrieve(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error)
}

type Paginator interface {
	List(ctx context.Context, params query.Params) ([]jsonapi.Linkable, error)
}

type Tokenable interface {
	GenerateToken(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error)
	Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error)
}

type CloudStorage interface {
	Upload(ctx context.Context, header *multipart.FileHeader) (string, error)
}

type Request interface {
	Creator
	Paginator
}
