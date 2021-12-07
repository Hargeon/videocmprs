package video

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRelationExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		userId       int64
		videoID      int64
		mock         func()
		expectedID   int64
		errorPresent bool
	}{
		{
			name:    "Invalid db connection",
			userId:  1,
			videoID: 1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", TableName)).
					WithArgs(1, 1).
					WillReturnRows(mock.NewRows([]string{"id"}))
			},
			errorPresent: true,
		},
		{
			name:    "Should return id",
			userId:  1,
			videoID: 1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", TableName)).
					WithArgs(1, 1).
					WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedID: 1,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewRepository(db)
			id, err := repo.RelationExists(context.Background(), testCase.userId, testCase.videoID)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if id != testCase.expectedID {
				t.Errorf("Invalid ID, expected: %d, got: %d\n",
					testCase.expectedID, id)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
