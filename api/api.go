// Package api uses for routing
package api

import (
	"database/sql"
	"os"

	"github.com/Hargeon/videocmprs/api/auth"
	"github.com/Hargeon/videocmprs/api/middleware"
	"github.com/Hargeon/videocmprs/api/request"
	"github.com/Hargeon/videocmprs/api/user"
	"github.com/Hargeon/videocmprs/api/video"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/Hargeon/videocmprs/pkg/service/cloud"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler struct {
	db        *sql.DB
	publisher service.Publisher
}

// NewHandler returns new Handler
func NewHandler(db *sql.DB, pb service.Publisher) *Handler {
	return &Handler{db: db, publisher: pb}
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

	storage := cloud.NewS3Storage(
		os.Getenv("AWS_BUCKET_NAME"),
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"))

	v1.Mount("/requests", request.NewHandler(h.db, storage, h.publisher).InitRoutes())
	v1.Mount("/videos", video.NewHandler(h.db).InitRoutes())

	return app
}
