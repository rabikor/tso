package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type Pagination struct {
	Limit int `query:"limit" param:"limit" json:"limit"`
	Page  int `query:"page" param:"page" json:"page"`
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func NewRouter(e *echo.Echo) {
	e.Validator = NewValidator()
}

func NewErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if err != nil {
		c.Logger().Error(err)

		if err := c.JSON(code, echo.Map{"status": false, "code": code, "error": err.Error()}); err != nil {
			c.Logger().Error(err)
		}
	}
}
