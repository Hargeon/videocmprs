package request

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"
)

type cloudMock struct{}

func (c *cloudMock) Upload(ctx context.Context, header *multipart.FileHeader) (string, error) {
	if header.Filename == "failed" {
		return "", errors.New("failed connection")
	}
	return "mock_service_id", nil
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	h := NewHandler(db, new(cloudMock))

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(1))
		return c.Next()
	})
	app.Post("/", h.create)

	cases := []struct {
		name        string
		requestMock func() *http.Request
		mock        func()

		expectedBody   string
		expectedStatus int
	}{
		{
			name: "Without file",
			requestMock: func() *http.Request {
				buf := new(bytes.Buffer)
				writer := multipart.NewWriter(buf)
				if err := writer.WriteField("test", "failed"); err != nil {
					t.Fatalf("Unexpected error when adding test filed, error: %s\n", err.Error())
				}
				if err := writer.Close(); err != nil {
					t.Fatalf("Unexpected error when closing writter, error: %s\n", err.Error())
				}

				req := httptest.NewRequest(http.MethodPost, "/", buf)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"Request does not include file"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid content-type",
			requestMock: func() *http.Request {
				buf := new(bytes.Buffer)
				writer := multipart.NewWriter(buf)
				part, err := writer.CreateFormFile("video", "test_name")
				if err != nil {
					t.Fatalf("Unexpected error when adding file to request, error: %s\n", err.Error())
				}

				fMock := bytes.NewReader([]byte("qwertyuiopasdfghjkl"))
				if _, err = io.Copy(part, fMock); err != nil {
					t.Fatalf("Unexpected error when copying body")
				}

				if err := writer.Close(); err != nil {
					t.Fatalf("Unexpected error when closing writter, error: %s\n", err.Error())
				}

				req := httptest.NewRequest(http.MethodPost, "/", buf)
				return req
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"Request does not include file"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid file body",
			requestMock: func() *http.Request {
				buf := new(bytes.Buffer)
				writer := multipart.NewWriter(buf)
				part, err := writer.CreateFormFile("video", "test_name")
				if err != nil {
					t.Fatalf("Unexpected error when adding file to request, error: %s\n", err.Error())
				}

				fMock := bytes.NewReader([]byte("qwertyuiopasdfghjkl"))
				if _, err = io.Copy(part, fMock); err != nil {
					t.Fatalf("Unexpected error when copying body")
				}

				if err := writer.Close(); err != nil {
					t.Fatalf("Unexpected error when closing writter, error: %s\n", err.Error())
				}

				req := httptest.NewRequest(http.MethodPost, "/", buf)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"File is not a video"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid file type",
			requestMock: func() *http.Request {
				buf := new(bytes.Buffer)
				writer := multipart.NewWriter(buf)
				part, err := writer.CreateFormFile("video", "test_image.jpeg")
				if err != nil {
					t.Fatalf("Unexpected error when adding file to request, error: %s\n", err.Error())
				}

				fMock, err := os.Open("test_image.jpeg")
				if err != nil {
					t.Fatalf("Unexpected error while readeing image, error: %s\n", err.Error())
				}

				if _, err = io.Copy(part, fMock); err != nil {
					t.Fatalf("Unexpected error when copying body")
				}

				if err := writer.Close(); err != nil {
					t.Fatalf("Unexpected error when closing writter, error: %s\n", err.Error())
				}

				req := httptest.NewRequest(http.MethodPost, "/", buf)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"File is not a video"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid request params",
			requestMock: func() *http.Request {
				buf := new(bytes.Buffer)
				writer := multipart.NewWriter(buf)
				fMock, err := os.Open("test_video.mkv")
				if err != nil {
					t.Fatalf("Unexpected error while readeing image, error: %s\n", err.Error())
				}

				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="video"; filename="%s"`, fMock.Name()))
				h.Set("Content-Type", "video/x-qwdq")
				part, err := writer.CreatePart(h)
				if err != nil {
					t.Fatalf("Unexpected error when adding file to request, error: %s\n", err.Error())
				}

				if _, err = io.Copy(part, fMock); err != nil {
					t.Fatalf("Unexpected error when copying body")
				}

				if err = writer.WriteField("requests", "qweqwe"); err != nil {
					t.Fatalf("Unexpected error while adding request, error: %s\n", err.Error())
				}

				if err := writer.Close(); err != nil {
					t.Errorf("Unexpected error when closing writter, error: %s\n", err.Error())
				}

				req := httptest.NewRequest(http.MethodPost, "/", buf)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"Invalid request params"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid request validation",
			requestMock: func() *http.Request {
				buf := new(bytes.Buffer)
				writer := multipart.NewWriter(buf)
				fMock, err := os.Open("test_video.mkv")
				if err != nil {
					t.Fatalf("Unexpected error while readeing image, error: %s\n", err.Error())
				}

				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="video"; filename="%s"`, fMock.Name()))
				h.Set("Content-Type", "video/x-qwdq")
				part, err := writer.CreatePart(h)
				if err != nil {
					t.Fatalf("Unexpected error when adding file to request, error: %s\n", err.Error())
				}

				if _, err = io.Copy(part, fMock); err != nil {
					t.Fatalf("Unexpected error when copying body")
				}

				r := &request.Resource{}

				bufReq := new(bytes.Buffer)

				if err = jsonapi.MarshalPayload(bufReq, r); err != nil {
					t.Fatalf("Unexpected error when marchaling request, error: %s\n", err.Error())
				}

				if err = writer.WriteField("requests", bufReq.String()); err != nil {
					t.Fatalf("Unexpected error while adding request, error: %s\n", err.Error())
				}

				if err := writer.Close(); err != nil {
					t.Errorf("Unexpected error when closing writter, error: %s\n", err.Error())
				}

				req := httptest.NewRequest(http.MethodPost, "/", buf)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"Validation failed"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Valid params",
			requestMock: func() *http.Request {
				buf := new(bytes.Buffer)
				writer := multipart.NewWriter(buf)
				fMock, err := os.Open("test_video.mkv")
				if err != nil {
					t.Fatalf("Unexpected error while readeing image, error: %s\n", err.Error())
				}

				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="video"; filename="%s"`, fMock.Name()))
				h.Set("Content-Type", "video/x-qwdq")
				part, err := writer.CreatePart(h)
				if err != nil {
					t.Fatalf("Unexpected error when adding file to request, error: %s\n", err.Error())
				}

				if _, err = io.Copy(part, fMock); err != nil {
					t.Fatalf("Unexpected error when copying body")
				}

				r := &request.Resource{
					Bitrate:    64000,
					Resolution: "800:600",
					Ratio:      "4:3",
				}

				bufReq := new(bytes.Buffer)

				if err = jsonapi.MarshalPayload(bufReq, r); err != nil {
					t.Fatalf("Unexpected error when marchaling request, error: %s\n", err.Error())
				}

				if err = writer.WriteField("requests", bufReq.String()); err != nil {
					t.Fatalf("Unexpected error while adding request, error: %s\n", err.Error())
				}

				if err := writer.Close(); err != nil {
					t.Errorf("Unexpected error when closing writter, error: %s\n", err.Error())
				}

				req := httptest.NewRequest(http.MethodPost, "/", buf)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", request.TableName)).
					WithArgs(64000, "800:600", "4:3", 1).
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution, ratio FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution", "ratio"}).
						AddRow(1, "original_in_review", "", 64000, "800:600", "4:3"))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs("test_video.mkv", 1441786, "mock_service_id").
					WillReturnRows(sqlxmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution, ratio, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlxmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution", "ratio", "service_id"}).
						AddRow(1, "my_name.mkv", 1441786, 0, "", "", "mock_service_id"))
			},
			expectedBody:   `{"data":{"type":"requests","id":"1","attributes":{"bitrate":64000,"ratio":"4:3","resolution":"800:600","status":"original_in_review"},"links":{"self":"need add"}}}` + "\n",
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid db connection",
			requestMock: func() *http.Request {
				buf := new(bytes.Buffer)
				writer := multipart.NewWriter(buf)
				fMock, err := os.Open("test_video.mkv")
				if err != nil {
					t.Fatalf("Unexpected error while readeing image, error: %s\n", err.Error())
				}

				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="video"; filename="%s"`, fMock.Name()))
				h.Set("Content-Type", "video/x-qwdq")
				part, err := writer.CreatePart(h)
				if err != nil {
					t.Fatalf("Unexpected error when adding file to request, error: %s\n", err.Error())
				}

				if _, err = io.Copy(part, fMock); err != nil {
					t.Fatalf("Unexpected error when copying body")
				}

				r := &request.Resource{
					Bitrate:    64000,
					Resolution: "800:600",
					Ratio:      "4:3",
				}

				bufReq := new(bytes.Buffer)

				if err = jsonapi.MarshalPayload(bufReq, r); err != nil {
					t.Fatalf("Unexpected error when marchaling request, error: %s\n", err.Error())
				}

				if err = writer.WriteField("requests", bufReq.String()); err != nil {
					t.Fatalf("Unexpected error while adding request, error: %s\n", err.Error())
				}

				if err := writer.Close(); err != nil {
					t.Errorf("Unexpected error when closing writter, error: %s\n", err.Error())
				}

				req := httptest.NewRequest(http.MethodPost, "/", buf)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mock:           func() {},
			expectedBody:   `{"errors":[{"title":"Can not create request"}]}` + "\n",
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
				t.Errorf("Invalid status code, expected: %d, got: %d\n",
					testCase.expectedStatus, res.StatusCode)
			}

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Unexpected error when reading response body, error: %s\n", err.Error())
			}

			strBody := string(resBody)
			if strBody != testCase.expectedBody {
				t.Errorf("Invalid body, expected: %#v, got: %#v\n",
					testCase.expectedBody, strBody)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	cases := []struct {
		name           string
		types          []string
		expectedResult bool
	}{
		{
			name:           "Valid video",
			types:          []string{"video/mp4"},
			expectedResult: true,
		},
		{
			name:           "Invalid video",
			types:          []string{"image/jpeg"},
			expectedResult: false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			h := new(Handler)
			result := h.isFile(testCase.types)
			if result != testCase.expectedResult {
				t.Errorf("Invalid result, expectedL %v, got: %v\n",
					testCase.expectedResult, result)
			}
		})
	}
}
