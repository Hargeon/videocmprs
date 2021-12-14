package video

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name                string
		video               *Resource
		mock                func()
		expectedID          int64
		expectedSize        int64
		expectedBitrate     int64
		expectedName        string
		expectedResolutionX int
		expectedResolutionY int
		expectedRatioX      int
		expectedRatioY      int
		expectedServiceID   string
		errorPresent        bool
	}{
		{
			name: "Should add video",
			video: &Resource{
				Name:      "my_name.mkv",
				Size:      1258000,
				ServiceID: "mock_service_id",
				UserID:    1,
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs("my_name.mkv", "mock_service_id", 1258000, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 1258000, 0, 0, 0, 0, 0, "mock_service_id"))
			},
			expectedID:          1,
			expectedSize:        1258000,
			expectedBitrate:     0,
			expectedName:        "my_name.mkv",
			expectedResolutionX: 0,
			expectedResolutionY: 0,

			expectedRatioX:    0,
			expectedRatioY:    0,
			expectedServiceID: "mock_service_id",
			errorPresent:      false,
		},
		{
			name: "Should not add video",
			video: &Resource{
				Name:      "qwe",
				Size:      125,
				ServiceID: "mock_service_id",
				UserID:    1,
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs("qwe", "mock_service_id", 125, 1).
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
			fields := testCase.video.BuildFields()
			linkable, err := repo.Create(context.Background(), fields)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err)
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				video, ok := linkable.(*Resource)
				if !ok {
					t.Fatalf("Invalid type assertion *video.Resource\n")
				}

				if video.ID != testCase.expectedID {
					t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedID, video.ID)
				}

				if video.Name != testCase.expectedName {
					t.Errorf("Invalid name, expected: %s, got: %s\n", testCase.expectedName, video.Name)
				}

				if video.Size != testCase.expectedSize {
					t.Errorf("Invalid size, expected: %d, got: %d\n", testCase.expectedSize, video.Size)
				}

				if video.ResolutionX != testCase.expectedResolutionX {
					t.Errorf("Invalid resolution, expected: %d, got: %d\n",
						testCase.expectedResolutionX, video.ResolutionX)
				}

				if video.ResolutionY != testCase.expectedResolutionY {
					t.Errorf("Invalid resolution, expected: %d, got: %d\n",
						testCase.expectedResolutionY, video.ResolutionY)
				}

				if video.RatioX != testCase.expectedRatioX {
					t.Errorf("Invalid ratio, expected: %d, got: %d\n",
						testCase.expectedRatioX, video.RatioX)
				}

				if video.RatioY != testCase.expectedRatioY {
					t.Errorf("Invalid ratio, expected: %d, got: %d\n",
						testCase.expectedRatioY, video.RatioY)
				}

				if video.Bitrate != testCase.expectedBitrate {
					t.Errorf("Invalid bitrate, expected: %d, got: %d\n",
						testCase.expectedBitrate, video.Bitrate)
				}

				if video.ServiceID != testCase.expectedServiceID {
					t.Errorf("Invalid ServiceID, expected: %s, got: %s\n",
						testCase.expectedServiceID, video.ServiceID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}
