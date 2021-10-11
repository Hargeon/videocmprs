package user

import (
	"fmt"
	"github.com/Hargeon/videocmprs/db/model/user"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"net/http"
)

type Handler struct {
	service service.UserService
}

func NewHandler(s service.UserService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/users", h.create)
	return router
}

func (h *Handler) create(c *fiber.Ctx) error {
	fmt.Println("Create users")

	u := new(user.Resource)
	if err := c.BodyParser(u); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	validation := validator.New()
	err := validation.Struct(u)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	c.Locals("user", u)

	newUser, err := h.service.Create(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	payload, err := jsonapi.Marshal(newUser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println(payload)
	return c.JSON(payload)
}
