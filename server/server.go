package server

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

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

func NewRouter() *echo.Echo {
	e := echo.New()
	e.Validator = NewValidator()

	return e
}
