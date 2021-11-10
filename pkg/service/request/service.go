// Package request uses for creating user request
package request

import (
	"context"
	"errors"
	"fmt"

	"github.com/Hargeon/videocmprs/api/query"
	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service"

	"github.com/google/jsonapi"
)

// Service for adding and changing requests
type Service struct {
	requestRepo  repository.RequestRepository
	videoRepo    repository.VideoRepository
	cloudStorage service.CloudStorage
}

// NewService initialize Service
func NewService(rRepo repository.RequestRepository, vRepo repository.VideoRepository, cS service.CloudStorage) *Service {
	return &Service{
		requestRepo:  rRepo,
		videoRepo:    vRepo,
		cloudStorage: cS,
	}
}

// Create function creates request in db, uploads video to cloud, creates video in db
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

	// upload video to cloud
	req.VideoRequest = videoFile
	srvVideoID, err := srv.cloudStorage.Upload(ctx, req.VideoRequest)

	if err != nil {
		// update request status
		fields := map[string]interface{}{"status": "failed", "details": "Can't upload video to cloud"}
		_, updateErr := srv.requestRepo.Update(ctx, req.ID, fields)

		if updateErr != nil {
			return nil, fmt.Errorf("can't upload video to cloud: %s, can't update request status: %s",
				err.Error(), updateErr.Error())
		}

		return nil, err
	}

	// create video in db
	videoRes.ServiceID = srvVideoID
	videoLinkable, err := srv.videoRepo.Create(ctx, videoRes)

	if err != nil {
		// update request status
		fields := map[string]interface{}{"status": "failed", "details": `Can't add video to database`}
		_, updateErr := srv.requestRepo.Update(ctx, req.ID, fields)

		if updateErr != nil {
			return nil, fmt.Errorf("can't add video to database: %s, can't update request status: %s",
				err, updateErr)
		}

		return nil, err
	}

	updatedVideo, ok := videoLinkable.(*video.Resource)
	if !ok {
		return nil, errors.New("invalid type assertion *video.Resource in service after video update")
	}

	// add original_file_id in request
	fields := map[string]interface{}{"original_file_id": updatedVideo.ID}
	linkable, err = srv.requestRepo.Update(ctx, req.ID, fields)

	if err != nil {
		return nil, err
	}

	req, ok = linkable.(*request.Resource)
	if !ok {
		return nil, errors.New("invalid type assertion")
	}

	req.OriginalVideo = updatedVideo

	return req, nil
}

// List returns []*request.Resource
func (srv *Service) List(ctx context.Context, params *query.Params) ([]interface{}, error) {
	requests, err := srv.requestRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return requests, nil
}
