package user

import (
	"bytes"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service"
	usersrv "github.com/Hargeon/videocmprs/pkg/service/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handler struct {
	srv service.UserService
}

func NewHandler(db *sqlx.DB) *Handler {
	repo := user.NewRepository(db)
	srv := usersrv.NewService(repo)
	return &Handler{srv: srv}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/", h.create)
	return router
}

// create not implemented yet
func (h *Handler) create(c *fiber.Ctx) error {
	usr := new(user.Resource)
	bodyReader := bytes.NewReader(c.Body())
	if err := jsonapi.UnmarshalPayload(bodyReader, usr); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	validation := validator.New()
	err := validation.Struct(usr)
	if err != nil {

	}
	//payload, err := jsonapi.Marshal(usr2)
	//if err != nil {
	//	return c.Status(http.StatusInternalServerError).SendString(err.Error())
	//}
	//return c.Status(http.StatusCreated).JSON(payload)
	return nil
}
