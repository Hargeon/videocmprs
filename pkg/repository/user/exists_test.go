package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		email        string
		password     string
		mock         func()
		expectedId   int64
		errorPresent bool
	}{
		{
			name:     "User exists",
			email:    "check@check.com",
			password: "qweqweqweqwe",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", TableName)).
					WithArgs("check@check.com", "qweqweqweqwe").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedId:   1,
			errorPresent: false,
		},
		{
			name:     "User doesn't exists",
			email:    "check@check.com",
			password: "qweqweqweqwe",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", TableName)).
					WithArgs("check@check.com", "qweqweqweqwe").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectedId:   0,
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewRepository(db)
			id, err := repo.Exists(context.Background(), testCase.email, testCase.password)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error, error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if id != testCase.expectedId {
				t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedId, id)
			}
		})
	}
}
