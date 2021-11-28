package video

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
		name   string
		id     int64
		fields map[string]interface{}
		mock   func()

		expectedID          int64
		expectedName        string
		expectedSize        int64
		expectedBitrate     int64
		expectedResolutionX int
		expectedResolutionY int
		expectedRatioX      int
		expectedRatioY      int
		expectedServiceID   string
		errorPresent        bool
	}{
		{
			name: "Invalid db connection",
			id:   1,
			fields: map[string]interface{}{
				"bitrate":      64000,
				"resolution_x": 800,
				"resolution_y": 600,
				"ratio_x":      4,
				"ratio_y":      3,
			},

			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", TableName)).
					WithArgs(64000, 4, 3, 800, 600, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},

			errorPresent: true,
		},
		{
			name: "Valid db connection",
			id:   1,
			fields: map[string]interface{}{
				"bitrate":      64000,
				"resolution_x": 800,
				"resolution_y": 600,
				"ratio_x":      4,
				"ratio_y":      3,
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", TableName)).
					WithArgs(64000, 4, 3, 800, 600, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 789569, 64000, 800, 600, 4, 3, "mock_service_id"))
			},
			expectedID:          1,
			expectedName:        "my_name.mkv",
			expectedBitrate:     64000,
			expectedSize:        789569,
			expectedResolutionX: 800,
			expectedResolutionY: 600,
			expectedRatioX:      4,
			expectedRatioY:      3,
			expectedServiceID:   "mock_service_id",
			errorPresent:        false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewRepository(db)
			linkable, err := repo.Update(context.Background(), testCase.id, testCase.fields)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error when updating a video, error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				res, ok := linkable.(*Resource)
				if !ok {
					t.Fatalf("Invalid type assertion\n")
				}

				if res.ID != testCase.expectedID {
					t.Errorf("Invalid id, expected: %d, got: %d\n",
						testCase.expectedID, res.ID)
				}

				if res.Name != testCase.expectedName {
					t.Errorf("Invalid name, expected: %s, got: %s\n",
						testCase.expectedName, res.Name)
				}

				if res.Size != testCase.expectedSize {
					t.Errorf("Invalid size, expected: %d, got: %d\n",
						testCase.expectedSize, res.Size)
				}

				if res.Bitrate != testCase.expectedBitrate {
					t.Errorf("Invalid bitrate, expected %d, got: %d\n",
						testCase.expectedBitrate, res.Bitrate)
				}

				if res.ResolutionX != testCase.expectedResolutionX {
					t.Errorf("Invalid resolution x, expected: %d, got: %d\n",
						testCase.expectedResolutionX, res.ResolutionX)
				}

				if res.ResolutionY != testCase.expectedResolutionY {
					t.Errorf("Invlaid resolution y, expected: %d, got: %d\n",
						testCase.expectedResolutionY, res.ResolutionY)
				}

				if res.RatioX != testCase.expectedRatioX {
					t.Errorf("Invalid ratio x, expected: %d, got: %d\n",
						testCase.expectedRatioX, res.RatioX)
				}

				if res.RatioY != testCase.expectedRatioY {
					t.Errorf("Invalid ratio y, expected: %d, got: %d\n",
						testCase.expectedRatioY, res.RatioY)
				}

				if res.ServiceID != testCase.expectedServiceID {
					t.Errorf("Invalid service id, expected %s, got: %s\n",
						testCase.expectedServiceID, res.ServiceID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}
