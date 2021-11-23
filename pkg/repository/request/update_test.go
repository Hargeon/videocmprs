package request

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

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
