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

type ExistAble interface {
	Exists(ctx context.Context, email, password string) (int64, error)
}

type Repository interface {
	Creator
	Retriever
}

type VideoRepository interface {
	Repository
}
