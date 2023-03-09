package server

import (
	"fmt"
	"net/http"
	"strconv"

	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type TreatmentSchemeHandler struct {
	p   Pagination
	sr  database.SchemesRepository
	tr  database.TreatmentsRepository
	tsr database.TreatmentSchemesRepository
}

func NewTreatmentSchemesHandler(
	p Pagination,
	sr database.SchemesRepository,
	tr database.TreatmentsRepository,
	tsr database.TreatmentSchemesRepository,
) TreatmentSchemeHandler {
	return TreatmentSchemeHandler{p: p, sr: sr, tr: tr, tsr: tsr}
}

func (h TreatmentSchemeHandler) AddRoutes(rtr *echo.Group) {
	schemeGroup := rtr.Group("/treatments")
	schemeGroup.GET("/:treatmentID/schemes", h.ByTreatment)
	schemeGroup.POST("/:treatmentID/schemes", h.CreateForTreatment)
}

func (h TreatmentSchemeHandler) ByTreatment(c echo.Context) error {
	if err := c.Bind(&h.p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	treatmentID, _ := strconv.Atoi(c.Param("treatmentID"))

	treatmentSchemes, err := h.tsr.ByTreatment(uint(treatmentID), h.p.Limit, h.p.Offset())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": treatmentSchemes, "meta": h.p})
}

type TreatmentSchemeData struct {
	SchemeID     uint `json:"scheme" validate:"required"`
	BeginFromDay uint `json:"beginFromDay" validate:"required"`
}

type createTreatmentSchemeRequest struct {
	TreatmentScheme TreatmentSchemeData `json:"treatmentScheme" validate:"required"`
}

func (r *createTreatmentSchemeRequest) Bind(c echo.Context, ts *database.TreatmentScheme) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	ts.SchemeID = r.TreatmentScheme.SchemeID
	ts.BeginFromDay = r.TreatmentScheme.BeginFromDay

	return nil
}

func (tsd TreatmentSchemeData) Validate(s database.Scheme) error {
	if tsd.BeginFromDay > s.Length {
		return fmt.Errorf("you cannot create treatment scheme than begin from not exists day")
	}

	return nil
}

func (h TreatmentSchemeHandler) CreateForTreatment(c echo.Context) error {
	var (
		ts database.TreatmentScheme
		r  createTreatmentSchemeRequest
	)

	if err := r.Bind(c, &ts); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	treatmentID, _ := strconv.Atoi(c.Param("treatmentID"))
	t, err := h.tr.ByID(uint(treatmentID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	s, err := h.sr.ByID(ts.SchemeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := r.TreatmentScheme.Validate(s); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err)
	}

	countTSS, err := h.tsr.CountByTreatment(uint(treatmentID))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var nextTSSOrder = countTSS + 1
	if _, err := h.tsr.Add(t.ID, s.ID, ts.BeginFromDay, nextTSSOrder); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
