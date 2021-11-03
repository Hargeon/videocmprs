package auth

import (
	"bytes"
	"database/sql"
	"net/http"

	"github.com/Hargeon/videocmprs/api/response"
	"github.com/Hargeon/videocmprs/pkg/repository/auth"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service"
	authsrv "github.com/Hargeon/videocmprs/pkg/service/auth"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
)

type Handler struct {
	srv service.Tokenable
}

func NewHandler(db *sql.DB) *Handler {
	repo := auth.NewRepository(db)
	srv := authsrv.NewService(repo)

	return &Handler{srv: srv}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/sign-in", h.signIn)

	return router
}

func (h *Handler) signIn(c *fiber.Ctx) error {
	u := new(user.Resource)
	bodyReader := bytes.NewReader(c.Body())

	if err := jsonapi.UnmarshalPayload(bodyReader, u); err != nil {
		errors := []string{"Bad request"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	validation := validator.New()
	if err := validation.StructPartial(u, "Email", "Password"); err != nil {
		errors := []string{"Validation failed"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	resource, err := h.srv.GenerateToken(c.Context(), u)
	if err != nil {
		errors := []string{err.Error()}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	err = jsonapi.MarshalPayload(c.Status(http.StatusCreated), resource)
	if err != nil {
		errors := []string{err.Error()}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	return nil
}
