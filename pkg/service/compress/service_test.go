package compress

import (
	"context"
	"fmt"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

func TestService_AddConvertedVideo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		video        *video.Resource
		mock         func()
		expectedID   int64
		errorPresent bool
	}{
		{
			name: "Valid db connection",
			video: &video.Resource{
				Name:        "converted_video.mkv",
				UserID:      1,
				Size:        12500,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				ServiceID:   "mock_service",
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs(64000, "converted_video.mkv", 4, 3, 800, 600, "mock_service", 12500, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "converted_video.mkv", 12500, 64000, 800, 600, 4, 3, "mock_service_id"))
			},
			expectedID:   1,
			errorPresent: false,
		},

		{
			name: "Invalid db connection",
			video: &video.Resource{
				Name:        "converted_video.mkv",
				UserID:      1,
				Size:        12500,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				ServiceID:   "mock_service",
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs(64000, "converted_video.mkv", 4, 3, 800, 600, "mock_service", 12500, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectedID:   0,
			errorPresent: true,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			reqRepo := request.NewRepository(db)
			vRepo := video.NewRepository(db)
			logger := zap.NewExample()
			defer logger.Sync()

			srv := NewService(reqRepo, vRepo, logger)

			id, err := srv.AddConvertedVideo(context.Background(), testCase.video)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error")
			}

			if id != testCase.expectedID {
				t.Errorf("Invalid ID, expected: %d, got: %d\n",
					testCase.expectedID, id)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		data         []byte
		mock         func()
		errorPresent bool
	}{
		{
			name:         "With invalid json",
			data:         []byte(`{ data: "check"`),
			mock:         func() {},
			errorPresent: true,
		},
		{
			name: "With error in response and invalid db connection",
			data: []byte(`{"request_id":1,"error":"Invalid ffmpeg path"}`),
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Invalid ffmpeg path", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name: "With error in response and valid db connection",
			data: []byte(`{"request_id":1,"error":"Invalid ffmpeg path"}`),
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Invalid ffmpeg path", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.user_id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.user_id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, 1, "failed", "Invalid ffmpeg path", 64000, 800, 600, 4, 3, "new_video", 1, "my_name.mkv",
						1258000, 0, 0, 0, 0, 0, "mock_service_id", nil, nil, nil, nil,
						nil, nil, nil, nil, nil))
			},
			errorPresent: true,
		},
		{
			name: "Without ConvertedVideo in response and invalid db connection",
			data: []byte(`{"request_id":1}`),
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Converted video does not present", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: false,
		},
		{
			name: "Without ConvertedVideo in response and valid db connection",
			data: []byte(`{"request_id":1}`),
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Converted video does not present", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.user_id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.user_id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, 1, "failed", "Converted video does not present", 64000, 800, 600, 4, 3, "new_video", 1, "my_name.mkv",
						1258000, 0, 0, 0, 0, 0, "mock_service_id", nil, nil, nil, nil,
						nil, nil, nil, nil, nil))
			},
			errorPresent: false,
		},
		{
			name: "With ConvertedVideo in response, invalid db connection for videos and requests",
			data: []byte(`{"request_id":1,"converted_video":{"name":"converted_video.mkv","user_id":1,"size":12500,"bitrate":64000,"resolution_x":800,"resolution_y":600,"ratio_x":4,"ratio_y":3,"service_id":"mock_service"}}`),
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs(64000, "converted_video.mkv", 4, 3, 800, 600, "mock_service", 12500, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Сan't add converted video to db, id: mock_service", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: false,
		},
		{
			name: "With ConvertedVideo in response, valid db connection for videos and invalid for requests",
			data: []byte(`{"request_id":1,"converted_video":{"name":"converted_video.mkv","user_id":1,"size":12500,"bitrate":64000,"resolution_x":800,"resolution_y":600,"ratio_x":4,"ratio_y":3,"service_id":"mock_service"}}`),
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs(64000, "converted_video.mkv", 4, 3, 800, 600, "mock_service", 12500, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "converted_video.mkv", 12500, 64000, 800, 600, 4, 3, "mock_service_id"))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs(1, completedStatus, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: false,
		},
		{
			name: "With ConvertedVideo in response, valid db connection for videos requests",
			data: []byte(`{"request_id":1,"converted_video":{"name":"converted_video.mkv","user_id":1,"size":12500,"bitrate":64000,"resolution_x":800,"resolution_y":600,"ratio_x":4,"ratio_y":3,"service_id":"mock_service"}}`),
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs(64000, "converted_video.mkv", 4, 3, 800, 600, "mock_service", 12500, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(2, "converted_video.mkv", 12500, 64000, 800, 600, 4, 3, "mock_service_id"))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs(2, completedStatus, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.user_id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.user_id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, 1, completedStatus, "", 64000, 800, 600, 4, 3, "new_video", 1, "my_name.mkv",
						1258000, 0, 0, 0, 0, 0, "mock_service_id", 2, "converted_video.mkv", 12500, 64000,
						800, 600, 4, 3, "mock_service_id"))
			},
			errorPresent: false,
		},
		{
			name: "With OriginalVideo in response, invalid db connection",
			data: []byte(`{"request_id":1,"original_video":{"id":1,"bitrate":64000,"resolution_x":800,"resolution_y":600,"ratio_x":4,"ratio_y":3}}`),
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", video.TableName)).
					WithArgs(64000, 4, 3, 800, 600, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name: "With OriginalVideo in response, valid db connection",
			data: []byte(`{"request_id":1,"original_video":{"id":1,"bitrate":64000,"resolution_x":800,"resolution_y":600,"ratio_x":4,"ratio_y":3}}`),
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", video.TableName)).
					WithArgs(64000, 4, 3, 800, 600, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 789569, 64000, 800, 600, 4, 3, "mock_service_id"))
			},
			errorPresent: false,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			reqRepo := request.NewRepository(db)
			vRepo := video.NewRepository(db)

			logger := zap.NewExample()
			defer logger.Sync()

			srv := NewService(reqRepo, vRepo, logger)

			err := srv.UpdateRequest(context.Background(), testCase.data)

			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}

