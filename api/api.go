// Package api uses for routing
package api

import (
	"database/sql"

	"github.com/Hargeon/videocmprs/api/auth"
	"github.com/Hargeon/videocmprs/api/middleware"
	"github.com/Hargeon/videocmprs/api/request"
	"github.com/Hargeon/videocmprs/api/user"
	"github.com/Hargeon/videocmprs/api/video"
	"github.com/Hargeon/videocmprs/pkg/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler struct {
	db        *sql.DB
	publisher service.Publisher
	cs        service.CloudStorage
}

// NewHandler returns new Handler
func NewHandler(db *sql.DB, pb service.Publisher, cs service.CloudStorage) *Handler {
	return &Handler{db: db, publisher: pb, cs: cs}
}

// InitRoutes initializes and returns *fiber.App
func (h *Handler) InitRoutes() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Use(middleware.AcceptHeader)
	v1.Mount("/users", user.NewHandler(h.db).InitRoutes())
	v1.Mount("/auth", auth.NewHandler(h.db).InitRoutes())
	v1.Use(middleware.UserIdentify)

	v1.Mount("/requests", request.NewHandler(h.db, h.cs, h.publisher).InitRoutes())
	v1.Mount("/videos", video.NewHandler(h.db).InitRoutes())

	return app
}
