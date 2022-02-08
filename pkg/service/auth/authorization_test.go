package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGenerateToken(t *testing.T) {
	db, mock, err := sqlmock.New()
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
				mock.ExpectQuery("SELECT count").
					WithArgs("check@check.com").
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))

				hashPass := encryption.GenerateHash([]byte("qweqweqwe"))
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", user.TableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", hashPass)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
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
				mock.ExpectQuery("SELECT count").
					WithArgs("check2@check.com").
					WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(0))
			},
			errorPresent: true,
			tokenPresent: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := user.NewRepository(db)
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
				mock.ExpectQuery(fmt.Sprintf("SELECT id, email FROM %s", user.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow("1", "check@check.com"))
			},
			expectedID:    1,
			expectedEmail: "check@check.com",
			errorPresent:  false,
		},

		{
			name: "Should not find user",
			id:   1,
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, email FROM %s", user.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))
			},
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := user.NewRepository(db)
			srv := NewService(repo)
			usrLinkable, err := srv.Retrieve(context.Background(), testCase.id)

			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error: %s\n", err.Error())
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error\n")
			}

			if err == nil {
				usr, ok := usrLinkable.(*user.Resource)
				if !ok {
					t.Fatalf("Invalid type assertion *user.Resource\n")
				}

				if usr.ID != testCase.expectedID {
					t.Errorf("Invalid id, expected: %d, got: %d\n",
						testCase.expectedID, usr.ID)
				}

				if usr.Email != testCase.expectedEmail {
					t.Errorf("Invalid email, expected: %s, got: %s\n",
						testCase.expectedEmail, usr.Email)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}