func TestService_UpdateOriginalVideo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		video        *video.Resource
		mock         func()
		errorPresent bool
	}{
		{
			name: "With invalid id",
			video: &video.Resource{
				ID:          0,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
			},
			mock: func() {

			},
			errorPresent: true,
		},
		{
			name: "With invalid db connection",
			video: &video.Resource{
				ID:          1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", video.TableName)).
					WithArgs(64000, 4, 3, 800, 600, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name: "With valid db connection",
			video: &video.Resource{
				ID:          1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", video.TableName)).
					WithArgs(64000, 4, 3, 800, 600, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 789569, 64000, 800, 600, 4, 3, "mock_service_id"))
			},
			errorPresent: false,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			reqRepo := request.NewRepository(db)
			vRepo := video.NewRepository(db)

			logger := zap.NewExample()
			defer logger.Sync()

			srv := NewService(reqRepo, vRepo, logger)

			err := srv.UpdateOriginalVideo(context.Background(), testCase.video)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}

func TestService_UpdateRequestStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		id           int64
		status       string
		details      string
		mock         func()
		errorPresent bool
	}{
		{
			name:         "With invalid id",
			id:           0,
			status:       "failed",
			details:      "bad connection",
			mock:         func() {},
			errorPresent: true,
		},
		{
			name:    "With invalid db connection",
			id:      1,
			status:  "failed",
			details: "Can't add video to database",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Can't add video to database", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name:    "With valid db connection",
			id:      1,
			status:  "failed",
			details: "Can't add video to database",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Can't add video to database", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT requests.id, requests.user_id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.user_id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, 1, "failed", "Can't add video to database", 64000, 800, 600, 4, 3, "new_video", 1, "my_name.mkv",
						1258000, 0, 0, 0, 0, 0, "mock_service_id", nil, nil, nil, nil,
						nil, nil, nil, nil, nil))
			},
			errorPresent: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			reqRepo := request.NewRepository(db)
			vRepo := video.NewRepository(db)

			logger := zap.NewExample()
			defer logger.Sync()

			srv := NewService(reqRepo, vRepo, logger)

			err := srv.UpdateRequestStatus(context.Background(), testCase.id, testCase.status, testCase.details)

			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Shoud be error\n")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}
