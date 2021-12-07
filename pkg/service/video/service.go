package video

import (
	"context"
	"errors"

	"github.com/Hargeon/videocmprs/pkg/repository"

	"github.com/google/jsonapi"
)

type Service struct {
	repo repository.VideoRepository
}

func NewService(repo repository.VideoRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Retrieve(ctx context.Context, userID, relationID int64) (jsonapi.Linkable, error) {
	id, err := s.repo.RelationExists(ctx, userID, relationID)
	if err != nil {
		return nil, err
	}

	if id == 0 {
		return nil, errors.New("video does not exists")
	}

	return s.repo.Retrieve(ctx, id)
}
