// Package user consists of handlers for users
package user

import (
	"bytes"
	"database/sql"
	"net/http"

	"github.com/Hargeon/videocmprs/api/response"
	"github.com/Hargeon/videocmprs/pkg/repository/user"
	"github.com/Hargeon/videocmprs/pkg/service"
	usersrv "github.com/Hargeon/videocmprs/pkg/service/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"go.uber.org/zap"
)

type Handler struct {
	srv    service.Creator
	logger *zap.Logger
}

// NewHandler initialize Handler
func NewHandler(db *sql.DB, logger *zap.Logger) *Handler {
	repo := user.NewRepository(db)
	srv := usersrv.NewService(repo)

	return &Handler{srv: srv, logger: logger}
}

// InitRoutes for users
func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/", h.create)

	return router
}

// create function validate request and create user
func (h *Handler) create(c *fiber.Ctx) error {
	usr := new(user.Resource)
	bodyReader := bytes.NewReader(c.Body())

	if err := jsonapi.UnmarshalPayload(bodyReader, usr); err != nil {
		h.logger.Error("Can't unmarshal request for creating user", zap.Error(err))
		errors := []string{"Bad request"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	validation := validator.New()
	err := validation.Struct(usr)

	if err != nil {
		h.logger.Error("Validation for creating user failed", zap.Error(err))
		errors := []string{"Validation failed"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	res, err := h.srv.Create(c.Context(), usr)

	if err != nil {
		h.logger.Error("Can't create user", zap.Error(err))
		errors := []string{err.Error()}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	err = jsonapi.MarshalPayload(c.Status(http.StatusCreated), res)

	if err != nil {
		h.logger.Error("Invalid response marshaling", zap.Error(err))
		errors := []string{err.Error()}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	return nil
}
