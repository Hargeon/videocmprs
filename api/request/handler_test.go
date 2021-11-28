package request

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
)

type cloudMock struct{}

func (c *cloudMock) Upload(ctx context.Context, header *multipart.FileHeader) (string, error) {
	if header.Filename == "failed" {
		return "", errors.New("failed connection")
	}

	return "mock_service_id", nil
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
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
					Bitrate:     64000,
					ResolutionX: 800,
					ResolutionY: 600,
					RatioX:      4,
					RatioY:      3,
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
					WithArgs(64000, 800, 600, 4, 3, 1, "test_video.mkv").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, video_name FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "video_name"}).
						AddRow(1, "original_in_review", "", 64000, 800, 600, 4, 3, "test_video.mkv"))

				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", video.TableName)).
					WithArgs("test_video.mkv", 1441786, "mock_service_id").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 1441786, 0, 0, 0, 0, 0, "mock_service_id"))

				mock.ExpectQuery(fmt.Sprintf("UPDATE %s", request.TableName)).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery(fmt.Sprintf("SELECT id, status, details, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, video_name FROM %s", request.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status", "details", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "video_name"}).
						AddRow(1, "original_in_review", "", 64000, 800, 600, 4, 3, "test_video.mkv"))
			},
			expectedBody:   `{"data":{"type":"requests","id":"1","attributes":{"bitrate":64000,"ratio_x":4,"ratio_y":3,"resolution_x":800,"resolution_y":600,"status":"original_in_review","video_name":"test_video.mkv"},"relationships":{"original_video":{"data":{"type":"videos","id":"1"}}},"links":{"self":"need add"}},"included":[{"type":"videos","id":"1","attributes":{"name":"my_name.mkv","size":1441786},"links":{"self":"need add"}}]}` + "\n",
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
					Bitrate:     64000,
					ResolutionX: 800,
					ResolutionY: 600,
					RatioX:      4,
					RatioY:      3,
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

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	h := NewHandler(db, new(cloudMock))

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", int64(1))

		return c.Next()
	})
	app.Get("/requests", h.list)

	cases := []struct {
		name           string
		mock           func()
		requestMock    func() *http.Request
		expectedBody   string
		expectedStatus int
	}{
		{
			name: "Invalid page number",
			mock: func() {},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/requests?page[number]=qwe&page[size]=5", nil)

				return req
			},
			expectedBody:   `{"errors":[{"title":"Invalid page number params"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Invalid page size",
			mock: func() {},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/requests?page[number]=5&page[size]=qwqwe", nil)

				return req
			},
			expectedBody:   `{"errors":[{"title":"Invalid page size params"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Zero requests",
			mock: func() {
				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}))
			},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/requests?page[number]=1&page[size]=10", nil)

				return req
			},
			expectedBody:   `{"data":[]}` + "\n",
			expectedStatus: http.StatusOK,
		},

		{
			name: "Without origin and converted video",
			mock: func() {
				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "", "", 64000, 800, 600, 4, 3, "new_video", nil, nil, nil,
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil))
			},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/requests?page[number]=1&page[size]=10", nil)

				return req
			},
			expectedBody:   `{"data":[{"type":"requests","id":"1","attributes":{"bitrate":64000,"ratio_x":4,"ratio_y":3,"resolution_x":800,"resolution_y":600,"video_name":"new_video"},"links":{"self":"need add"}}]}` + "\n",
			expectedStatus: http.StatusOK,
		},

		{
			name: "With origin video",
			mock: func() {
				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "", "", 64000, 800, 600, 4, 3, "new_video", 1, "new_video", 15000,
						78000, 1200, 800, 6, 5, "new_service_id", nil, nil, nil, nil, nil, nil, nil, nil, nil))
			},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/requests?page[number]=1&page[size]=10", nil)

				return req
			},
			expectedBody:   `{"data":[{"type":"requests","id":"1","attributes":{"bitrate":64000,"ratio_x":4,"ratio_y":3,"resolution_x":800,"resolution_y":600,"video_name":"new_video"},"relationships":{"original_video":{"data":{"type":"videos","id":"1"}}},"links":{"self":"need add"}}],"included":[{"type":"videos","id":"1","attributes":{"bitrate":78000,"name":"new_video","ratio_x":6,"ratio_y":5,"resolution_x":1200,"resolution_y":800,"size":15000},"links":{"self":"need add"}}]}` + "\n",
			expectedStatus: http.StatusOK,
		},

		{
			name: "With origin and converted video",
			mock: func() {
				mock.ExpectQuery("SELECT requests.id, requests.status, requests.details, requests.bitrate, requests.resolution_x, requests.resolution_y, requests.ratio_x, requests.ratio_y, requests.video_name, origin_video.id, origin_video.name, origin_video.size, origin_video.bitrate, origin_video.resolution_x, origin_video.resolution_y, origin_video.ratio_x, origin_video.ratio_y, origin_video.service_id, converted_video.id, converted_video.name, converted_video.size, converted_video.bitrate, converted_video.resolution_x, converted_video.resolution_y, converted_video.ratio_x, converted_video.ratio_y, converted_video.service_id FROM requests LEFT JOIN videos AS origin_video ON requests.original_file_id = origin_video.id LEFT JOIN videos AS converted_video ON requests.converted_file_id = converted_video.id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"requests.id", "requests.status",
						"requests.details", "requests.bitrate", "requests.resolution_x",
						"requests.resolution_y", "requests.ratio_x", "requests.ratio_y",
						"requests.video_name", "origin_video.id", "origin_video.name",
						"origin_video.size", "origin_video.bitrate", "origin_video.resolution_x",
						"origin_video.resolution_y", "origin_video.ratio_x", "origin_video.ratio_y",
						"origin_video.service_id", "converted_video.id", "converted_video.name",
						"converted_video.size", "converted_video.bitrate", "converted_video.resolution_x",
						"converted_video.resolution_y", "converted_video.ratio_x",
						"converted_video.ratio_y", "converted_video.service_id"}).AddRow(
						1, "", "", 64000, 800, 600, 4, 3, "new_video", 1, "new_video", 15000,
						78000, 1200, 800, 6, 5, "new_service_id", 2, "converted_video", 12000, 64000,
						800, 600, 4, 3, "converted_service_id"))
			},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/requests?page[number]=1&page[size]=10", nil)

				return req
			},
			expectedBody:   `{"data":[{"type":"requests","id":"1","attributes":{"bitrate":64000,"ratio_x":4,"ratio_y":3,"resolution_x":800,"resolution_y":600,"video_name":"new_video"},"relationships":{"converted_video":{"data":{"type":"videos","id":"2"}},"original_video":{"data":{"type":"videos","id":"1"}}},"links":{"self":"need add"}}],"included":[{"type":"videos","id":"1","attributes":{"bitrate":78000,"name":"new_video","ratio_x":6,"ratio_y":5,"resolution_x":1200,"resolution_y":800,"size":15000},"links":{"self":"need add"}},{"type":"videos","id":"2","attributes":{"bitrate":64000,"name":"converted_video","ratio_x":4,"ratio_y":3,"resolution_x":800,"resolution_y":600,"size":12000},"links":{"self":"need add"}}]}` + "\n",
			expectedStatus: http.StatusOK,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			req := testCase.requestMock()

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error when creating stub request, error: %s\n", err.Error())
			}

			if resp.StatusCode != testCase.expectedStatus {
				t.Errorf("Invalid status code, expected: %d, got: %d\n",
					testCase.expectedStatus, resp.StatusCode)
			}

			resBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Unexpected error when reading response body, error: %s\n", err.Error())
			}

			if string(resBody) != testCase.expectedBody {
				t.Errorf("Invalid body,\nexpected: %#v\ngot: %#v\n",
					testCase.expectedBody, string(resBody))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s\n", err)
			}
		})
	}
}
