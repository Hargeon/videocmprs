package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
)

func TestSignIn(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	handler := NewHandler(db)
	app := fiber.New()

	app.Post("/", handler.signIn)

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
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", user.TableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", hashPass)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
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
				mock.ExpectQuery(fmt.Sprintf("SELECT id FROM %s", user.TableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", hashPass)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
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

func TestRetrieve(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	handler := NewHandler(db)
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(1))

		return c.Next()
	})
	app.Get("/", handler.retrieve)

	cases := []struct {
		name           string
		mock           func()
		requestMock    func() *http.Request
		expectedBody   string
		expectedStatus int
	}{
		{
			name: "Should find user",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, email FROM %s", user.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow("1", "check@check.com"))
			},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)

				return req
			},
			expectedBody:   `{"data":{"type":"users","id":"1","attributes":{"email":"check@check.com"},"links":{"self":"/api/v1/auth/me"}}}` + "\n",
			expectedStatus: http.StatusOK,
		},

		{
			name: "Should not find user",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, email FROM %s", user.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))
			},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)

				return req
			},
			expectedBody:   `{"errors":[{"title":"sql: no rows in result set"}]}` + "\n",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			req := testCase.requestMock()

			res, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error when creating a stub request, error: %s\n", err.Error())
			}

			if res.StatusCode != testCase.expectedStatus {
				t.Errorf("Invaid status code, expected: %d, got: %d\n",
					testCase.expectedStatus, res.StatusCode)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Unexpected error when reading a response body, error: %s\n", err.Error())
			}

			if string(body) != testCase.expectedBody {
				t.Errorf("Invalid body,\nexpected: %#v\ngot: %#v\n",
					testCase.expectedBody, string(body))
			}
		})
	}
}
