package handler

import (
	"github.com/labstack/echo/v4"
)

type Pagination struct {
	Page  uint `form:"page"`
	Limit uint `form:"limit"`
}

func (p *Pagination) GetFromRequest(c echo.Context) error {
	if err := c.Bind(&p); err != nil {
		return err
	}

	return nil
}
