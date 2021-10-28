package request

import (
	"context"
	"errors"
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/google/jsonapi"
)

type Service struct {
	requestRepo  repository.UpdaterRepository
	videoRepo    repository.UpdaterRepository
	cloudStorage service.CloudStorage
}

func NewService(rRepo, vRepo repository.UpdaterRepository, cS service.CloudStorage) *Service {
	return &Service{
		requestRepo:  rRepo,
		videoRepo:    vRepo,
		cloudStorage: cS,
	}
}

func (srv *Service) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	res, ok := resource.(*request.Resource)
	if !ok {
		return nil, errors.New("invalid type assertion *request.Resource in service")
	}

	videoRes := res.OriginalVideo
	videoFile := res.VideoRequest

	linkable, err := srv.requestRepo.Create(ctx, resource)
	if err != nil {
		return nil, err
	}

	req, ok := linkable.(*request.Resource)
	if !ok {
		return nil, errors.New("invalid type assertion *request.Resource in service")
	}

	req.VideoRequest = videoFile
	srvVideoId, err := srv.cloudStorage.Upload(ctx, req.VideoRequest)
	if err != nil {
		fields := map[string]interface{}{"status": "failed", "details": `Can't upload video to cloud'`}
		reqLinkable, updateErr := srv.requestRepo.Update(ctx, req.ID, fields)
		if updateErr != nil {
			return nil, fmt.Errorf("can't upload video to cloud: %s, can't update request status: %s",
				err.Error(), updateErr.Error())
		}
		return reqLinkable, err
	}

	videoRes.ServiceId = srvVideoId
	videoLinkable, err := srv.videoRepo.Create(ctx, videoRes)
	if err != nil {
		fields := map[string]interface{}{"status": "failed", "details": `Can't add video to database'`}
		reqLinkable, updateErr := srv.requestRepo.Update(ctx, req.ID, fields)
		if updateErr != nil {
			return nil, fmt.Errorf("can't add video to database: %s, can't update request status: %s",
				err, updateErr)
		}
		return reqLinkable, err
	}

	updatedVideo, ok := videoLinkable.(*video.Resource)
	if !ok {
		return nil, errors.New("invalid type assertion *video.Resource in service after video update")
	}

	req.OriginalVideo = updatedVideo
	return req, nil
}
