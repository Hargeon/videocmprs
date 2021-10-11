package video

import (
	"fmt"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.VideoService
}

// NewHandler returns new Handler
func NewHandler(s *service.VideoService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/videos", h.create)
	return router
}

func (h *Handler) create(c *fiber.Ctx) error {
	fmt.Println("Create videos")
	return nil
}
