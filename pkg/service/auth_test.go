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
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	now := time.Now()

	cases := []struct {
		name         string
		user         *model.User
		mock         func()
		expectedId   int64
		errorPresent bool
	}{
		{
			name: "With valid email, password, created_at",
			user: &model.User{
				Email:     "check@gmail.com",
				Password:  "123456789",
				CreatedAt: now,
			},
			mock: func() {
				passHash := generateHash([]byte("123456789"), []byte(os.Getenv("DB_SECRET")))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", model.UserTableName)).
					WithArgs("check@gmail.com", fmt.Sprintf("%x", passHash), now).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(1))
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
			mock: func() {
				passHash := generateHash([]byte("123456789"), []byte(os.Getenv("DB_SECRET")))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", model.UserTableName)).
					WithArgs("", fmt.Sprintf("%x", passHash), now).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
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
			mock: func() {
				passHash := generateHash([]byte(""), []byte(os.Getenv("DB_SECRET")))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", model.UserTableName)).
					WithArgs("check@gmail.com", fmt.Sprintf("%x", passHash), now).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			expectedId:   0,
			errorPresent: true,
		},
		{
			name: "With invalid created_at",
			user: &model.User{
				Email:     "check@gmail.com",
				Password:  "123456789",
				CreatedAt: now,
			},
			mock: func() {
				passHash := generateHash([]byte("123456789"), []byte(os.Getenv("DB_SECRET")))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", model.UserTableName)).
					WithArgs("check@gmail.com", fmt.Sprintf("%x", passHash), now).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			expectedId:   0,
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
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
				hashPassword := generateHash([]byte("pokopkopkpo"), []byte(os.Getenv("DB_SECRET")))
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", model.UserTableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", hashPassword)).
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
				hashPassword := generateHash([]byte("pokopkopkpo"), []byte(os.Getenv("DB_SECRET")))
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", model.UserTableName)).
					WithArgs("", fmt.Sprintf("%x", hashPassword)).
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
