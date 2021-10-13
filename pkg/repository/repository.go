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

type UserRepository interface {
	Creator
	Retriever
}
