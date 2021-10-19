// Package api uses for routing
package api

import (
	"github.com/Hargeon/videocmprs/api/auth"
	"github.com/Hargeon/videocmprs/api/middleware"
	"github.com/Hargeon/videocmprs/api/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db *sqlx.DB
}

// NewHandler returns new Handler
func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

// InitRoutes initializes and returns *fiber.App
func (h *Handler) InitRoutes() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Use(middleware.AcceptHeader)
	v1.Mount("/users", user.NewHandler(h.db).InitRoutes())
	v1.Mount("/auth", auth.NewHandler(h.db).InitRoutes())

	return app
}
