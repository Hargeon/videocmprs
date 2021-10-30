package request

import (
	"context"
	"fmt"
	"github.com/google/jsonapi"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

type invalidResource struct{}

func (r *invalidResource) JSONAPILinks() *jsonapi.Links {
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
		name               string
		mock               func()
		req                jsonapi.Linkable
		expectedID         int64
		expectedStatus     string
		expectedDetails    string
		expectedBitrate    int64
		expectedResolution string
		expectedRatio      string
		errorPresent       bool
	}{
		{
			name: "Should add request to db",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs(64000, "800:600", "4:3", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "original_in_review", "", 64000, "800:600", "4:3"))
			},
			req: &Resource{
				UserID:     1,
				Bitrate:    64000,
				Resolution: "800:600",
				Ratio:      "4:3",
			},
			expectedID:         1,
			expectedStatus:     "original_in_review",
			expectedDetails:    "",
			expectedBitrate:    64000,
			expectedResolution: "800:600",
			expectedRatio:      "4:3",
			errorPresent:       false,
		},
		{
			name: "Should not add request to db",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs(64000, "800:600", "4:3", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			req: &Resource{
				UserID:     1,
				Bitrate:    64000,
				Resolution: "800:600",
				Ratio:      "4:3",
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

				if request.Resolution != testCase.expectedResolution {
					t.Errorf("Invalid resolution, expected: %s, got: %s\n",
						testCase.expectedResolution, request.Resolution)
				}

				if request.Ratio != testCase.expectedRatio {
					t.Errorf("Invalid ratio, expected: %s, got: %s\n",
						testCase.expectedRatio, request.Ratio)
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
		expectedID         int64
		expectedStatus     string
		expectedDetails    string
		expectedBitrate    int64
		expectedResolution string
		expectedRatio      string
		errorPresent       bool
	}{
		{
			name: "Should return request",
			id:   1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "original_in_review", "", 1589875, "800:600", "4:3"))
			},
			expectedID:         1,
			expectedStatus:     "original_in_review",
			expectedDetails:    "",
			expectedBitrate:    1589875,
			expectedResolution: "800:600",
			expectedRatio:      "4:3",
			errorPresent:       false,
		},
		{
			name: "Should not return request",
			id:   1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}))
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

				if request.Resolution != testCase.expectedResolution {
					t.Errorf("Invalid resolution, expected: %s, got: %s\n",
						testCase.expectedResolution, request.Resolution)
				}

				if request.Ratio != testCase.expectedRatio {
					t.Errorf("Invalid ratio, expected: %s, got: %s\n",
						testCase.expectedRatio, request.Ratio)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name               string
		id                 int64
		fields             map[string]interface{}
		mock               func()
		expectedID         int64
		expectedStatus     string
		expectedDetails    string
		expectedBitrate    int64
		expectedResolution string
		expectedRatio      string
		errorPresent       bool
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
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "failed", "Can't add video to database", 64000, "800:600", "4:3"))
			},
			expectedID:         1,
			expectedStatus:     "failed",
			expectedDetails:    "Can't add video to database",
			expectedBitrate:    64000,
			expectedResolution: "800:600",
			expectedRatio:      "4:3",
			errorPresent:       false,
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
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
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

				if request.Resolution != testCase.expectedResolution {
					t.Errorf("Invalid resolution, expected: %s, got: %s\n",
						testCase.expectedResolution, request.Resolution)
				}

				if request.Ratio != testCase.expectedRatio {
					t.Errorf("Invalid ratio, expected: %s, got: %s\n",
						testCase.expectedRatio, request.Ratio)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}
