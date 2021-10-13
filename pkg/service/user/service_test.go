package user

import (
	"context"
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"
	"github.com/google/jsonapi"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name          string
		usr           jsonapi.Linkable
		password      string
		mock          func()
		expectedId    int64
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
				passHash := encryption.GenerateHash([]byte("qjwpeqwpoekpqwe"))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", user.UserTableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", passHash)).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, email FROM %s", user.UserTableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "email"}).AddRow("1", "check@check.com"))
			},
			expectedId:    1,
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
				passHash := encryption.GenerateHash([]byte("qjwpeqwpoekpqwe"))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", user.UserTableName)).
					WithArgs("", fmt.Sprintf("%x", passHash)).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			expectedId:    0,
			expectedEmail: "",
			errorPresent:  true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
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

				if usr.ID != testCase.expectedId {
					t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedId, usr.ID)
				}

				if usr.Email != testCase.expectedEmail {
					t.Errorf("Invalid email, expected: %s, got: %s\n", testCase.expectedEmail, usr.Email)
				}
			}
		})
	}
}
