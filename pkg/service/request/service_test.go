package request

import (
	"context"
	"errors"
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/google/jsonapi"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"mime/multipart"
	"testing"
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

func TestCreate(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name     string
		resource jsonapi.Linkable
		mock     func()

		expectedRequestId         int64
		expectedRequestStatus     string
		expectedRequestDetails    string
		expectedRequestBitrate    int64
		expectedRequestResolution string
		expectedRequestRatio      string

		expectedOriginalVideoId         int64
		expectedOriginalVideoSize       int64
		expectedOriginalVideoBitrate    int64
		expectedOriginalVideoName       string
		expectedOriginalVideoResolution string
		expectedOriginalVideoRatio      string
		expectedOriginalVideoServiceId  string

		errorPresent bool
	}{
		{
			name:         "Invalid jsonapi.Linkable",
			resource:     new(invalidLinkable),
			mock:         func() {},
			errorPresent: true,
		},
		{
			name: "Invalid db connection to create request",
			resource: &request.Resource{
				UserID:     1,
				Bitrate:    64000,
				Resolution: "800:600",
				Ratio:      "4:3",
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, "800:600", "4:3", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name: "Valid db connection to create request, invalid cloud connection, invalid db connection to update request",
			resource: &request.Resource{
				UserID:     1,
				Bitrate:    64000,
				Resolution: "800:600",
				Ratio:      "4:3",
				VideoRequest: &multipart.FileHeader{
					Filename: "failed",
				},
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, "800:600", "4:3", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "original_in_review", "", 64000, "800:600", "4:3"))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs(`Can't upload video to cloud`, "failed", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name: "Valid db connection to create request, invalid cloud connection, valid db connection to update request",
			resource: &request.Resource{
				UserID:     1,
				Bitrate:    64000,
				Resolution: "800:600",
				Ratio:      "4:3",
				VideoRequest: &multipart.FileHeader{
					Filename: "failed",
				},
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, "800:600", "4:3", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "original_in_review", "", 64000, "800:600", "4:3"))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs(`Can't upload video to cloud`, "failed", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "failed", `Can't upload video to cloud`, 64000, "800:600", "4:3"))
			},
			errorPresent: true,
		},
		{
			name: "Valid db connection to create request, invalid db connection to create video, invalid db connection to update request",
			resource: &request.Resource{
				UserID:     1,
				Bitrate:    64000,
				Resolution: "800:600",
				Ratio:      "4:3",
				VideoRequest: &multipart.FileHeader{
					Filename: "good",
				},
				OriginalVideo: &video.Resource{
					Name:      "my_name.mkv",
					Size:      1258000,
					ServiceId: "mock_service_id",
				},
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, "800:600", "4:3", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "original_in_review", "", 64000, "800:600", "4:3"))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs("my_name.mkv", 1258000, "mock_service_id").
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Can't add video to database", "failed", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name: "Valid db connection to create request, invalid db connection to create video, valid db connection to update request",
			resource: &request.Resource{
				UserID:     1,
				Bitrate:    64000,
				Resolution: "800:600",
				Ratio:      "4:3",
				VideoRequest: &multipart.FileHeader{
					Filename: "good",
				},
				OriginalVideo: &video.Resource{
					Name:      "my_name.mkv",
					Size:      1258000,
					ServiceId: "mock_service_id",
				},
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, "800:600", "4:3", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "original_in_review", "", 64000, "800:600", "4:3"))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs("my_name.mkv", 1258000, "mock_service_id").
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs("Can't add video to database", "failed", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "failed", "Can't add video to database", 64000, "800:600", "4:3"))
			},
			errorPresent: true,
		},
		{
			name: "Should add request, video. Upload video to cloud. Update request",
			resource: &request.Resource{
				UserID:     1,
				Bitrate:    64000,
				Resolution: "800:600",
				Ratio:      "4:3",
				VideoRequest: &multipart.FileHeader{
					Filename: "good",
				},
				OriginalVideo: &video.Resource{
					Name:      "my_name.mkv",
					Size:      1258000,
					ServiceId: "mock_service_id",
				},
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, "800:600", "4:3", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "original_in_review", "", 64000, "800:600", "4:3"))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs("my_name.mkv", 1258000, "mock_service_id").
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution, ratio, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution", "ratio", "service_id"}).
						AddRow(1, "my_name.mkv", 1258000, 0, "", "", "mock_service_id"))
			},

			expectedRequestId:         1,
			expectedRequestStatus:     "original_in_review",
			expectedRequestDetails:    "",
			expectedRequestBitrate:    64000,
			expectedRequestResolution: "800:600",
			expectedRequestRatio:      "4:3",

			expectedOriginalVideoId:         1,
			expectedOriginalVideoSize:       1258000,
			expectedOriginalVideoName:       "my_name.mkv",
			expectedOriginalVideoResolution: "",
			expectedOriginalVideoBitrate:    0,
			expectedOriginalVideoRatio:      "",
			expectedOriginalVideoServiceId:  "mock_service_id",

			errorPresent: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			rRepo := request.NewRepository(db)
			vRepo := video.NewRepository(db)
			cs := new(cloudMock)
			srv := NewService(rRepo, vRepo, cs)

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

				if req.Resolution != testCase.expectedRequestResolution {
					t.Errorf("Invalid request resolution, expected: %s, got: %s\n",
						testCase.expectedRequestResolution, req.Resolution)
				}

				if req.Ratio != testCase.expectedRequestRatio {
					t.Errorf("Invalid request ratio, expected: %s, got: %s\n",
						testCase.expectedRequestRatio, req.Ratio)
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

				if originVideo.Resolution != testCase.expectedOriginalVideoResolution {
					t.Errorf("Invalid original video resolution, expected: %s, got: %s\n",
						testCase.expectedOriginalVideoResolution, originVideo.Resolution)
				}

				if originVideo.Bitrate != testCase.expectedOriginalVideoBitrate {
					t.Errorf("Invalid original video bitrate, expected: %d, got: %d\n",
						testCase.expectedOriginalVideoBitrate, originVideo.Bitrate)
				}

				if originVideo.Ratio != testCase.expectedOriginalVideoRatio {
					t.Errorf("Invalid original video ratio, expected: %s, got: %s\n",
						testCase.expectedOriginalVideoRatio, originVideo.Ratio)
				}

				if originVideo.ServiceId != testCase.expectedOriginalVideoServiceId {
					t.Errorf("Invalid original service id, expected: %s, got: %s\n",
						testCase.expectedOriginalVideoServiceId, originVideo.ServiceId)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}
