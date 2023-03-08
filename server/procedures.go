package server

import (
	"net/http"

	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type ProcedureHandler struct {
	p  Pagination
	pr database.ProceduresRepository
}

func NewProceduresHandler(p Pagination, pr database.ProceduresRepository) ProcedureHandler {
	return ProcedureHandler{p: p, pr: pr}
}

func (h ProcedureHandler) AddRoutes(rg *echo.Group) {
	router := rg.Group("/procedures")
	router.GET("", h.All)
	router.POST("", h.Create)
}

func (h ProcedureHandler) All(c echo.Context) error {
	if err := c.Bind(&h.p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	procedures, err := h.pr.All(h.p.Limit, h.p.Offset())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": procedures, "meta": h.p})
}

type createProcedureRequest struct {
	Procedure struct {
		Title string `json:"title" validate:"required"`
	} `json:"procedure" validate:"required"`
}

func (r *createProcedureRequest) Bind(c echo.Context, p *database.Procedure) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	p.Title = r.Procedure.Title

	return nil
}

func (h ProcedureHandler) Create(c echo.Context) error {
	var (
		r createProcedureRequest
		p database.Procedure
	)

	if err := r.Bind(c, &p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := h.pr.Add(p.Title); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
