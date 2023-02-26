package server

import (
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type ProcedureHandler struct {
	env *config.Env
	db  *database.DB
}

func NewProceduresHandler(env *config.Env, db *database.DB) ProcedureHandler {
	return ProcedureHandler{env: env, db: db}
}

func (ph ProcedureHandler) AddRoutes(rg *echo.Group) {
	router := rg.Group("/procedures")
	router.GET("", ph.GetAll)
	router.POST("", ph.Create)
}

func (ph ProcedureHandler) GetAll(c echo.Context) error {
	p := Pagination{Limit: ph.env.API.Request.Limit, Page: ph.env.API.Request.Page}

	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	procedures, err := ph.db.Procedures.GetAll(p.Limit, p.GetOffset())

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": procedures, "meta": p})
}

func (ph ProcedureHandler) Create(c echo.Context) error {
	var req createProcedureRequest

	p := &database.Procedure{}
	if err := req.bind(c, p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "procedure.create.bind-json", "error": err.Error()})
	}

	if err := ph.db.Procedures.Add(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "procedure.create.service-request"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
