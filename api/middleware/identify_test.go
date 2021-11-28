package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hargeon/videocmprs/pkg/service/jwt"

	"github.com/gofiber/fiber/v2"
)

func TestUserIdentify(t *testing.T) {
	app := fiber.New()
	app.Use(UserIdentify)
	app.All("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).SendString("")
	})

	cases := []struct {
		name           string
		generateToken  func() string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "With not Bearer token",
			generateToken: func() string {
				return "pojjpjpio.[pk[pkp[k["
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":[{"title":"Should be Bearer token"}]}` + "\n",
		},
		{
			name: "Invalid Bearer token",
			generateToken: func() string {
				return "Bearer pojjpjpio.[pk[pkp[k["
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":[{"title":"token contains an invalid number of segments"}]}` + "\n",
		},
		{
			name: "Valid Bearer token",
			generateToken: func() string {
				token, err := jwt.SignedString(64)
				if err != nil {
					t.Fatalf("Unexpected error while generating jwt token")
				}
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			token := testCase.generateToken()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", token)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Unexpected error when creating a stub request\n")
			}

			if resp.StatusCode != testCase.expectedStatus {
				t.Errorf("Invelid status code, expected: %d, got: %d\n", testCase.expectedStatus,
					resp.StatusCode)
			}

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Unexpected error while reading response body, error: %v\n", err.Error())
			}

			if string(respBody) != testCase.expectedBody {
				t.Errorf("Invalid response body,\nexpected: %#v,\ngot: %#v\n", testCase.expectedBody,
					string(respBody))
			}
		})
	}
}
