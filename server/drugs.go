package server

import (
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type DrugHandler struct {
	env config.Env
	db  *database.DB
}

func NewDrugsHandler(env config.Env, db *database.DB) DrugHandler {
	return DrugHandler{env: env, db: db}
}

func (h DrugHandler) AddRoutes(rg *echo.Group) {
	router := rg.Group("/drugs")
	router.GET("", h.All)
	router.POST("", h.Create)
}

func (h DrugHandler) All(c echo.Context) error {
	p := NewPagination(h.env)
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	drugs, err := h.db.Drugs.All(p.Limit, p.Offset())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": drugs, "meta": p})
}

type createDrugRequest struct {
	Drug struct {
		Title string `json:"title" validate:"required"`
	} `json:"drug" validate:"required"`
}

func (r createDrugRequest) Bind(c echo.Context, d *database.Drug) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	d.Title = r.Drug.Title

	return nil
}

func (h DrugHandler) Create(c echo.Context) error {
	var (
		r createDrugRequest
		d database.Drug
	)

	if err := r.Bind(c, &d); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := h.db.Drugs.Add(d); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
