package middleware

import (
	"net/http"
	"strings"

	"github.com/Hargeon/videocmprs/api/response"
	"github.com/Hargeon/videocmprs/pkg/service/jwt"
	"github.com/gofiber/fiber/v2"
)

const tokenPrefix = "Bearer "

func UserIdentify(c *fiber.Ctx) error {
	header := string(c.Request().Header.Peek("Authorization"))
	if !strings.HasPrefix(header, tokenPrefix) {
		errors := []string{"Should be Bearer token"}

		return response.ErrorJsonApiResponse(c, http.StatusUnauthorized, errors)
	}

	token := strings.TrimPrefix(header, tokenPrefix)
	id, err := jwt.ParseToken(token)

	if err != nil {
		errors := []string{err.Error()}

		return response.ErrorJsonApiResponse(c, http.StatusUnauthorized, errors)
	}

	c.Locals("user_id", id)

	return c.Next()
}
