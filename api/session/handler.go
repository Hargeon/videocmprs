package session

import (
	"github.com/Hargeon/videocmprs/db/model/user"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Handler struct {
	service service.SessionService
}

func NewHandler(s service.SessionService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *fiber.App {
	routes := fiber.New()
	routes.Post("/sessions", h.create)
	return routes
}

func (h *Handler) create(c *fiber.Ctx) error {
	u := new(user.Resource)

	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	validation := validator.New()
	if err := validation.StructPartial(u, "Email", "Password"); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
		})
	}
	payload, err := h.service.GenerateToken(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
		})
	}
	return c.JSON(payload)
}
