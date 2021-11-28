package request

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"testing"

	"github.com/Hargeon/videocmprs/api/query"
	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/jsonapi"
)

type invalidLinkable struct{}

func (r *invalidLinkable) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "",
	}
}

type cloudMock struct{}

func (c *cloudMock) Upload(ctx context.Context, header *multipart.FileHeader) (string, error) {
	if header.Filename == "failed" {
		return "", errors.New("failed connection")
	}

	return "mock_service_id", nil
}

type rabbitSuccess struct{}

type rabbitError struct{}

func (r *rabbitSuccess) Publish(body []byte) error {
	return nil
}

func (r *rabbitError) Publish(body []byte) error {
	return errors.New("mock error")
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name      string
		resource  jsonapi.Linkable
		publisher service.Publisher
		mock      func()

		expectedRequestId          int64
		expectedRequestStatus      string
		expectedRequestDetails     string
		expectedRequestBitrate     int64
		expectedRequestResolutionX int
		expectedRequestResolutionY int
		expectedRequestRatioX      int
		expectedRequestRatioY      int
		expectedRequestVideoName   string

		expectedOriginalVideoId          int64
		expectedOriginalVideoSize        int64
		expectedOriginalVideoBitrate     int64
		expectedOriginalVideoName        string
		expectedOriginalVideoResolutionX int
		expectedOriginalVideoResolutionY int
		expectedOriginalVideoRatioX      int
		expectedOriginalVideoRatioY      int
		expectedOriginalVideoServiceId   string

		errorPresent bool
	}{
		{
			name:         "Invalid jsonapi.Linkable",
			resource:     new(invalidLinkable),
			publisher:    &rabbitSuccess{},
			mock:         func() {},
			errorPresent: true,
		},
		{
			name: "Invalid db connection to create request",
			resource: &request.Resource{
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				VideoName:   "new_video",
			},
			publisher: &rabbitSuccess{},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, 800, 600, 4, 3, 1, "new_video").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name: "Valid db connection to create request, invalid cloud connection, invalid db connection to update request",
			resource: &request.Resource{
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				VideoName:   "new_video",
				VideoRequest: &multipart.FileHeader{
					Filename: "failed",
				},
			},
			publisher: &rabbitSuccess{},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, 800, 600, 4, 3, 1, "new_video").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "original_in_review", "", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs(`Can't upload video to cloud`, "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name: "Valid db connection to create request, invalid cloud connection, valid db connection to update request",
			resource: &request.Resource{
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				VideoName:   "new_video",
				VideoRequest: &multipart.FileHeader{
					Filename: "failed",
				},
			},
			publisher: &rabbitSuccess{},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, 800, 600, 4, 3, 1, "new_video").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "original_in_review", "", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs(`Can't upload video to cloud`, "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "failed", "Can't upload video to cloud", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil))
			},
			errorPresent: true,
		},
		{
			name: "Valid db connection to create request, invalid db connection to create video, invalid db connection to update request",
			resource: &request.Resource{
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				VideoName:   "new_video",
				VideoRequest: &multipart.FileHeader{
					Filename: "good",
				},
				OriginalVideo: &video.Resource{
					Name:      "my_name.mkv",
					Size:      1258000,
					ServiceID: "mock_service_id",
				},
			},
			publisher: &rabbitSuccess{},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, 800, 600, 4, 3, 1, "new_video").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "original_in_review", "", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs("my_name.mkv", "mock_service_id", 1258000).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Can't add video to database", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name: "Valid db connection to create request, invalid db connection to create video, valid db connection to update request",
			resource: &request.Resource{
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				VideoName:   "new_video",
				VideoRequest: &multipart.FileHeader{
					Filename: "good",
				},
				OriginalVideo: &video.Resource{
					Name:      "my_name.mkv",
					Size:      1258000,
					ServiceID: "mock_service_id",
				},
			},
			publisher: &rabbitSuccess{},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, 800, 600, 4, 3, 1, "new_video").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "original_in_review", "", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs("my_name.mkv", "mock_service_id", 1258000).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Can't add video to database", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "failed", "Can't add video to database", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil))
			},
			errorPresent: true,
		},
		{
			name: "Should add request, video. Upload video to cloud. Update request",
			resource: &request.Resource{
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				VideoName:   "new_video",
				VideoRequest: &multipart.FileHeader{
					Filename: "good",
				},
				OriginalVideo: &video.Resource{
					Name:      "my_name.mkv",
					Size:      1258000,
					ServiceID: "mock_service_id",
				},
			},
			publisher: &rabbitSuccess{},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, 800, 600, 4, 3, 1, "new_video").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "original_in_review", "", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs("my_name.mkv", "mock_service_id", 1258000).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 1258000, 0, 0, 0, 0, 0, "mock_service_id"))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "original_in_review", "", 64000, 800, 600, 4, 3, "new_video", 1, "my_name.mkv",
						1258000, 0, 0, 0, 0, 0, "mock_service_id", nil, nil, nil, nil,
						nil, nil, nil, nil, nil))
			},

			expectedRequestId:          1,
			expectedRequestStatus:      "original_in_review",
			expectedRequestDetails:     "",
			expectedRequestBitrate:     64000,
			expectedRequestResolutionX: 800,
			expectedRequestResolutionY: 600,
			expectedRequestRatioX:      4,
			expectedRequestRatioY:      3,
			expectedRequestVideoName:   "new_video",

			expectedOriginalVideoId:          1,
			expectedOriginalVideoSize:        1258000,
			expectedOriginalVideoName:        "my_name.mkv",
			expectedOriginalVideoResolutionX: 0,
			expectedOriginalVideoResolutionY: 0,
			expectedOriginalVideoBitrate:     0,
			expectedOriginalVideoRatioX:      0,
			expectedOriginalVideoRatioY:      0,
			expectedOriginalVideoServiceId:   "mock_service_id",

			errorPresent: false,
		},
		{
			name: "With invalid rabbit connection, should update request status",
			resource: &request.Resource{
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				VideoName:   "new_video",
				VideoRequest: &multipart.FileHeader{
					Filename: "good",
				},
				OriginalVideo: &video.Resource{
					Name:      "my_name.mkv",
					Size:      1258000,
					ServiceID: "mock_service_id",
				},
			},
			publisher: &rabbitError{},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, 800, 600, 4, 3, 1, "new_video").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "original_in_review", "", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs("my_name.mkv", "mock_service_id", 1258000).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 1258000, 0, 0, 0, 0, 0, "mock_service_id"))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "original_in_review", "", 64000, 800, 600, 4, 3, "new_video", 1, "my_name.mkv",
						1258000, 0, 0, 0, 0, 0, "mock_service_id", nil, nil, nil, nil,
						nil, nil, nil, nil, nil))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Failed connection to worker", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "failed", "Failed connection to worker", 64000, 800, 600, 4, 3, "new_video", 1, "my_name.mkv",
						1258000, 0, 0, 0, 0, 0, "mock_service_id", nil, nil, nil, nil,
						nil, nil, nil, nil, nil))
			},

			expectedRequestId:          1,
			expectedRequestStatus:      "failed",
			expectedRequestDetails:     "Failed connection to worker",
			expectedRequestBitrate:     64000,
			expectedRequestResolutionX: 800,
			expectedRequestResolutionY: 600,
			expectedRequestRatioX:      4,
			expectedRequestRatioY:      3,
			expectedRequestVideoName:   "new_video",

			expectedOriginalVideoId:          1,
			expectedOriginalVideoSize:        1258000,
			expectedOriginalVideoName:        "my_name.mkv",
			expectedOriginalVideoResolutionX: 0,
			expectedOriginalVideoResolutionY: 0,
			expectedOriginalVideoBitrate:     0,
			expectedOriginalVideoRatioX:      0,
			expectedOriginalVideoRatioY:      0,
			expectedOriginalVideoServiceId:   "mock_service_id",

			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			rRepo := request.NewRepository(db)
			vRepo := video.NewRepository(db)
			cs := new(cloudMock)
			srv := NewService(rRepo, vRepo, cs, testCase.publisher)

			linkable, err := srv.Create(context.Background(), testCase.resource)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				req, ok := linkable.(*request.Resource)
				if !ok {
					t.Fatalf("Invalid type assertion *request.Resource\n")
				}

				if req.ID != testCase.expectedRequestId {
					t.Errorf("Invalid request id, expected: %d, got: %d\n",
						testCase.expectedRequestId, req.ID)
				}

				if req.Status != testCase.expectedRequestStatus {
					t.Errorf("Invalid request status, expected: %s, got: %s\n",
						testCase.expectedRequestStatus, req.Status)
				}

				if req.Details != testCase.expectedRequestDetails {
					t.Errorf("Invalid request details, expected: %s, got: %s\n",
						testCase.expectedRequestDetails, req.Details)
				}

				if req.Bitrate != testCase.expectedRequestBitrate {
					t.Errorf("Invalid request bitrate, expected: %d, got: %d\n",
						testCase.expectedRequestBitrate, req.Bitrate)
				}

				if req.ResolutionX != testCase.expectedRequestResolutionX {
					t.Errorf("Invalid request resolution, expected: %d, got: %d\n",
						testCase.expectedRequestResolutionX, req.ResolutionX)
				}

				if req.ResolutionY != testCase.expectedRequestResolutionY {
					t.Errorf("Invalid request resolution, expected: %d, got: %d\n",
						testCase.expectedRequestResolutionY, req.ResolutionY)
				}

				if req.RatioX != testCase.expectedRequestRatioX {
					t.Errorf("Invalid request ratio, expected: %d, got: %d\n",
						testCase.expectedRequestRatioX, req.RatioX)
				}

				if req.RatioY != testCase.expectedRequestRatioY {
					t.Errorf("Invalid request ratio, expected: %d, got: %d\n",
						testCase.expectedRequestRatioY, req.RatioY)
				}

				if req.VideoName != testCase.expectedRequestVideoName {
					t.Errorf("Invalid reqest name, expected: %s, got: %s\n",
						testCase.expectedRequestVideoName, req.VideoName)
				}

				originVideo := req.OriginalVideo
				if originVideo.ID != testCase.expectedOriginalVideoId {
					t.Errorf("Invalid original video id, expected: %d, got: %d\n",
						testCase.expectedOriginalVideoId, originVideo.ID)
				}

				if originVideo.Size != testCase.expectedOriginalVideoSize {
					t.Errorf("Invalid origin video size, expected: %d, got: %d\n",
						testCase.expectedOriginalVideoSize, originVideo.Size)
				}

				if originVideo.Name != testCase.expectedOriginalVideoName {
					t.Errorf("Invalid original video name, expected: %s, got: %s\n",
						testCase.expectedOriginalVideoName, originVideo.Name)
				}

				if originVideo.ResolutionX != testCase.expectedOriginalVideoResolutionX {
					t.Errorf("Invalid original video resolution, expected: %d, got: %d\n",
						testCase.expectedOriginalVideoResolutionX, originVideo.ResolutionX)
				}

				if originVideo.ResolutionY != testCase.expectedOriginalVideoResolutionY {
					t.Errorf("Invalid original video resolution, expected: %d, got: %d\n",
						testCase.expectedOriginalVideoResolutionY, originVideo.ResolutionY)
				}

				if originVideo.Bitrate != testCase.expectedOriginalVideoBitrate {
					t.Errorf("Invalid original video bitrate, expected: %d, got: %d\n",
						testCase.expectedOriginalVideoBitrate, originVideo.Bitrate)
				}

				if originVideo.RatioX != testCase.expectedOriginalVideoRatioX {
					t.Errorf("Invalid original video ratio, expected: %d, got: %d\n",
						testCase.expectedOriginalVideoRatioX, originVideo.RatioX)
				}

				if originVideo.RatioY != testCase.expectedOriginalVideoRatioY {
					t.Errorf("Invalid original video ratio, expected: %d, got: %d\n",
						testCase.expectedOriginalVideoRatioY, originVideo.RatioY)
				}

				if originVideo.ServiceID != testCase.expectedOriginalVideoServiceId {
					t.Errorf("Invalid original service id, expected: %s, got: %s\n",
						testCase.expectedOriginalVideoServiceId, originVideo.ServiceID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		params       *query.Params
		mock         func()
		expectedLen  int
		errorPresent bool
	}{
		{
			name: "Zero requests",
			params: &query.Params{
				RelationID: 1,
				PageNumber: 0,
				PageSize:   10,
			},
			mock: func() {
				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}))
			},
			expectedLen:  0,
			errorPresent: false,
		},

		{
			name: "One request",
			params: &query.Params{
				RelationID: 1,
				PageNumber: 0,
				PageSize:   10,
			},
			mock: func() {
				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "", "", 64000, 800, 600, 4, 3, "new_video", 1, "new_video", 15000,
						78000, 1200, 800, 6, 5, "new_service_id", 2, "converted_video", 12000, 64000,
						800, 600, 4, 3, "converted_service_id"))
			},
			expectedLen:  1,
			errorPresent: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			rRepo := request.NewRepository(db)
			vRepo := video.NewRepository(db)
			cs := new(cloudMock)
			srv := NewService(rRepo, vRepo, cs, &rabbitSuccess{})
			res, err := srv.List(context.Background(), testCase.params)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				if len(res) != testCase.expectedLen {
					t.Fatalf("Invalid number of requests, expected: %d, got: %d\n",
						testCase.expectedLen, len(res))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}

