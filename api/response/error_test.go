package response

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type ErrorObject struct {
	Title string `json:"title"`
}

type ErrorResponse struct {
	Errors []ErrorObject `json:"errors"`
}

func TestErrorJsonApiResponse(t *testing.T) {
	app := fiber.New()

	cases := []struct {
		name           string
		errors         []string
		status         int
		expectedTitles []string
		expectedStatus int
	}{
		{
			name:           "With status 500",
			status:         http.StatusInternalServerError,
			errors:         []string{"check", "validation"},
			expectedTitles: []string{"check", "validation"},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "With status 400",
			status:         http.StatusBadRequest,
			errors:         []string{"check2", "validation2", "new err"},
			expectedTitles: []string{"check2", "validation2", "new err"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			app.All("/", func(ctx *fiber.Ctx) error {
				return ErrorJsonApiResponse(ctx, testCase.status, testCase.errors)
			})
			req := httptest.NewRequest("POST", "/", nil)

			resp, err := app.Test(req)
			if err != nil {
				t.Errorf("Error occured when creating a stub reguest ti fiber, error: %s\n", err.Error())
			}

			if resp.StatusCode != testCase.expectedStatus {
				t.Errorf("Invalid status code, expected: %d, got: %d\n", testCase.expectedStatus, resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error occured when reading response body, error: %s\n", err.Error())
			}

			errResp := new(ErrorResponse)
			err = json.Unmarshal(body, errResp)
			if err != nil {
				t.Fatalf("Error occured when unmarshal respone, error: %s\n", err.Error())
			}

			for _, title := range testCase.expectedTitles {
				var exists bool
				var indexExistsTitle int
				for i, err := range errResp.Errors {
					if title == err.Title {
						exists = true
						indexExistsTitle = i
					}
				}
				if !exists {
					t.Fatalf("Error doesn't present in response, title: %s\n", title)
				}
				errResp.Errors = append(errResp.Errors[:indexExistsTitle], errResp.Errors[indexExistsTitle+1:]...)
			}
		})
	}
}
