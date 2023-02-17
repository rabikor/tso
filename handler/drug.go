package handler

import (
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/module"

	"github.com/labstack/echo/v4"
)

type DrugHandler struct {
	ds module.DrugStore
}

type createDrugRequest struct {
	Title string `json:"title" binding:"required"`
}

func NewDrugHandler(ds module.DrugStore) *DrugHandler {
	return &DrugHandler{ds: ds}
}

func (dh DrugHandler) GetAll(c echo.Context) error {
	p := Pagination{Page: config.Env.API.Request.Page, Limit: config.Env.API.Request.PerPage}

	if err := p.GetFromRequest(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	drugs, err := dh.ds.GetAll(p.Limit, p.Page)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": drugs})
}

func (dh DrugHandler) Create(c echo.Context) error {
	var input createDrugRequest

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	drug := module.Drug{Title: input.Title}

	if err := dh.ds.Add(&drug); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": drug})
}
