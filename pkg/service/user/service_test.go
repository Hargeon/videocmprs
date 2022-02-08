package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/jsonapi"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name          string
		usr           jsonapi.Linkable
		password      string
		mock          func()
		expectedID    int64
		expectedEmail string
		errorPresent  bool
	}{
		{
			name: "With valid email and password",
			usr: &user.Resource{
				Email:    "check@check.com",
				Password: "qjwpeqwpoekpqwe",
			},
			mock: func() {
				mock.ExpectQuery("SELECT count").
					WithArgs("check@check.com").
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(0))

				passHash := encryption.GenerateHash([]byte("qjwpeqwpoekpqwe"))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", user.TableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", passHash)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, email FROM %s", user.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow("1", "check@check.com"))
			},
			expectedID:    1,
			expectedEmail: "check@check.com",
			errorPresent:  false,
		},
		{
			name: "With invalid email",
			usr: &user.Resource{
				Email:    "",
				Password: "qjwpeqwpoekpqwe",
			},
			mock: func() {
				mock.ExpectQuery("SELECT count").
					WithArgs("").
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(0))

				passHash := encryption.GenerateHash([]byte("qjwpeqwpoekpqwe"))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", user.TableName)).
					WithArgs("", fmt.Sprintf("%x", passHash)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectedID:    0,
			expectedEmail: "",
			errorPresent:  true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			testCase.mock()
			repo := user.NewRepository(db)
			srv := NewService(repo)
			linkable, err := srv.Create(context.Background(), testCase.usr)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err)
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}

			if !testCase.errorPresent {
				usr, ok := linkable.(*user.Resource)
				if !ok {
					t.Fatalf("Invalid assertion\n")
				}

				if usr.ID != testCase.expectedID {
					t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedID, usr.ID)
				}

				if usr.Email != testCase.expectedEmail {
					t.Errorf("Invalid email, expected: %s, got: %s\n", testCase.expectedEmail, usr.Email)
				}
			}
		})
	}
}
