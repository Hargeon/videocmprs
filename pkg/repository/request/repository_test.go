package request

import (
	"context"
	"fmt"
	"testing"

	"github.com/Hargeon/videocmprs/api/query"
	"github.com/Hargeon/videocmprs/pkg/repository/video"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/jsonapi"
)

type invalidResource struct{}

func (r *invalidResource) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "",
	}
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name                string
		mock                func()
		req                 jsonapi.Linkable
		expectedID          int64
		expectedStatus      string
		expectedDetails     string
		expectedBitrate     int64
		expectedResolutionX int
		expectedResolutionY int
		expectedRatioX      int
		expectedRatioY      int
		expectedVideoName   string
		errorPresent        bool
	}{
		{
			name: "Should add request to db",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs(64000, 800, 600, 4, 3, 1, "new_video").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, video_name FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "video_name"}).
						AddRow(1, "original_in_review", "", 64000, 800, 600, 4, 3, "new_video"))
			},
			req: &Resource{
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				VideoName:   "new_video",
			},
			expectedID:          1,
			expectedStatus:      "original_in_review",
			expectedDetails:     "",
			expectedBitrate:     64000,
			expectedResolutionX: 800,
			expectedResolutionY: 600,
			expectedRatioX:      4,
			expectedRatioY:      3,
			expectedVideoName:   "new_video",
			errorPresent:        false,
		},
		{
			name: "Should not add request to db",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs(64000, 800, 600, 4, 3, 1, "new_video").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			req: &Resource{
				UserID:      1,
				Bitrate:     64000,
				ResolutionX: 800,
				ResolutionY: 600,
				RatioX:      4,
				RatioY:      3,
				VideoName:   "new_video",
			},
			errorPresent: true,
		},
		{
			name: "With invalid json.Linkable",
			mock: func() {
			},
			req:          &invalidResource{},
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			testCase.mock()
			repo := NewRepository(db)
			linkable, err := repo.Create(context.Background(), testCase.req)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				request, ok := linkable.(*Resource)
				if !ok {
					t.Fatalf("Invalid type assetrion *reqest.Resource\n")
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

				if request.VideoName != testCase.expectedVideoName {
					t.Errorf("Invalid name, expected: %s, got: %s\n",
						testCase.expectedVideoName, request.VideoName)
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
				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, video_name FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "video_name"}).
						AddRow(1, "original_in_review", "", 1589875, 800, 600, 4, 3, "new_video"))
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
				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, video_name FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "video_name"}))
			},
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			testCase.mock()
			repo := NewRepository(db)
			linkable, err := repo.Retrieve(context.Background(), testCase.id)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				request, ok := linkable.(*Resource)
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

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name                string
		id                  int64
		fields              map[string]interface{}
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
			name: "Should update request",
			id:   1,
			fields: map[string]interface{}{
				"details": "Can't add video to database",
				"status":  "failed",
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", TableName)).
					WithArgs("Can't add video to database", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, video_name FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "video_name"}).
						AddRow(1, "failed", "Can't add video to database", 64000, 800, 600, 4, 3, "new_video"))
			},
			expectedID:          1,
			expectedStatus:      "failed",
			expectedDetails:     "Can't add video to database",
			expectedBitrate:     64000,
			expectedResolutionX: 800,
			expectedResolutionY: 600,
			expectedRatioX:      4,
			expectedRatioY:      3,
			expectedName:        "new_video",
			errorPresent:        false,
		},

		{
			name: "With bad db connection",
			id:   1,
			fields: map[string]interface{}{
				"details": "Can't add video to database",
				"status":  "failed",
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", TableName)).
					WithArgs("Can't add video to database", "failed", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			testCase.mock()
			repo := NewRepository(db)
			linkable, err := repo.Update(context.Background(), testCase.id, testCase.fields)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				request, ok := linkable.(*Resource)
				if !ok {
					t.Fatalf("Invalid type assertion for *request.Resource\n")
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
		name             string
		params           *query.Params
		expectedRequests []*Resource
		mock             func()
		errorPresent     bool
	}{
		{
			name: "Zero requests",
			params: &query.Params{
				RelationID: 1,
				PageNumber: 0,
				PageSize:   10,
			},
			expectedRequests: []*Resource{},
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
			errorPresent: false,
		},

		{
			name: "Without origin and converted video",
			params: &query.Params{
				RelationID: 1,
				PageNumber: 0,
				PageSize:   10,
			},
			expectedRequests: []*Resource{
				{
					ID:          1,
					Status:      "",
					Details:     "",
					Bitrate:     64000,
					ResolutionX: 800,
					ResolutionY: 600,
					RatioX:      4,
					RatioY:      3,
					VideoName:   "new_video",
				},
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
						1, "", "", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil))
			},
			errorPresent: false,
		},

		{
			name: "With origin video",
			params: &query.Params{
				RelationID: 1,
				PageNumber: 0,
				PageSize:   10,
			},
			expectedRequests: []*Resource{
				{
					ID:          1,
					Status:      "",
					Details:     "",
					Bitrate:     64000,
					ResolutionX: 800,
					ResolutionY: 600,
					RatioX:      4,
					RatioY:      3,
					VideoName:   "new_video",
					OriginalVideo: &video.Resource{
						ID:          1,
						Name:        "new_video",
						Size:        15000,
						Bitrate:     78000,
						ResolutionX: 1200,
						ResolutionY: 800,
						RatioX:      6,
						RatioY:      5,
						ServiceID:   "new_service_id",
					},
				},
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
						78000, 1200, 800, 6, 5, "new_service_id", nil, nil, nil, nil, nil, nil, nil, nil, nil))
			},
			errorPresent: false,
		},

		{
			name: "With origin and converted video",
			params: &query.Params{
				RelationID: 1,
				PageNumber: 0,
				PageSize:   10,
			},
			expectedRequests: []*Resource{
				{
					ID:          1,
					Status:      "",
					Details:     "",
					Bitrate:     64000,
					ResolutionX: 800,
					ResolutionY: 600,
					RatioX:      4,
					RatioY:      3,
					VideoName:   "new_video",
					OriginalVideo: &video.Resource{
						ID:          1,
						Name:        "new_video",
						Size:        15000,
						Bitrate:     78000,
						ResolutionX: 1200,
						ResolutionY: 800,
						RatioX:      6,
						RatioY:      5,
						ServiceID:   "new_service_id",
					},
					ConvertedVideo: &video.Resource{
						ID:          2,
						Name:        "converted_video",
						Size:        12000,
						Bitrate:     64000,
						ResolutionX: 800,
						ResolutionY: 600,
						RatioX:      4,
						RatioY:      3,
						ServiceID:   "converted_service_id",
					},
				},
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
			errorPresent: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			repo := NewRepository(db)
			requests, err := repo.List(context.Background(), testCase.params)

			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				if len(testCase.expectedRequests) != len(requests) {
					t.Fatalf("Invalid number of requests, expected: %d, got: %d\n",
						len(testCase.expectedRequests), len(requests))
				}

				for i := 0; i < len(requests); i++ {
					request, ok := requests[i].(*Resource)
					if !ok {
						t.Fatalf("Invalid type assertion")
					}
					expRequest := testCase.expectedRequests[i]

					// check requests
					if request.ID != expRequest.ID {
						t.Errorf("Invalid id, expected: %d, got: %d\n",
							expRequest.ID, request.ID)
					}

					if request.Status != expRequest.Status {
						t.Errorf("Invalid status, expected: %s, got: %s\n",
							expRequest.Status, request.Status)
					}

					if request.Details != expRequest.Details {
						t.Errorf("Invalid details, expected: %s, got: %s\n",
							expRequest.Details, request.Details)
					}

					if request.Bitrate != expRequest.Bitrate {
						t.Errorf("Invalid bitrate, expected: %d, got: %d\n",
							expRequest.Bitrate, request.Bitrate)
					}

					if request.ResolutionX != expRequest.ResolutionX {
						t.Errorf("Invalid resolution, expected: %d, got: %d\n",
							expRequest.ResolutionX, request.ResolutionX)
					}

					if request.ResolutionY != expRequest.ResolutionY {
						t.Errorf("Invalid resolution, expected: %d, got: %d\n",
							expRequest.ResolutionY, request.ResolutionY)
					}

					if request.RatioX != expRequest.RatioX {
						t.Errorf("Invalid ratio, expected: %d, got: %d\n",
							expRequest.RatioY, request.RatioX)
					}

					if request.RatioY != expRequest.RatioY {
						t.Errorf("Invalid ratio, expected: %d, got: %d\n",
							expRequest.RatioY, request.RatioY)
					}

					if request.VideoName != expRequest.VideoName {
						t.Errorf("Invalid name, expected: %s, got: %s\n",
							expRequest.VideoName, request.VideoName)
					}

					// check original video
					if request.OriginalVideo != nil && expRequest.OriginalVideo == nil {
						t.Errorf("Invalid Original video, expected: nil, got: %#v\n",
							request.OriginalVideo)
					}

					if request.OriginalVideo == nil && expRequest.OriginalVideo != nil {
						t.Errorf("Invalid Original video, expected: %#v, got: nil\n",
							expRequest.OriginalVideo)
					}

					if request.OriginalVideo != nil && expRequest.OriginalVideo != nil {
						gotVideo := request.OriginalVideo
						expVideo := expRequest.OriginalVideo
						testVideo(t, gotVideo, expVideo)
					}

					// check converted video
					if request.ConvertedVideo != nil && expRequest.ConvertedVideo == nil {
						t.Errorf("Invalid ConvertedVideo video, expected: nil, got: %#v\n",
							request.ConvertedVideo)
					}

					if request.ConvertedVideo == nil && expRequest.ConvertedVideo != nil {
						t.Errorf("Invalid ConvertedVideo video, expected: %#v, got: nil\n",
							expRequest.ConvertedVideo)
					}

					if request.ConvertedVideo != nil && expRequest.ConvertedVideo != nil {
						gotVideo := request.ConvertedVideo
						expVideo := expRequest.ConvertedVideo
						testVideo(t, gotVideo, expVideo)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}

func testVideo(t *testing.T, gotVideo, expVideo *video.Resource) {
	t.Helper()

	if gotVideo.ID != expVideo.ID {
		t.Errorf("Invalid id, expected: %d, got: %d\n",
			expVideo.ID, gotVideo.ID)
	}

	if gotVideo.Name != expVideo.Name {
		t.Errorf("Invalid name, expected: %s, got: %s\n",
			expVideo.Name, gotVideo.Name)
	}

	if gotVideo.Size != expVideo.Size {
		t.Errorf("Invalid size, expected: %d, got: %d\n",
			expVideo.Size, gotVideo.Size)
	}

	if gotVideo.ResolutionX != expVideo.ResolutionX {
		t.Errorf("Invalid resolution, expected: %d, got: %d\n",
			expVideo.ResolutionX, gotVideo.ResolutionX)
	}

	if gotVideo.ResolutionY != expVideo.ResolutionY {
		t.Errorf("Invalid resolution, expected: %d, got: %d\n",
			expVideo.ResolutionY, gotVideo.ResolutionY)
	}

	if gotVideo.RatioX != expVideo.RatioX {
		t.Errorf("Invalid ratio, expected: %d, got: %d\n",
			expVideo.RatioX, gotVideo.RatioX)
	}

	if gotVideo.RatioY != expVideo.RatioY {
		t.Errorf("Invalid ratio, expected: %d, got: %d\n",
			expVideo.RatioY, gotVideo.RatioY)
	}

	if gotVideo.Bitrate != expVideo.Bitrate {
		t.Errorf("Invalid bitrate, expected: %d, got: %d\n",
			expVideo.Bitrate, gotVideo.Bitrate)
	}

	if gotVideo.ServiceID != expVideo.ServiceID {
		t.Errorf("Invalid ServiceID, expected: %s, got: %s\n",
			expVideo.ServiceID, gotVideo.ServiceID)
	}
}
