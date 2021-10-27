package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignIn(t *testing.T) {
	db, mock, err := sqlxmock.Newx()

	handler := NewHandler(db)
	app := fiber.New()

	app.Post("/", handler.signIn)

	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	cases := []struct {
		name           string
		user           interface{}
		marshalUser    func(user interface{}) []byte
		mock           func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "With bad unmarshal",
			user: &struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    "check@check.com",
				Password: "97987897989",
			},
			marshalUser: func(user interface{}) []byte {
				result, err := json.Marshal(user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return result
			},
			mock:           func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":[{"title":"Bad request"}]}` + "\n",
		},
		{
			name: "With invalid email",
			user: &user.Resource{
				Email:    "qwe",
				Password: "qweqweqwe",
			},
			marshalUser: func(user interface{}) []byte {
				var reqBody []byte
				reqBuf := bytes.NewBuffer(reqBody)
				err := jsonapi.MarshalPayload(reqBuf, user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return reqBuf.Bytes()
			},
			mock:           func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":[{"title":"Validation failed"}]}` + "\n",
		},
		{
			name: "With invalid password",
			user: &user.Resource{
				Email:    "check@check.com",
				Password: "qwe",
			},
			marshalUser: func(user interface{}) []byte {
				var reqBody []byte
				reqBuf := bytes.NewBuffer(reqBody)
				err := jsonapi.MarshalPayload(reqBuf, user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return reqBuf.Bytes()
			},
			mock:           func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":[{"title":"Validation failed"}]}` + "\n",
		},
		{
			name: "With valid db params",
			user: &user.Resource{
				Email:    "check@check.com",
				Password: "qweqweqwe",
			},
			marshalUser: func(user interface{}) []byte {
				var reqBody []byte
				reqBuf := bytes.NewBuffer(reqBody)
				err := jsonapi.MarshalPayload(reqBuf, user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return reqBuf.Bytes()
			},
			mock: func() {
				hashPass := encryption.GenerateHash([]byte("qweqweqwe"))
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", user.UserTableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", hashPass)).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "With invalid db params",
			user: &user.Resource{
				Email:    "check@check.com",
				Password: "qweqweqwe",
			},
			marshalUser: func(user interface{}) []byte {
				var reqBody []byte
				reqBuf := bytes.NewBuffer(reqBody)
				err := jsonapi.MarshalPayload(reqBuf, user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return reqBuf.Bytes()
			},
			mock: func() {
				hashPass := encryption.GenerateHash([]byte("qweqweqwe"))
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", user.UserTableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", hashPass)).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":[{"title":"sql: no rows in result set"}]}` + "\n",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			userReqBody := testCase.marshalUser(testCase.user)
			userReqBuf := bytes.NewBuffer(userReqBody)
			req := httptest.NewRequest(http.MethodPost, "/", userReqBuf)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error when creating a stub request to handler, error: %s\n", err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if resp.StatusCode != testCase.expectedStatus {
				t.Errorf("Invalid status, expected: %d, got: %d\n", testCase.expectedStatus, resp.StatusCode)
			}

			if testCase.expectedStatus != http.StatusCreated {
				respBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Error occured when reading response body, error: %s\n", err.Error())
				}

				if string(respBody) != testCase.expectedBody {
					t.Errorf("Invalid response body,\nexpected: %#v,\ngot: %#v\n", testCase.expectedBody,
						string(respBody))
				}
			}
		})
	}
}
