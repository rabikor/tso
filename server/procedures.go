package server

import (
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type ProcedureHandler struct {
	env config.Env
	db  *database.DB
}

func NewProceduresHandler(env config.Env, db *database.DB) ProcedureHandler {
	return ProcedureHandler{env: env, db: db}
}

func (ph ProcedureHandler) AddRoutes(rg *echo.Group) {
	router := rg.Group("/procedures")
	router.GET("", ph.All)
	router.POST("", ph.Create)
}

func (ph ProcedureHandler) All(c echo.Context) error {
	p := Pagination{Limit: ph.env.API.Request.Limit, Page: ph.env.API.Request.Page}
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	procedures, err := ph.db.Procedures.All(p.Limit, p.Offset())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": procedures, "meta": p})
}

type createProcedureRequest struct {
	Procedure struct {
		Title string `json:"title" validate:"required"`
	} `json:"procedure" validate:"required"`
}

func (r createProcedureRequest) Bind(c echo.Context, p *database.Procedure) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	p.Title = r.Procedure.Title

	return nil
}

func (ph ProcedureHandler) Create(c echo.Context) error {
	var (
		r createProcedureRequest
		p database.Procedure
	)

	if err := r.Bind(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := ph.db.Procedures.Add(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
