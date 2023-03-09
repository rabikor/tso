package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"treatment-scheme-organizer/database"
)

type TreatmentHandler struct {
	p  Pagination
	ir database.IllnessesRepository
	sr database.SchemesRepository
	tr database.TreatmentsRepository
}

func NewTreatmentsHandler(
	p Pagination,
	ir database.IllnessesRepository,
	sr database.SchemesRepository,
	tr database.TreatmentsRepository,
) TreatmentHandler {
	return TreatmentHandler{p: p, ir: ir, sr: sr, tr: tr}
}

func (h TreatmentHandler) AddRoutes(rtr *echo.Group) {
	illnessGroup := rtr.Group("/illnesses")
	illnessGroup.GET("/:illnessID/treatments", h.ByIllness)

	treatmentGroup := rtr.Group("/treatments")
	treatmentGroup.POST("", h.Create)
}

func (h TreatmentHandler) ByIllness(c echo.Context) error {
	if err := c.Bind(&h.p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	illnessID, _ := strconv.Atoi(c.Param("illnessID"))

	treatments, err := h.tr.ByIllness(uint(illnessID), h.p.Limit, h.p.Offset())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": treatments, "meta": h.p})
}

type createTreatmentRequest struct {
	Treatment struct {
		IllnessID uint   `json:"illness" validate:"required"`
		BegunAt   string `json:"begunAt" validate:"required"`
		EndedAt   string `json:"endedAt" validate:"required"`
	} `json:"treatment" validate:"required"`
	Schemes []TreatmentSchemeData `json:"schemes" validate:"dive,required"`
}

func (r *createTreatmentRequest) Bind(c echo.Context, t *database.Treatment) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	t.IllnessID = r.Treatment.IllnessID
	t.BegunAt = r.Treatment.BegunAt
	t.EndedAt = r.Treatment.EndedAt
	t.Schemes = []database.TreatmentScheme{}

	for loop, ts := range r.Schemes {
		t.Schemes = append(
			t.Schemes,
			database.TreatmentScheme{
				SchemeID:     ts.SchemeID,
				BeginFromDay: ts.BeginFromDay,
				Order:        uint(loop + 1),
			},
		)
	}

	return nil
}

func (h TreatmentHandler) Create(c echo.Context) error {
	var (
		r createTreatmentRequest
		t database.Treatment
	)

	if err := r.Bind(c, &t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	for _, ts := range r.Schemes {
		s, err := h.sr.ByID(ts.SchemeID)
		if err != nil {
			return echo.NewHTTPError(http.StatusConflict, err)
		}

		if err := ts.Validate(s); err != nil {
			return echo.NewHTTPError(http.StatusConflict, err)
		}
	}

	_, err := h.ir.ByID(r.Treatment.IllnessID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := h.tr.Add(t.IllnessID, t.BegunAt, t.EndedAt, t.Schemes); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
