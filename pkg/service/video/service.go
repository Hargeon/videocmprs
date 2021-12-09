package video

import (
	"context"

	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service"

	"github.com/google/jsonapi"
)

type Service struct {
	repo  repository.VideoRepository
	cloud service.CloudStorage
}

func NewService(repo repository.VideoRepository, cloud service.CloudStorage) *Service {
	return &Service{repo: repo, cloud: cloud}
}

// Retrieve video by userID and videoID
func (s *Service) Retrieve(ctx context.Context, userID, relationID int64) (jsonapi.Linkable, error) {
	id, err := s.repo.RelationExists(ctx, userID, relationID)
	if err != nil {
		return nil, err
	}

	if id == 0 {
		return nil, ErrVideoNotPresent
	}

	return s.repo.Retrieve(ctx, id)
}

// DownloadURL returns url for downloading video from cloud
func (s *Service) DownloadURL(ctx context.Context, userID, videoID int64) (string, error) {
	v, err := s.Retrieve(ctx, userID, videoID)
	if err != nil {
		return "", err
	}

	vid, ok := v.(*video.Resource)
	if !ok {
		return "", ErrVideoAssertion
	}

	if vid.ServiceID == "" {
		return "", ErrVideoNotInCloud
	}

	return s.cloud.URL(vid.ServiceID)
}
