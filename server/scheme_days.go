package server

import (
	"net/http"
	"strconv"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type SchemeDayHandler struct {
	env *config.Env
	db  *database.DB
}

func NewSchemeDaysHandler(env *config.Env, db *database.DB) SchemeDayHandler {
	return SchemeDayHandler{env: env, db: db}
}

func (dh SchemeDayHandler) AddRoutes(rtr *echo.Group) {
	schemeGroup := rtr.Group("/schemes")
	schemeGroup.GET("/:schemeID/days", dh.GetByScheme)
	schemeGroup.POST("/:schemeID/days", dh.CreateForScheme)
}

func (dh SchemeDayHandler) GetByScheme(c echo.Context) error {
	p := Pagination{Limit: dh.env.API.Request.Limit, Page: dh.env.API.Request.Page}

	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	schemeID, _ := strconv.Atoi(c.Param("schemeID"))

	schemeDays, err := dh.db.SchemeDays.GetByScheme(schemeID, p.Limit, p.GetOffset())

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": schemeDays, "meta": p})
}

type createSchemeDayRequest struct {
	DrugID      uint `json:"drugId"`
	ProcedureID uint `json:"procedureId"`
	DayNumber   uint `json:"dayNumber"`
	Times       uint `json:"times"`
	EachHours   uint `json:"eachHours"`
}

func (dh SchemeDayHandler) CreateForScheme(c echo.Context) error {
	var req createSchemeDayRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.bind-json", "error": err})
	}

	drug, err := dh.db.Drugs.GetById(req.DrugID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.not-found-drug", "error": err.Error()})
	}

	procedure, err := dh.db.Procedures.GetById(req.ProcedureID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.not-found-procedure", "error": err.Error()})
	}

	schemeID, _ := strconv.Atoi(c.Param("schemeID"))

	scheme, err := dh.db.Schemes.GetById(uint(schemeID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.not-found-illness", "error": err.Error()})
	}

	schemeDay := database.SchemeDay{
		DayNumber: req.DayNumber,
		Times:     req.Times,
		EachHours: req.EachHours,
		Drug:      *drug,
		Procedure: *procedure,
		Scheme:    *scheme,
	}
	if err := dh.db.SchemeDays.Add(&schemeDay); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.service-request"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
