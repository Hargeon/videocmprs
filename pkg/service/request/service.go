package request

import (
	"context"
	"github.com/google/jsonapi"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (srv *Service) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	// create video

	// create request

	// send video to cloud
	return nil, nil
}
