package video

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/repository/video"

	"github.com/DATA-DOG/go-sqlmock"
)

type cloudMock struct{}

func (c *cloudMock) Upload(ctx context.Context, header *multipart.FileHeader) (string, error) {
	if header.Filename == "failed" {
		return "", errors.New("failed connection")
	}

	return "mock_service_id", nil
}

func (c *cloudMock) URL(filename string) (string, error) {
	if filename == "error" {
		return "", errors.New("mock error")
	}

	return filename, nil
}

func TestRetrieve(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name                string
		id                  int64
		userID              int64
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
			name:   "Should find video",
			id:     1,
			userID: 1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", video.TableName)).
					WithArgs(1, 1).
					WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 1258000, 789569, 700, 600, 4, 3, "mock_service_id"))
			},
			expectedID:          1,
			expectedName:        "my_name.mkv",
			expectedSize:        1258000,
			expectedBitrate:     789569,
			expectedResolutionX: 700,
			expectedResolutionY: 600,
			expectedRatioX:      4,
			expectedRatioY:      3,
			expectedServiceID:   "mock_service_id",
			errorPresent:        false,
		},
		{
			name:   "Should not find video",
			id:     1,
			userID: 1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", video.TableName)).
					WithArgs(1, 1).
					WillReturnRows(mock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			testCase.mock()
			repo := video.NewRepository(db)
			srv := NewService(repo, &cloudMock{})
			linkable, err := srv.Retrieve(context.Background(), testCase.userID, testCase.id)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err)
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				video, ok := linkable.(*video.Resource)
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

				if video.Bitrate != testCase.expectedBitrate {
					t.Errorf("Invalid bitrate, expected: %d, got: %d\n", testCase.expectedBitrate, video.Bitrate)
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
					t.Errorf("Invalid ration, expected: %d, got: %d\n",
						testCase.expectedRatioX, video.RatioX)
				}

				if video.RatioY != testCase.expectedRatioY {
					t.Errorf("Invalid ration, expected: %d, got: %d\n",
						testCase.expectedRatioY, video.RatioY)
				}

				if video.ServiceID != testCase.expectedServiceID {
					t.Errorf("Invalid service id, expected: %s, got: %s\n",
						testCase.expectedServiceID, video.ServiceID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}

func TestDownloadURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		userID       int64
		videoID      int64
		mock         func()
		expectedURL  string
		errorPresent bool
	}{
		{
			name:         "Video doesn't exists",
			userID:       1,
			videoID:      1,
			errorPresent: true,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", video.TableName)).
					WithArgs(1, 1).
					WillReturnRows(mock.NewRows([]string{"id"}))
			},
		},
		{
			name:         "Invalid request to cloud",
			userID:       1,
			videoID:      1,
			errorPresent: true,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", video.TableName)).
					WithArgs(1, 1).
					WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 1258000, 789569, 700, 600, 4, 3, "error"))
			},
		},
		{
			name:         "Video doesn't exists",
			userID:       1,
			videoID:      1,
			errorPresent: true,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", video.TableName)).
					WithArgs(1, 1).
					WillReturnRows(mock.NewRows([]string{"id"}))
			},
		},
		{
			name:    "Should return url",
			userID:  1,
			videoID: 1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", video.TableName)).
					WithArgs(1, 1).
					WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 1258000, 789569, 700, 600, 4, 3, "https://video.com"))
			},
			expectedURL: "https://video.com",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := video.NewRepository(db)
			srv := NewService(repo, &cloudMock{})

			url, err := srv.DownloadURL(context.Background(), testCase.userID, testCase.videoID)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if url != testCase.expectedURL {
				t.Errorf("Invalid URL, expected: %s, got: %s\n",
					testCase.expectedURL, url)
			}
		})
	}
}
