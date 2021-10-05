package service

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/Hargeon/videocmprs/db/model"
	"github.com/Hargeon/videocmprs/pkg/repository"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"os"
	"testing"
	"time"
)

func TestGenerateHash(t *testing.T) {
	cases := []struct {
		name     string
		password []byte
		salt     []byte
		output   []byte
	}{
		{
			name:     "Hashing",
			password: []byte{15, 78, 21, 69, 89, 32},
			salt:     []byte{98, 12, 45, 89, 125, 36, 78, 12, 45, 45},
		},
	}

	for i, testCase := range cases {
		hash := sha1.New()
		hash.Write(testCase.password)
		output := hash.Sum(testCase.salt)
		testCase.output = output
		cases[i] = testCase
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			output := generateHash(testCase.password, testCase.salt)
			if !bytes.Equal(output, testCase.output) {
				t.Errorf("invalid hashing, expected: %x, got: %x\n", testCase.output, output)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	type caseTest struct {
		name         string
		user         *model.User
		mock         func(c caseTest)
		expectedId   int64
		errorPresent bool
	}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	now := time.Now()

	cases := []caseTest{
		{
			name: "With valid email, password, created_at",
			user: &model.User{
				Email:     "check@gmail.com",
				Password:  "123456789",
				CreatedAt: now,
			},
			expectedId:   1,
			errorPresent: false,
		},
		{
			name: "With invalid email",
			user: &model.User{
				Email:     "",
				Password:  "123456789",
				CreatedAt: now,
			},
			expectedId:   0,
			errorPresent: true,
		},
		{
			name: "With invalid password",
			user: &model.User{
				Email:     "check@gmail.com",
				Password:  "",
				CreatedAt: now,
			},
			expectedId:   0,
			errorPresent: true,
		},
		{
			name: "With valid created_at",
			user: &model.User{
				Email:     "check@gmail.com",
				Password:  "123456789",
				CreatedAt: now,
			},
			expectedId:   0,
			errorPresent: true,
		},
	}

	for i, testCase := range cases {
		testMock := func(c caseTest) {
			passHash := generateHash([]byte(c.user.Password), []byte(os.Getenv("DB_SECRET")))
			rows := sqlxmock.NewRows([]string{"id"})
			if !c.errorPresent {
				rows = rows.AddRow(1)
			}
			mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", model.UserTableName)).
				WithArgs(c.user.Email, fmt.Sprintf("%x", passHash), c.user.CreatedAt).
				WillReturnRows(rows)
		}
		testCase.mock = testMock
		cases[i] = testCase
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock(testCase)
			repo := repository.NewRepository(db)
			srv := NewService(repo)
			id, err := srv.Authorization.CreateUser(testCase.user)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error, error: %s\n", err)
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error")
			}
			if id != testCase.expectedId {
				t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedId, id)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name         string
		email        string
		password     string
		mock         func()
		tokenPresent bool
		errorPresent bool
	}{
		{
			name:     "Should find user",
			email:    "check@check.com",
			password: "pokopkopkpo",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", model.UserTableName)).
					WithArgs("check@check.com", "pokopkopkpo").
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(1))
			},
			tokenPresent: true,
			errorPresent: false,
		},
		{
			name:     "Should not find user",
			email:    "",
			password: "pokopkopkpo",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", model.UserTableName)).
					WithArgs("", "pokopkopkpo").
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			tokenPresent: false,
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := repository.NewRepository(db)
			srv := NewAuthService(repo)
			token, err := srv.GenerateToken(testCase.email, testCase.password)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error, error: %s\n", err)
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error")
			}
			if token != "" && !testCase.tokenPresent {
				t.Errorf("Invalid token. expected empty token, got: %s\n", token)
			}

			if token == "" && testCase.tokenPresent {
				t.Errorf("Invalid token. token should be not empty")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
