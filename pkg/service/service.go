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
	Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error)
}

type RetrieveRelation interface {
	Retrieve(ctx context.Context, userID, relationID int64) (jsonapi.Linkable, error)
}

type Paginator interface {
	List(ctx context.Context, params *query.Params) ([]interface{}, error)
}

type Tokenable interface {
	GenerateToken(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error)
	Retriever
}

type CloudStorage interface {
	Upload(ctx context.Context, header *multipart.FileHeader) (string, error)
	URL(filename string) (string, error)
}

type Request interface {
	Creator
	RetrieveRelation
	Paginator
}

type Video interface {
	RetrieveRelation

	DownloadURL(ctx context.Context, userID, videoID int64) (string, error)
}

type Publisher interface {
	Publish(body []byte) error
	Ping() error
}
