package middleware

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAcceptHeader(t *testing.T) {
	app := fiber.New()

	app.Use(AcceptHeader)
	app.All("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).SendString("OK")
	})

	cases := []struct {
		name           string
		accept         string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET with accept: text/html",
			accept:         "text/html",
			method:         http.MethodGet,
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "POST with accept: text/html",
			accept:         "text/html",
			method:         http.MethodPost,
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "GET with accept: application/vnd.api+json",
			accept:         "application/vnd.api+json",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST with accept: application/vnd.api+json",
			accept:         "application/vnd.api+json",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest(testCase.method, "/", nil)
			req.Header.Set("Accept", testCase.accept)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error occured when creating stub request\n")
			}

			if resp.StatusCode != testCase.expectedStatus {
				t.Errorf("Invalid status, expected: %d, got: %d\n", testCase.expectedStatus, resp.StatusCode)
			}
		})
	}
}
