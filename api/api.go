// Package api uses for routing
package api

import (
	"github.com/Hargeon/videocmprs/api/middleware"
	"github.com/Hargeon/videocmprs/api/session"
	"github.com/Hargeon/videocmprs/api/user"
	"github.com/Hargeon/videocmprs/api/video"
	sessionrepo "github.com/Hargeon/videocmprs/pkg/repository/session"
	sessionsrv "github.com/Hargeon/videocmprs/pkg/service/session"
	"github.com/gofiber/fiber/v2"
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
	videos := h.initVideoRoutes()
	sessions := h.initSessionRoutes()

	app := fiber.New()
	api := app.Group("/api")
	// BASE_URL/api/health: check the health of the ecosystem (all dependencies: DB, RABBITMQ, ...)
	// - {postgres: connection failed}

	// BASE_URL/api/ready: {ready: OK}
	v1 := api.Group("/v1")

	v1.Mount("/users", user.NewHandler(h.db).InitRoutes())
	v1.Mount("/sessions", sessions)
	v1.Use(h.initSessionMiddleware())
	v1.Mount("/videos", videos)

	return app
}

func (h *Handler) initVideoRoutes() *fiber.App {
	router := video.NewHandler(nil)
	return router.InitRoutes()
}

func (h *Handler) initSessionRoutes() *fiber.App {
	r := sessionrepo.NewRepository(h.db)
	s := sessionsrv.NewService(r)
	router := session.NewHandler(s)
	return router.InitRoutes()
}

func (h *Handler) initSessionMiddleware() func(c *fiber.Ctx) error {
	r := sessionrepo.NewRepository(h.db)
	s := sessionsrv.NewService(r)
	return middleware.NewSessionMiddleware(s)
}
