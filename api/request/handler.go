package request

import (
	"bytes"
	"database/sql"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Hargeon/videocmprs/api/query"
	"github.com/Hargeon/videocmprs/api/response"
	reqrepo "github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/Hargeon/videocmprs/pkg/service/request"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"go.uber.org/zap"
)

const (
	IDBase    = 10
	IDBitSize = 64
)

type Handler struct {
	srv    service.Request
	logger *zap.Logger
}

func NewHandler(db *sql.DB, cS service.CloudStorage, pb service.Publisher, logger *zap.Logger) *Handler {
	reqRepo := reqrepo.NewRepository(db)
	vRepo := video.NewRepository(db)
	srv := request.NewService(reqRepo, vRepo, cS, pb)

	return &Handler{srv: srv, logger: logger}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/", h.create)
	router.Get("/", h.list)
	router.Get("/:id", h.retrieve)

	return router
}

func (h *Handler) create(c *fiber.Ctx) error {
	uID, ok := c.Locals("user_id").(int64)

	if !ok {
		h.logger.Error("Invalid type assertion for User ID")
		errors := []string{"Invalid user ID"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	file, err := c.FormFile("video")

	if err != nil {
		h.logger.Error("Can't Read video from request", zap.Error(err),
			zap.Int64("User ID", uID))
		errors := []string{"Request does not include file"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	ok = h.isFile(file.Header.Values("Content-Type"))

	if !ok {
		h.logger.Warn("File is not a video", zap.Int64("User ID", uID))
		errors := []string{"File is not a video"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	reqData := c.FormValue("requests")
	buf := bytes.NewBufferString(reqData)
	res := new(reqrepo.Resource)

	if err := jsonapi.UnmarshalPayload(buf, res); err != nil {
		h.logger.Error("can't unmarshal request for creating request", zap.Error(err),
			zap.Int64("User ID", uID))
		errors := []string{"Invalid request params"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	res.VideoName = file.Filename

	validation := validator.New()

	if err = validation.Struct(res); err != nil {
		h.logger.Error("Validation Failed", zap.Error(err), zap.Int64("User ID", uID))
		errors := []string{"Validation failed"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	res.VideoRequest = file
	res.OriginalVideo = &video.Resource{
		Name: file.Filename,
		Size: file.Size,
	}

	res.UserID = uID

	r, err := h.srv.Create(c.Context(), res)

	if err != nil {
		h.logger.Error("Create request", zap.Error(err),
			zap.Int64("User ID", uID))
		errors := []string{"Can not create request"}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	return jsonapi.MarshalPayload(c.Status(http.StatusCreated), r)
}

func (h *Handler) list(c *fiber.Ctx) error {
	uID, ok := c.Locals("user_id").(int64)

	if !ok {
		h.logger.Error("Invalid type assertion for User ID")
		errors := []string{"Invalid user ID"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	pageNumS := c.Query("page[number]", "0")
	pageNumI, err := strconv.Atoi(pageNumS)

	if err != nil {
		h.logger.Error("Can't convert pageNum to int", zap.Error(err),
			zap.Int64("User ID", uID), zap.String("pageNumS", pageNumS))
		errors := []string{"Invalid page number params"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	if pageNumI == 1 {
		pageNumI = 0
	}

	pageSizeS := c.Query("page[size]", "10")
	pageSizeI, err := strconv.Atoi(pageSizeS)

	if err != nil {
		h.logger.Error("Can't convert pageSize to int", zap.Error(err),
			zap.Int64("User ID", uID), zap.String("pageSizeS", pageSizeS))
		errors := []string{"Invalid page size params"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	q := &query.Params{
		RelationID: uID,
		PageNumber: uint64(pageNumI),
		PageSize:   uint64(pageSizeI),
	}

	res, err := h.srv.List(c.Context(), q)

	if err != nil {
		h.logger.Error("List requests", zap.Error(err),
			zap.Int64("User ID", uID))
		errors := []string{err.Error()}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	return jsonapi.MarshalPayload(c.Status(http.StatusOK), res)
}

func (h *Handler) isFile(types []string) bool {
	re, err := regexp.Compile(`video/.+`)

	if err != nil {
		return false
	}

	var present bool

	for _, cType := range types {
		ok := re.MatchString(cType)

		if ok {
			present = true

			break
		}
	}

	return present
}

func (h *Handler) retrieve(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, IDBase, IDBitSize)

	if err != nil || id <= 0 {
		h.logger.Error("Invalid request ID", zap.Error(err), zap.Int64("Request ID", id))
		errors := []string{"Invalid ID"}

		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	res, err := h.srv.Retrieve(c.Context(), id)

	if err != nil {
		h.logger.Error("Retrieve request", zap.Error(err),
			zap.Int64("Request ID", id))
		errors := []string{"Can not fetch video"}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	err = jsonapi.MarshalPayload(c.Status(http.StatusOK), res)

	if err != nil {
		h.logger.Error("Invalid response marshaling", zap.Error(err),
			zap.Int64("Request ID", id))
		errors := []string{err.Error()}

		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	return nil
}