func TestRetrieve(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name                string
		id                  int64
		mock                func()
		expectedID          int64
		expectedStatus      string
		expectedDetails     string
		expectedBitrate     int64
		expectedResolutionX int
		expectedResolutionY int
		expectedRatioX      int
		expectedRatioY      int
		expectedName        string
		errorPresent        bool
	}{
		{
			name: "Should return request",
			id:   1,
			mock: func() {
				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "original_in_review", "", 1589875, 800, 600, 4, 3, "new_video", 1, "new_video", 15000,
						78000, 1200, 800, 6, 5, "new_service_id", 2, "converted_video", 12000, 64000,
						800, 600, 4, 3, "converted_service_id"))
			},
			expectedID:          1,
			expectedStatus:      "original_in_review",
			expectedDetails:     "",
			expectedBitrate:     1589875,
			expectedResolutionX: 800,
			expectedResolutionY: 600,
			expectedRatioX:      4,
			expectedRatioY:      3,
			expectedName:        "new_video",
			errorPresent:        false,
		},
		{
			name: "Should not return request",
			id:   1,
			mock: func() {
				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}))
			},
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			testCase.mock()

			rRepo := request.NewRepository(db)
			vRepo := video.NewRepository(db)
			cs := new(cloudMock)

			srv := NewService(rRepo, vRepo, cs, &rabbitSuccess{})
			linkable, err := srv.Retrieve(context.Background(), testCase.id)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				request, ok := linkable.(*request.Resource)
				if !ok {
					t.Fatalf("Invalid type assertion for *request.Resource")
				}

				if request.ID != testCase.expectedID {
					t.Errorf("Invalid id, expected: %d, got: %d\n",
						testCase.expectedID, request.ID)
				}

				if request.Status != testCase.expectedStatus {
					t.Errorf("Invalid status, expected: %s, got: %s\n",
						testCase.expectedStatus, request.Status)
				}

				if request.Details != testCase.expectedDetails {
					t.Errorf("Invalid details, expected: %s, got: %s\n",
						testCase.expectedDetails, request.Details)
				}

				if request.Bitrate != testCase.expectedBitrate {
					t.Errorf("Invalid bitrate, expected: %d, got: %d\n",
						testCase.expectedBitrate, request.Bitrate)
				}

				if request.ResolutionX != testCase.expectedResolutionX {
					t.Errorf("Invalid resolution, expected: %d, got: %d\n",
						testCase.expectedResolutionX, request.ResolutionX)
				}

				if request.ResolutionY != testCase.expectedResolutionY {
					t.Errorf("Invalid resolution, expected: %d, got: %d\n",
						testCase.expectedResolutionY, request.ResolutionY)
				}

				if request.RatioX != testCase.expectedRatioX {
					t.Errorf("Invalid ratio, expected: %d, got: %d\n",
						testCase.expectedRatioX, request.RatioX)
				}

				if request.RatioY != testCase.expectedRatioY {
					t.Errorf("Invalid ratio, expected: %d, got: %d\n",
						testCase.expectedRatioY, request.RatioY)
				}

				if request.VideoName != testCase.expectedName {
					t.Errorf("Invalid name, expected: %s, got: %s\n",
						testCase.expectedName, request.VideoName)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
