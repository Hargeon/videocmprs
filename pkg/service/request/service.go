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
	// create request

	// send video to cloud

	// create video

	// update request
	return nil, nil
}
