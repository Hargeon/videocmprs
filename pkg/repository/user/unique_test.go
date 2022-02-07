package user

import (
	"context"
	"fmt"
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
	}{
		{
			name:  "User exists in db",
			email: "check@gmail.com",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", TableName)).
					WithArgs("check@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
			},
			expectedResult: false,
		},
		{
			name:  "User is not exists in db",
			email: "not_exists@gmail.com",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", TableName)).
					WithArgs("not_exists@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectedResult: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewRepository(db)
			result := repo.Unique(context.Background(), testCase.email)
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
