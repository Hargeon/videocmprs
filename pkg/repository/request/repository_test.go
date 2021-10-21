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
		expectedId         int64
		expectedBitrate    int64
		expectedResolution string
		expectedRatio      string
		errorPresent       bool
	}{
		{
			name: "Should add request to db",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs(64000, "800:600", "4:3").
					WillReturnRows(sqlxmock.NewRows([]string{"id", "bitrate", "resolution", "ratio"}).
						AddRow(1, 64000, "800:600", "4:3"))
			},
			req: &Resource{
				Bitrate:    64000,
				Resolution: "800:600",
				Ration:     "4:3",
			},
			expectedId:         1,
			expectedBitrate:    64000,
			expectedResolution: "800:600",
			expectedRatio:      "4:3",
			errorPresent:       false,
		},
		{
			name: "With failed db connection",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", TableName)).
					WithArgs(64000, "800:600", "4:3").
					WillReturnRows(sqlxmock.NewRows([]string{"id", "bitrate", "resolution", "ratio"}))
			},
			req: &Resource{
				Bitrate:    64000,
				Resolution: "800:600",
				Ration:     "4:3",
			},
			expectedId:         0,
			expectedBitrate:    0,
			expectedResolution: "",
			expectedRatio:      "",
			errorPresent:       true,
		},
		{
			name: "With invalid jsonapi.Linkable struct",
			mock: func() {
			},
			req:                &invalidResource{},
			expectedId:         0,
			expectedBitrate:    0,
			expectedResolution: "",
			expectedRatio:      "",
			errorPresent:       true,
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
				res, ok := linkable.(*Resource)
				if !ok {
					t.Fatalf("Invalid type assetrion *reqest.Resource\n")
				}

				if res.Resolution != testCase.expectedResolution {
					t.Errorf("Invalid resoulution, expected: %s, got: %s\n",
						testCase.expectedResolution, res.Resolution)
				}

				if res.Ration != testCase.expectedRatio {
					t.Errorf("Invalid ration, expected: %s, got: %s\n",
						testCase.expectedRatio, res.Ration)
				}

				if res.Bitrate != testCase.expectedBitrate {
					t.Errorf("Invalid bitrate, expected: %d, got: %d\n",
						testCase.expectedBitrate, res.Bitrate)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}
