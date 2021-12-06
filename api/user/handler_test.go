package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service/encryption"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"go.uber.org/zap"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	logger := zap.NewExample()
	defer logger.Sync()

	handler := NewHandler(db, logger)
	app := fiber.New()

	app.Post("/", handler.create)

	cases := []struct {
		name           string
		method         string
		user           interface{}
		marshalUser    func(user interface{}) []byte
		mock           func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "With bad unmarshal",
			user: &struct {
				Email                string `json:"email"`
				Password             string `json:"password"`
				PasswordConfirmation string `json:"password_confirmation"`
			}{
				Email:                "check@check.com",
				Password:             "97987897989",
				PasswordConfirmation: "97987897989",
			},
			marshalUser: func(user interface{}) []byte {
				result, err := json.Marshal(user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return result
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"Bad request"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "With invalid email",
			user: &user.Resource{
				Email:                "check",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
			},
			marshalUser: func(user interface{}) []byte {
				var req []byte
				reqBuf := bytes.NewBuffer(req)
				err := jsonapi.MarshalPayload(reqBuf, user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return reqBuf.Bytes()
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"Validation failed"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "With invalid password",
			user: &user.Resource{
				Email:                "check@check.com",
				Password:             "12",
				PasswordConfirmation: "12",
			},
			marshalUser: func(user interface{}) []byte {
				var req []byte
				reqBuf := bytes.NewBuffer(req)
				err := jsonapi.MarshalPayload(reqBuf, user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return reqBuf.Bytes()
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"Validation failed"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "With invalid password confirmation",
			user: &user.Resource{
				Email:                "check@check.com",
				Password:             "123456789",
				PasswordConfirmation: "123456789f",
			},
			marshalUser: func(user interface{}) []byte {
				var req []byte
				reqBuf := bytes.NewBuffer(req)
				err := jsonapi.MarshalPayload(reqBuf, user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return reqBuf.Bytes()
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"Validation failed"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "With valid params",
			user: &user.Resource{
				Email:                "check@check.com",
				Password:             "123456789",
				PasswordConfirmation: "123456789",
			},
			marshalUser: func(user interface{}) []byte {
				var req []byte
				reqBuf := bytes.NewBuffer(req)
				err := jsonapi.MarshalPayload(reqBuf, user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return reqBuf.Bytes()
			},
			mock: func() {
				passHash := encryption.GenerateHash([]byte("123456789"))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", user.TableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", passHash)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, email FROM %s", user.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow("1", "check@check.com"))
			},
			expectedBody:   `{"data":{"type":"users","id":"1","attributes":{"email":"check@check.com"},"links":{"self":"/api/v1/auth/me"}}}` + "\n",
			expectedStatus: http.StatusCreated,
		},
		{
			name: "With failed db connection",
			user: &user.Resource{
				Email:                "check@check.com",
				Password:             "123456789",
				PasswordConfirmation: "123456789",
			},
			marshalUser: func(user interface{}) []byte {
				var req []byte
				reqBuf := bytes.NewBuffer(req)
				err := jsonapi.MarshalPayload(reqBuf, user)
				if err != nil {
					t.Fatalf("Error occured when marshaling user, error: %s\n", err.Error())
				}

				return reqBuf.Bytes()
			},
			mock: func() {
				passHash := encryption.GenerateHash([]byte("123456789"))
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", user.TableName)).
					WithArgs("check@check.com", fmt.Sprintf("%x", passHash)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			expectedBody:   `{"errors":[{"title":"sql: no rows in result set"}]}` + "\n",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			userReqBody := testCase.marshalUser(testCase.user)
			userBuf := bytes.NewBuffer(userReqBody)

			req := httptest.NewRequest(http.MethodPost, "/", userBuf)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error occured when creating stub request, error: %s\n", err.Error())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if resp.StatusCode != testCase.expectedStatus {
				t.Errorf("Invalid status code, expected: %d, got: %d\n", testCase.expectedStatus, resp.StatusCode)
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error occured when reading response body, error: %s\n", err.Error())
			}

			if string(respBody) != testCase.expectedBody {
				t.Errorf("Invalid response body, expected: %#v\n got: %#v\n", testCase.expectedBody, string(respBody))
			}
		})
	}
}
