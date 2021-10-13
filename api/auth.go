package api

// Deprecated
/*
import (
	"github.com/Hargeon/videocmprs/db/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

// signUpValidationError uses for invalid user field
type signUpValidationError struct {
	FailedField string
	Tag         string
	Value       string
}

// signUpErrorResponse ...
type signUpErrorResponse struct {
	Validation []*signUpValidationError
	Msg        string
}

// signIn uses for user authorization
func (h *Handler) signIn(c *fiber.Ctx) error {
	u := new(model.User)

	if err := c.BodyParser(u); err != nil {
		return signInBadResponse(c, "Invalid email or password")
	}

	validation := validator.New()
	if err := validation.StructPartial(u, "Email", "Password"); err != nil {
		return signInBadResponse(c, "Invalid email or password")
	}

	token, err := h.service.Authorization.GenerateToken(u.Email, u.Password)
	if err != nil {
		return signInBadResponse(c, err.Error())
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

// signUp uses for user registration
func (h *Handler) signUp(c *fiber.Ctx) error {
	u := new(model.User)
	if err := c.BodyParser(u); err != nil {
		errResponse := signUpErrorResponse{Msg: "Invalid params"}
		return c.Status(http.StatusBadRequest).JSON(errResponse)
	}

	errors := validate(u)
	if errors != nil {
		errResponse := signUpErrorResponse{Validation: errors}
		return c.Status(http.StatusBadRequest).JSON(errResponse)
	}

	u.CreatedAt = time.Now()
	id, err := h.service.Authorization.CreateUser(u)
	if err != nil {
		errResponse := signUpErrorResponse{Msg: "Can't create user"}
		return c.Status(http.StatusBadRequest).JSON(errResponse)
	}

	return c.JSON(fiber.Map{
		"id": id,
	})
}

// validate *model.User for registration
func validate(u *model.User) []*signUpValidationError {
	var errors []*signUpValidationError
	validation := validator.New()
	err := validation.Struct(u)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element signUpValidationError
			element.FailedField = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func signInBadResponse(c *fiber.Ctx, msg string) error {
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"msg": msg,
	})
}
*/
