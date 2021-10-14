package response

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"net/http"
)

func ErrorJsonApiResponse(c *fiber.Ctx, status int, errors []string) error {
	var respBody []byte
	errBuf := bytes.NewBuffer(respBody)
	errObjects := make([]*jsonapi.ErrorObject, 0, len(errors))
	for _, err := range errors {
		errObject := &jsonapi.ErrorObject{Title: err}
		errObjects = append(errObjects, errObject)
	}

	err := jsonapi.MarshalErrors(errBuf, errObjects)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(status).Send(errBuf.Bytes())
}
