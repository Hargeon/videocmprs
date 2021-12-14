// Package api uses for routing
package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Hargeon/videocmprs/api/auth"
	"github.com/Hargeon/videocmprs/api/middleware"
	"github.com/Hargeon/videocmprs/api/request"
	"github.com/Hargeon/videocmprs/api/user"
	"github.com/Hargeon/videocmprs/api/video"
	"github.com/Hargeon/videocmprs/pkg/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

type Handler struct {
	db        *sql.DB
	publisher service.Publisher
	cs        service.CloudStorage
	logger    *zap.Logger
}

// NewHandler returns new Handler
func NewHandler(db *sql.DB, pb service.Publisher, cs service.CloudStorage, logger *zap.Logger) *Handler {
	return &Handler{db: db, publisher: pb, cs: cs, logger: logger}
}

// InitRoutes initializes and returns *fiber.App
func (h *Handler) InitRoutes() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Static("/docs/v1", "./docs/v1")

	api := app.Group("/api")

	api.Get("/ready", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(http.StatusOK)
	})

	api.Get("/health", h.health)

	v1 := api.Group("/v1")
	v1.Use(middleware.AcceptHeader)
	v1.Mount("/users", user.NewHandler(h.db, h.logger).InitRoutes())
	v1.Mount("/auth", auth.NewHandler(h.db, h.logger).InitRoutes())
	v1.Use(middleware.UserIdentify)

	v1.Mount("/requests", request.NewHandler(h.db, h.cs, h.publisher, h.logger).InitRoutes())
	v1.Mount("/videos", video.NewHandler(h.db, h.cs, h.logger).InitRoutes())

	return app
}

func (h *Handler) health(c *fiber.Ctx) error {
	dbStatus := "OK"
	if err := h.db.Ping(); err != nil {
		dbStatus = fmt.Sprintf("ERROR: %s", err.Error())
	}

	rabbitStatus := "OK"

	if err := h.publisher.Ping(); err != nil {
		rabbitStatus = fmt.Sprintf("ERROR: %s", err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"DB":     dbStatus,
		"Rabbit": rabbitStatus,
	})
}
