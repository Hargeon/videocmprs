package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRetrieve(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name          string
		id            int64
		mock          func()
		expectedID    int64
		expectedEmail string
		errorPresent  bool
	}{
		{
			name: "Should find user",
			id:   1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, email FROM %s WHERE", TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow("1", "check@check.com"))
			},
			expectedID:    1,
			expectedEmail: "check@check.com",
			errorPresent:  false,
		},
		{
			name: "Should not find user",
			id:   58,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, email FROM %s WHERE", TableName)).
					WithArgs(58).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))
			},
			expectedID:    0,
			expectedEmail: "",
			errorPresent:  true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewRepository(db)
			linkable, err := repo.Retrieve(context.Background(), testCase.id)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err)
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if !testCase.errorPresent {
				usr, ok := linkable.(*Resource)
				if !ok {
					t.Fatalf("Invalid type assertion.\n")
				}

				if usr.ID != testCase.expectedID {
					t.Errorf("Invalid user id, expected: %d, got: %d\n", testCase.expectedID, usr.ID)
				}

				if usr.Email != testCase.expectedEmail {
					t.Errorf("Invalid user email, expected: %s, got: %s\n", testCase.expectedEmail, usr.Email)
				}
			}
		})
	}
}
