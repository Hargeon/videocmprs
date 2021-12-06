package video

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Hargeon/videocmprs/api/response"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service"
	videosrv "github.com/Hargeon/videocmprs/pkg/service/video"

	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"go.uber.org/zap"
)

const (
	IDBase    = 10
	IDBitSize = 64
)

type Handler struct {
	srv    service.Retriever
	logger *zap.Logger
}

func NewHandler(db *sql.DB, logger *zap.Logger) *Handler {
	repo := video.NewRepository(db)
	vSrv := videosrv.NewService(repo)

	return &Handler{srv: vSrv, logger: logger}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Get("/:id", h.retrieve)

	return router
}

func (h *Handler) retrieve(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, IDBase, IDBitSize)

	if err != nil || id <= 0 {
		errors := []string{"Invalid ID"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	res, err := h.srv.Retrieve(c.Context(), id)

	if err != nil {
		h.logger.Error("Get video", zap.String("Error", err.Error()),
			zap.Int64("Video ID", id))

		errors := []string{"Can not fetch video"}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	err = jsonapi.MarshalPayload(c.Status(http.StatusOK), res)

	if err != nil {
		errors := []string{err.Error()}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	return nil
}
