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

func (dh DrugHandler) AddRoutes(rg *echo.Group) {
	router := rg.Group("/drugs")
	router.GET("", dh.All)
	router.POST("", dh.Create)
}

func (dh DrugHandler) All(c echo.Context) error {
	p := Pagination{Limit: dh.env.API.Request.Limit, Page: dh.env.API.Request.Page}
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	drugs, err := dh.db.Drugs.All(p.Limit, p.Offset())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

func (dh DrugHandler) Create(c echo.Context) error {
	var (
		r createDrugRequest
		d database.Drug
	)

	if err := r.Bind(c, &d); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			echo.Map{"status": false, "error": err.Error()},
		)
	}

	if _, err := dh.db.Drugs.Add(d); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			echo.Map{"status": false, "error": err.Error()},
		)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
