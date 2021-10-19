package auth

import (
	"context"
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/repository/auth"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestExists(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		user         *user.Resource
		mock         func()
		errorPresent bool
		tokenPresent bool
	}{
		{
			name: "Should find user",
			user: &user.Resource{
				Email:    "check@check.com",
				Password: "qweqweqwe",
			},
			mock: func() {
				hashPass := encryption.GenerateHash([]byte("qweqweqwe"))
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", user.UserTableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", hashPass)).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(1))
			},
			errorPresent: false,
			tokenPresent: true,
		},
		{
			name: "Should not find user",
			user: &user.Resource{
				Email:    "check2@check.com",
				Password: "qweqweqwe",
			},
			mock: func() {
				hashPass := encryption.GenerateHash([]byte("qweqweqwe"))
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", user.UserTableName)).
					WithArgs("check2@check.com", fmt.Sprintf("%x", hashPass)).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			errorPresent: true,
			tokenPresent: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := auth.NewRepository(db)
			srv := NewService(repo)
			linkable, err := srv.GenerateToken(context.Background(), testCase.user)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				usr, ok := linkable.(*user.Resource)
				if !ok {
					t.Fatalf("Can't type assertion for auth.Resource\n")
				}

				tokenPresent := usr.Token > ""

				if tokenPresent != testCase.tokenPresent {
					t.Errorf("Invalid token\n")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}

		})
	}
}
