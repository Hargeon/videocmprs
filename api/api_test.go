package api

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/service"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

type cloudMock struct{}

func (c *cloudMock) Upload(ctx context.Context, header *multipart.FileHeader) (string, error) {
	if header.Filename == "failed" {
		return "", errors.New("failed connection")
	}

	return "mock_service_id", nil
}

func (c *cloudMock) URL(filename string) (string, error) {
	return filename, nil
}

type rabbitSuccess struct{}

func (r *rabbitSuccess) Publish(body []byte) error {
	return nil
}

func (r *rabbitSuccess) Ping() error {
	return nil
}

type rabbitError struct{}

func (r *rabbitError) Publish(body []byte) error {
	return errors.New("mock error")
}

func (r *rabbitError) Ping() error {
	return errors.New("mock error")
}

func TestReady(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
	}

	logger := zap.NewExample()
	defer logger.Sync()

	h := NewHandler(db, new(rabbitSuccess), new(cloudMock), logger)

	app := h.InitRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/ready", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error when creating a stub request, error: %s\n",
			err.Error())
	}

	if resp.StatusCode != 200 {
		t.Errorf("Invalid status code. expected: %d, got: %d\n",
			200, resp.StatusCode)
	}
}

func TestHealth(t *testing.T) {
	cases := []struct {
		name         string
		dbMock       func() *sql.DB
		rabbitConn   service.Publisher
		expectedBody string
	}{
		{
			name: "Valid connection",
			dbMock: func() *sql.DB {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
				}

				return db
			},
			rabbitConn:   new(rabbitSuccess),
			expectedBody: `{"DB":"OK","Rabbit":"OK"}`,
		},
		{
			name: "Invalid rabbit connection",
			dbMock: func() *sql.DB {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Unexpected error when opening a stub db connection, error: %s\n", err)
				}

				return db
			},
			rabbitConn:   new(rabbitError),
			expectedBody: `{"DB":"OK","Rabbit":"ERROR: mock error"}`,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			db := testCase.dbMock()
			logger := zap.NewExample()
			defer logger.Sync()

			h := NewHandler(db, testCase.rabbitConn, new(cloudMock), logger)

			app := h.InitRoutes()

			req := httptest.NewRequest(http.MethodGet, "/api/health", nil)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error when creating a stub request, error: %s\n",
					err.Error())
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Unexpected error when reading response body, error: %s\n",
					err.Error())
			}

			if string(body) != testCase.expectedBody {
				t.Errorf("Invalid body\nexpected: %v\ngot: %s\n",
					testCase.expectedBody, string(body))
			}
		})
	}
}
