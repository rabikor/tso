package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"treatment-scheme-organizer/config"
)

type Pagination struct {
	Limit int `query:"limit" param:"limit" json:"limit"`
	Page  int `query:"page" param:"page" json:"page"`
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func NewPagination(env config.Env) Pagination {
	return Pagination{Limit: env.API.Request.Limit, Page: env.API.Request.Page}
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

func ErrorHandler(err error, c echo.Context) {
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

func NewErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = ErrorHandler
}
