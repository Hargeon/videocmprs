package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
)

const headerAccept = "Accept"

// AcceptHeader validate Accept header
func AcceptHeader(c *fiber.Ctx) error {
	if string(c.Request().Header.Peek(headerAccept)) != jsonapi.MediaType {
		return c.Status(http.StatusUnsupportedMediaType).SendString("Unsupported Media Type")
	}

	return c.Next()
}
