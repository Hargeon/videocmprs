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
	return app
}
