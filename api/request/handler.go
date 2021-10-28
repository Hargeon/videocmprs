package request

import (
	"bytes"
	"github.com/Hargeon/videocmprs/api/response"
	reqrepo "github.com/Hargeon/videocmprs/pkg/repository/request"
	"github.com/Hargeon/videocmprs/pkg/repository/video"
	"github.com/Hargeon/videocmprs/pkg/service"
	"github.com/Hargeon/videocmprs/pkg/service/cloud"
	"github.com/Hargeon/videocmprs/pkg/service/request"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
)

type Handler struct {
	srv service.Creator
}

func NewHandler(db *sqlx.DB) *Handler {
	reqRepo := reqrepo.NewRepository(db)
	vRepo := video.NewRepository(db)
	storage := cloud.NewS3Storage(
		os.Getenv("AWS_BUCKET_NAME"),
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"))
	srv := request.NewService(reqRepo, vRepo, storage)
	return &Handler{srv: srv}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()
	router.Post("/", h.create)
	return router
}

func (h *Handler) create(c *fiber.Ctx) error {
	file, err := c.FormFile("video")
	if err != nil {
		errors := []string{"Request does not include file"}
		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}
	reqData := c.FormValue("requests")
	buf := bytes.NewBufferString(reqData)
	res := new(reqrepo.Resource)
	if err := jsonapi.UnmarshalPayload(buf, res); err != nil {
		errors := []string{"Bad request"}
		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	validation := validator.New()
	err = validation.RegisterValidation("resolution", reqrepo.ValidateResolution)
	if err != nil {
		errors := []string{"Broken validator"}
		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	if err = validation.Struct(res); err != nil {
		errors := []string{"Validation failed"}
		return response.ErrorJsonApiResponse(c, http.StatusBadRequest, errors)
	}

	res.VideoRequest = file
	res.OriginalVideo = &video.Resource{
		Name: file.Filename,
		Size: file.Size,
	}

	res.UserID = c.Locals("user_id").(int64)

	r, err := h.srv.Create(c.Context(), res)
	if err != nil {
		errors := []string{"Can not create request"}
		return response.ErrorJsonApiResponse(c, http.StatusInternalServerError, errors)
	}

	return jsonapi.MarshalPayload(c.Status(http.StatusCreated), r)
}
