package video

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/repository/video"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
)

func TestRetrieve(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	h := NewHandler(db)
	app := fiber.New()
	app.Mount("/videos", h.InitRoutes())

	cases := []struct {
		name           string
		mock           func()
		requestMock    func() *http.Request
		expectedBody   string
		expectedStatus int
	}{
		{
			name: "Invalid id",
			mock: func() {},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/videos/name", nil)

				return req
			},
			expectedBody:   `{"errors":[{"title":"Invalid ID"}]}` + "\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid db connection",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}))
			},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/videos/1", nil)

				return req
			},
			expectedBody:   `{"errors":[{"title":"Can not fetch video"}]}` + "\n",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Valid db connection",
			mock: func() {
				mock.ExpectQuery(fmt.Sprintf("SELECT id, name, size, bitrate, resolution_x, resolution_y, ratio_x, ratio_y, service_id FROM %s", video.TableName)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "size", "bitrate", "resolution_x", "resolution_y", "ratio_x", "ratio_y", "service_id"}).
						AddRow(1, "my_name.mkv", 1258000, 789569, 700, 600, 4, 3, "mock_service_id"))
			},
			requestMock: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/videos/1", nil)

				return req
			},
			expectedBody:   `{"data":{"type":"videos","id":"1","attributes":{"bitrate":789569,"name":"my_name.mkv","ratio_x":4,"ratio_y":3,"resolution_x":700,"resolution_y":600,"size":1258000},"links":{"self":"/api/v1/videos/1"}}}` + "\n",
			expectedStatus: http.StatusOK,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			req := testCase.requestMock()

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error when creating a stub request, error: %s\n",
					err.Error())
			}

			if resp.StatusCode != testCase.expectedStatus {
				t.Errorf("Invalid status code, expected: %d, got: %d\n",
					testCase.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Unexpected error when reading a body, error: %s\n",
					err.Error())
			}

			if string(body) != testCase.expectedBody {
				t.Errorf("Invalid body\nexpected: %v\ngot: %v\n",
					testCase.expectedBody, string(body))
			}
		})
	}
}
