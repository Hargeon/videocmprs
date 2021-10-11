// Package api uses for routing
package api

import (
	"github.com/Hargeon/videocmprs/api/middleware"
	"github.com/Hargeon/videocmprs/api/session"
	"github.com/Hargeon/videocmprs/api/user"
	"github.com/Hargeon/videocmprs/api/video"
	sessionrepo "github.com/Hargeon/videocmprs/pkg/repository/session"
	urepo "github.com/Hargeon/videocmprs/pkg/repository/user"
	sessionsrv "github.com/Hargeon/videocmprs/pkg/service/session"
	uservice "github.com/Hargeon/videocmprs/pkg/service/user"
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
	users := h.initUserRoutes()
	videos := h.initVideoRoutes()
	sessions := h.initSessionRoutes()

	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Mount("/", users)
	v1.Mount("/", sessions)
	v1.Use(h.initSessionMiddleware())
	v1.Mount("/", videos)

	return app
}

func (h *Handler) initUserRoutes() *fiber.App {
	r := urepo.NewRepository(h.db)
	s := uservice.NewService(r)
	router := user.NewHandler(s)
	return router.InitRoutes()
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
