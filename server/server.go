package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewErrorHandler(e *echo.Echo) {
	e.HTTPErrorHandler = ErrorHandler
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
