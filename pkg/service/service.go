// Package service represent business logic
package service

import (
	"context"
	"github.com/google/jsonapi"
	"mime/multipart"
)

type Creator interface {
	Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error)
}

type Retriever interface {
	Retrieve(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error)
}

type Tokenable interface {
	GenerateToken(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error)
}

type UserService interface {
	Creator
	Retriever
}

type CloudStorage interface {
	Upload(ctx context.Context, header *multipart.FileHeader) (string, error)
}
