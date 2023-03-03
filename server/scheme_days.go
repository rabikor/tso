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
	env config.Env
	dr  database.DrugsRepository
	pr  database.ProceduresRepository
	sr  database.SchemesRepository
	sdr database.SchemeDaysRepository
}

func NewSchemeDaysHandler(env config.Env, sdr database.SchemeDaysRepository) SchemeDayHandler {
	return SchemeDayHandler{env: env, sdr: sdr}
}

func (h SchemeDayHandler) AddRoutes(rtr *echo.Group) {
	schemeGroup := rtr.Group("/schemes")
	schemeGroup.GET("/:schemeID/days", h.ByScheme)
	schemeGroup.POST("/:schemeID/days", h.CreateForScheme)
}

func (h SchemeDayHandler) ByScheme(c echo.Context) error {
	p := NewPagination(h.env)
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	schemeID, _ := strconv.Atoi(c.Param("schemeID"))

	schemeDays, err := h.sdr.ByScheme(schemeID, p.Limit, p.Offset())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": schemeDays, "meta": p})
}

type SchemeDayData struct {
	DrugID      uint `json:"drugID" validate:"required"`
	ProcedureID uint `json:"procedureID" validate:"required"`
	Order       uint `json:"order" validate:"required"`
	Times       uint `json:"times" validate:"required"`
	Frequency   uint `json:"frequency" validate:"required"`
}

type createSchemeDayRequest struct {
	SchemeDay SchemeDayData `json:"schemeDay" validate:"required"`
}

func (r createSchemeDayRequest) Bind(c echo.Context, sd *database.SchemeDay) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	sd.DrugID = r.SchemeDay.DrugID
	sd.ProcedureID = r.SchemeDay.ProcedureID
	sd.Order = r.SchemeDay.Order
	sd.Times = r.SchemeDay.Times
	sd.Frequency = r.SchemeDay.Frequency

	return nil
}

func (sdd SchemeDayData) Validate(s database.Scheme) error {
	if sdd.Order > s.Length {
		return fmt.Errorf("you cannot create day with order higher than length of scheme")
	}

	if sdd.Times*sdd.Frequency > 24 {
		return fmt.Errorf("you cannot create scheme day with count of medications higher length of the day")
	}

	return nil
}

func (h SchemeDayHandler) CreateForScheme(c echo.Context) error {
	var (
		sd database.SchemeDay
		r  createSchemeDayRequest
	)

	schemeID, _ := strconv.Atoi(c.Param("schemeID"))

	if err := r.Bind(c, &sd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	s, err := h.sr.ByID(uint(schemeID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := r.SchemeDay.Validate(s); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err)
	}

	if _, err := h.dr.ByID(r.SchemeDay.DrugID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := h.pr.ByID(r.SchemeDay.ProcedureID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := h.sdr.Add(sd); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
