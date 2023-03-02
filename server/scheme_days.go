package server

import (
	"fmt"
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

// ValidateSchemeDay todo: refactor
func ValidateSchemeDay(sd *database.SchemeDay) error {
	if sd.Order > sd.Scheme.Length {
		return fmt.Errorf("you cannot create day with order higher than length of scheme")
	}

	if sd.Times*sd.Frequency > 24 {
		return fmt.Errorf("you cannot create scheme day with count of medications higher length of the day")
	}

	return nil
}

func (dh SchemeDayHandler) CreateForScheme(c echo.Context) error {
	var req createSchemeDayRequest

	sd := &database.SchemeDay{}

	schemeID, _ := strconv.Atoi(c.Param("schemeID"))

	s, err := dh.db.Schemes.GetById(uint(schemeID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.not-found-scheme", "error": err.Error()})
	}

	if err := req.bind(c, sd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.bind-json", "error": err.Error()})
	}
	sd.Scheme = *s

	if err := ValidateSchemeDay(sd); err != nil {
		return echo.NewHTTPError(http.StatusConflict, echo.Map{"status": false, "slug": "scheme.day.create.validation", "error": err.Error()})
	}

	d, err := dh.db.Drugs.GetById(req.SchemeDay.DrugID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.not-found-drug", "error": err.Error()})
	}
	sd.Drug = *d

	p, err := dh.db.Procedures.GetById(req.SchemeDay.ProcedureID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.not-found-procedure", "error": err.Error()})
	}
	sd.Procedure = *p

	if err := dh.db.SchemeDays.Add(sd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.day.create.service-request", "error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
