// Package repository represent database connection
package repository

import (
	"context"

	"github.com/Hargeon/videocmprs/api/query"

	"github.com/google/jsonapi"
)

type Creator interface {
	Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error)
}

type Retriever interface {
	Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error)
}

type Paginator interface {
	List(ctx context.Context, params *query.Params) ([]interface{}, error)
}

type RelationExistable interface {
	RelationExists(ctx context.Context, userID, relationID int64) (int64, error)
}

type Updater interface {
	Update(ctx context.Context, id int64, fields map[string]interface{}) (jsonapi.Linkable, error)
}

type CreatorRetriever interface {
	Creator
	Retriever
}

type UserRepository interface {
	Creator
	Retriever

	Exists(ctx context.Context, email, password string) (int64, error)
	Unique(ctx context.Context, email string) (bool, error)
}

type VideoRepository interface {
	Retriever
	Updater
	RelationExistable

	Create(ctx context.Context, fields map[string]interface{}) (jsonapi.Linkable, error)
}

type RequestRepository interface {
	CreatorRetriever
	Updater
	Paginator
	RelationExistable
}
