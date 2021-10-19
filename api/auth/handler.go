package auth

import (
	"bytes"
	"github.com/Hargeon/videocmprs/api/response"
	"github.com/Hargeon/videocmprs/pkg/repository/auth"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service"
	authsrv "github.com/Hargeon/videocmprs/pkg/service/auth"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handler struct {
	srv service.TokenAble
}

func NewHandler(db *sqlx.DB) *Handler {
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
	usr := new(user.Resource)
	bodyReader := bytes.NewReader(c.Body())
	if err := jsonapi.UnmarshalPayload(bodyReader, usr); err != nil {
		errors := []string{"Bad request"}
		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	validation := validator.New()
	if err := validation.StructPartial(usr, "Email", "Password"); err != nil {
		errors := []string{"Validation failed"}
		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	linkable, err := h.srv.GenerateToken(c.Context(), usr)
	if err != nil {
		errors := []string{err.Error()}
		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	payload, err := jsonapi.Marshal(linkable)
	if err != nil {
		errors := []string{err.Error()}
		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}
	return c.Status(http.StatusCreated).JSON(payload)
}
