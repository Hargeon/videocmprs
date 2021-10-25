package video

import (
	"context"
	"fmt"
	"github.com/google/jsonapi"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

type invalidLinkable struct{}

func (r *invalidLinkable) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": "",
	}
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		video        jsonapi.Linkable
		mock         func()
		expectedId   int64
		expectedSize int64
		expectedName string
		errorPresent bool
	}{
		{
			name: "Should add video",
			video: &Resource{
				Name:    "my_name.mkv",
				Size:    1258000,
				Bitrate: 64000,
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs("my_name.mkv", 1258000, 64000).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution, ratio, service_id FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution", "ratio", "service_id"}).
						AddRow(1, "my_name.mkv", 1258000, 0, "", "", ""))
			},
			expectedId:   1,
			expectedSize: 1258000,
			expectedName: "my_name.mkv",
			errorPresent: false,
		},
		{
			name: "Should not add video",
			video: &Resource{
				Name:    "qwe",
				Size:    125,
				Bitrate: 78787,
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs("qwe", 125, 78787).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			expectedId:   0,
			expectedSize: 0,
			expectedName: "",
			errorPresent: true,
		},
		{
			name:         "With invalid jsonapi.Linkable",
			video:        &invalidLinkable{},
			mock:         func() {},
			expectedId:   0,
			expectedSize: 0,
			expectedName: "",
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewRepository(db)
			linkable, err := repo.Create(context.Background(), testCase.video)
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

				if video.ID != testCase.expectedId {
					t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedId, video.ID)
				}

				if video.Name != testCase.expectedName {
					t.Errorf("Invalid name, expected: %s, got: %s\n", testCase.expectedName, video.Name)
				}

				if video.Size != testCase.expectedSize {
					t.Errorf("Invalid size, expected: %d, got: %d\n", testCase.expectedSize, video.Size)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}

func TestRetrieve(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name               string
		id                 int64
		mock               func()
		expectedId         int64
		expectedSize       int64
		expectedBitrate    int64
		expectedName       string
		expectedResolution string
		expectedRation     string
		expectedServiceId  string
		errorPresent       bool
	}{
		{
			name: "Should find video",
			id:   1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution, ratio, service_id FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution", "ratio", "service_id"}).
						AddRow(1, "check_name.mkv", 185000, 789569, "700:600", "4:3", "qweqweqwqfqw"))
			},
			expectedId:         1,
			expectedName:       "check_name.mkv",
			expectedSize:       185000,
			expectedBitrate:    789569,
			expectedResolution: "700:600",
			expectedRation:     "4:3",
			expectedServiceId:  "qweqweqwqfqw",
			errorPresent:       false,
		},
		{
			name: "Should not find video",
			id:   1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution, ratio, service_id FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution", "ratio", "service_id"}))
			},
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewRepository(db)
			linkable, err := repo.Retrieve(context.Background(), testCase.id)
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

				if video.ID != testCase.expectedId {
					t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedId, video.ID)
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

				if video.Resolution != testCase.expectedResolution {
					t.Errorf("Invalid resolution, expected: %s, got: %s\n",
						testCase.expectedResolution, video.Resolution)
				}

				if video.Ratio != testCase.expectedRation {
					t.Errorf("Invalid ration, expected: %s, got: %s\n",
						testCase.expectedRation, video.Ratio)
				}

				if video.ServiceId != testCase.expectedServiceId {
					t.Errorf("Invalid service id, expected: %s, got: %s\n",
						testCase.expectedServiceId, video.ServiceId)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}