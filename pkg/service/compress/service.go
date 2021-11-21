package compress

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Hargeon/videocmprs/pkg/repository"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
)

const (
	failedStatus    = "failed"
	completedStatus = "success"
)

// Service for updating request and original video in db.
// And adding converted video to db.
type Service struct {
	reqRepo repository.Updater
	vRepo   repository.VideoRepository
}

// NewService initialize Service
func NewService(reqRepo repository.Updater, vRepo repository.VideoRepository) *Service {
	return &Service{
		reqRepo: reqRepo,
		vRepo:   vRepo,
	}
}

// UpdateRequest request status and add converted video to db
func (srv *Service) UpdateRequest(ctx context.Context, data []byte) error {
	res := new(Response)
	if err := json.Unmarshal(data, res); err != nil {
		log.Println("Unmarshal from rabbit", err) // TODO need add logger

		return ErrInvalidResponse
	}

	// check if compress worker got and error
	if res.Error != "" {
		err := srv.UpdateRequestStatus(ctx, res.RequestID, failedStatus, res.Error)

		if err != nil {
			log.Println("update request status", err) // TODO need add logger

			return err
		}

		return ErrCompressWorker
	}

	// add converted video to db
	if res.ConvertedVideo != nil {
		id, err := srv.AddConvertedVideo(ctx, res.ConvertedVideo)

		if err == nil {
			fields := map[string]interface{}{"status": completedStatus, "converted_file_id": id}
			_, reqErr := srv.reqRepo.Update(context.Background(), res.RequestID, fields)
			if reqErr != nil {
				log.Println("updating request status", reqErr) // TODO need add logger
			}
		} else {
			msg := fmt.Sprintf("Сan't add converted video to db, id: %s",
				res.ConvertedVideo.ServiceID)
			err = srv.UpdateRequestStatus(ctx, res.RequestID, failedStatus, msg)
			if err != nil {
				log.Println(err) // TODO need add logger
			}
		}
	} else {
		msg := "Converted video does not present"
		err := srv.UpdateRequestStatus(ctx, res.RequestID, failedStatus, msg)
		if err != nil {
			log.Println(err) // TODO need add logger
		}
	}

	// update original video in db
	if res.OriginalVideo != nil {
		err := srv.UpdateOriginalVideo(ctx, res.OriginalVideo)

		return err
	}

	return nil
}

// UpdateRequestStatus function update status and details for request
func (srv *Service) UpdateRequestStatus(ctx context.Context, id int64, status, details string) error {
	if id <= 0 {
		return ErrInvalidID
	}

	fields := map[string]interface{}{"status": status, "details": details}
	_, err := srv.reqRepo.Update(ctx, id, fields)

	return err
}

// AddConvertedVideo function add converted video to db
func (srv *Service) AddConvertedVideo(ctx context.Context, v *video.Resource) (int64, error) {
	convertedVideo, err := srv.vRepo.Create(ctx, v.BuildFields())
	if err != nil {
		return 0, err
	}

	res, ok := convertedVideo.(*video.Resource)
	if !ok {
		return 0, ErrInvalidTypeAssertion
	}

	return res.ID, nil
}

// UpdateOriginalVideo function update bitrate, resolution and ratio for original video
func (srv *Service) UpdateOriginalVideo(ctx context.Context, v *video.Resource) error {
	id := v.ID
	if id <= 0 {
		return ErrInvalidID
	}

	_, err := srv.vRepo.Update(ctx, id, v.BuildFields())

	return err
}
