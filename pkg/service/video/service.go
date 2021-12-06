package video

import (
	"context"

	"github.com/Hargeon/videocmprs/pkg/repository"

	"github.com/google/jsonapi"
)

type Service struct {
	repo repository.Retriever
}

func NewService(repo repository.Retriever) *Service {
	return &Service{repo: repo}
}

func (s *Service) Retrieve(ctx context.Context, id int64) (jsonapi.Linkable, error) {
	return s.repo.Retrieve(ctx, id)
}
