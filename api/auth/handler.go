// Package auth uses for authorization
package auth

import (
	"bytes"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"github.com/jmoiron/sqlx"

	"github.com/Hargeon/videocmprs/api/middleware"
	"github.com/Hargeon/videocmprs/api/response"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/Hargeon/videocmprs/pkg/service/auth"
)

type Handler struct {
	srv service.Tokenable
}

func NewHandler(db *sqlx.DB) *Handler {
	repo := user.NewRepository(db)
	srv := auth.NewService(repo)
	return &Handler{srv: srv}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/sign-in", h.signIn)
	router.Use(middleware.UserIdentify)
	router.Get("/me", h.retrieve)
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

// retrieve return user params
func (h *Handler) retrieve(c *fiber.Ctx) error {
	id := c.Locals("user_id").(int64)
	res, err := h.srv.Retrieve(c.Context(), id)
	if err != nil {
		errors := []string{err.Error()}
		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}
	return jsonapi.MarshalPayload(c.Status(http.StatusOK), res)
}
