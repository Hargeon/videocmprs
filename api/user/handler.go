package user

import (
	"fmt"
	"github.com/Hargeon/videocmprs/db/model/user"
	urepo "github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service"
	uservice "github.com/Hargeon/videocmprs/pkg/service/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handler struct {
	service service.UserService
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		service: uservice.NewService(urepo.NewRepository(db)),
	}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/", h.create)
	return router
}

func (h *Handler) create(c *fiber.Ctx) error {
	fmt.Println("Create users")

	u := new(user.Resource) // user.Resource{}
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

	newUser, err := h.service.Create(c, u)
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
