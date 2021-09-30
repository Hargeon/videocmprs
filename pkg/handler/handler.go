package handler

import (
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *fiber.App {
	app := fiber.New()
	api := app.Group("/api")

	v1 := api.Group("/v1")
	auth := v1.Group("/auth")
	auth.Post("/sign-in", h.signIn)
	auth.Post("/sign-up", h.signUp)

	return app
}
