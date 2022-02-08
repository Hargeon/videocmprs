package user

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name           string
		email          string
		mock           func()
		expectedResult bool
		errorPresent   bool
	}{
		{
			name:  "User exists in db",
			email: "check@gmail.com",
			mock: func() {
				mock.ExpectQuery("SELECT count").
					WithArgs("check@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))
			},
			expectedResult: false,
			errorPresent:   false,
		},
		{
			name:  "Sql error",
			email: "not_exists@gmail.com",
			mock: func() {
				mock.ExpectQuery("SELECT count").
					WithArgs("not_exists@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"total"}))
			},
			expectedResult: true,
			errorPresent:   true,
		},
		{
			name:  "User is not exists in db",
			email: "check@gmail.com",
			mock: func() {
				mock.ExpectQuery("SELECT count").
					WithArgs("check@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(0))
			},
			expectedResult: true,
			errorPresent:   false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewRepository(db)
			result, err := repo.Unique(context.Background(), testCase.email)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error, %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if result != testCase.expectedResult {
				t.Errorf("invalid result, expected: %v, got: %v\n",
					testCase.expectedResult, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
