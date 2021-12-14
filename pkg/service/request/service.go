// Package request uses for creating user request
package request

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"

	"github.com/Hargeon/videocmprs/api/query"
	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/Hargeon/videocmprs/pkg/service/compress"

	"github.com/google/jsonapi"
	"go.uber.org/zap"
)

// Service for adding and changing requests
type Service struct {
	requestRepo  repository.RequestRepository
	videoRepo    repository.VideoRepository
	cloudStorage service.CloudStorage
	publisher    service.Publisher
	logger       *zap.Logger
}

// NewService initialize Service
func NewService(rRepo repository.RequestRepository, vRepo repository.VideoRepository, cS service.CloudStorage, pb service.Publisher, logger *zap.Logger) *Service {
	return &Service{
		requestRepo:  rRepo,
		videoRepo:    vRepo,
		cloudStorage: cS,
		publisher:    pb,
		logger:       logger,
	}
}

// Create function creates request in db, uploads video to cloud, creates video in db
func (srv *Service) Create(ctx context.Context, resource jsonapi.Linkable) (jsonapi.Linkable, error) {
	res, ok := resource.(*request.Resource)
	if !ok {
		return nil, errors.New("invalid type assertion *request.Resource in service")
	}

	vid := res.OriginalVideo
	videoFile := res.VideoRequest

	linkable, err := srv.requestRepo.Create(ctx, resource)
	if err != nil {
		return nil, err
	}

	req, ok := linkable.(*request.Resource)
	if !ok {
		return nil, errors.New("invalid type assertion *request.Resource in service")
	}

	go srv.addVideo(ctx, *req, *vid, *videoFile)

	return req, nil
}

// addVideo to cloud and db
func (srv *Service) addVideo(ctx context.Context, req request.Resource, vid video.Resource, videoFile multipart.FileHeader) {
	cloudVideoID, err := srv.cloudStorage.Upload(ctx, &videoFile)
	if err != nil {
		srv.logger.Error("can't upload video to cloud", zap.Error(err))
		// update request status
		fields := map[string]interface{}{"status": "failed", "details": "Can't upload video to cloud"}
		_, updateErr := srv.requestRepo.Update(ctx, req.ID, fields)

		if updateErr != nil {
			srv.logger.Error("can't update request status", zap.Error(updateErr))
		}

		return
	}

	// create video in db
	vid.ServiceID = cloudVideoID
	videoLinkable, err := srv.videoRepo.Create(ctx, vid.BuildFields())

	if err != nil {
		srv.logger.Error("can't add video to database", zap.Error(err))
		// update request status
		fields := map[string]interface{}{"status": "failed", "details": `Can't add video to database`}
		_, updateErr := srv.requestRepo.Update(ctx, req.ID, fields)

		if updateErr != nil {
			srv.logger.Error("can't update request status", zap.Error(err))
		}

		return
	}

	createdVideo, ok := videoLinkable.(*video.Resource)
	if !ok {
		srv.logger.Error("invalid type assertion *video.Resource in service after creating video")

		return
	}

	// add original_file_id in request
	fields := map[string]interface{}{"original_file_id": createdVideo.ID}
	requestLinkable, err := srv.requestRepo.Update(ctx, req.ID, fields)

	if err != nil {
		srv.logger.Error("Can't add original_file_id to request", zap.Error(err))

		return
	}

	updatedReq, ok := requestLinkable.(*request.Resource)
	if !ok {
		srv.logger.Error("invalid type assertion for *request.Resource after adding original_file_id")

		return
	}

	// send requests to rabbit
	err = srv.rabbitPublish(updatedReq)
	if err != nil {
		srv.logger.Error("Can't add request to rabbit", zap.Error(err))

		fields := map[string]interface{}{"status": "failed", "details": `Failed connection to worker`}

		_, updateErr := srv.requestRepo.Update(ctx, req.ID, fields)

		if updateErr != nil {
			srv.logger.Error("can't update request status", zap.Error(err))
		}

		return
	}
}

// List returns []*request.Resource
func (srv *Service) List(ctx context.Context, params *query.Params) ([]interface{}, error) {
	requests, err := srv.requestRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

// Retrieve function check if user has request and return it
func (srv *Service) Retrieve(ctx context.Context, userID, relationID int64) (jsonapi.Linkable, error) {
	id, err := srv.requestRepo.RelationExists(ctx, userID, relationID)
	if err != nil {
		return nil, err
	}

	if id == 0 {
		return nil, errors.New("request does not exists")
	}

	return srv.requestRepo.Retrieve(ctx, id)
}

func (srv *Service) rabbitPublish(res *request.Resource) error {
	req := compress.NewRequest(res)
	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	return srv.publisher.Publish(body)
}
