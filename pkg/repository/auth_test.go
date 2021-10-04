package repository

import (
	"fmt"
	"github.com/Hargeon/videocmprs/db/model"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

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
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", model.UserTableName)).
					WithArgs("check@gmail.com", "123456789", now).
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
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", model.UserTableName)).
					WithArgs("", "123456789", "").
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
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", model.UserTableName)).
					WithArgs("check@gmail.com", "", now).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			expectedId:   0,
			errorPresent: true,
		},
		{
			name: "With valid created_at",
			user: &model.User{
				Email:    "check@gmail.com",
				Password: "123456789",
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", model.UserTableName)).
					WithArgs("check@gmail.com", "123456789").
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			expectedId:   0,
			errorPresent: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			repo := NewRepository(db)
			id, err := repo.CreateUser(testCase.user)
			if err != nil && !testCase.errorPresent {
				t.Errorf("Unexpected error, error: %s\n", err)
			}

			if err == nil && testCase.errorPresent {
				t.Errorf("Should be error")
			}
			if id != testCase.expectedId {
				t.Errorf("Invalid id, expected: %d, got: %d\n", testCase.expectedId, id)
			}
		})
	}
}
