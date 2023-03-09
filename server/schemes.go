package server

import (
	"net/http"
	"strconv"

	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type SchemeHandler struct {
	p  Pagination
	ir database.IllnessesRepository
	sr database.SchemesRepository
}

func NewSchemesHandler(p Pagination, ir database.IllnessesRepository, sr database.SchemesRepository) SchemeHandler {
	return SchemeHandler{p: p, ir: ir, sr: sr}
}

func (h SchemeHandler) AddRoutes(rtr *echo.Group) {
	illnessGroup := rtr.Group("/illnesses")
	illnessGroup.GET("/:illnessID/schemes", h.ByIllness)

	schemeGroup := rtr.Group("/schemes")
	schemeGroup.POST("", h.Create)
}

func (h SchemeHandler) ByIllness(c echo.Context) error {
	if err := c.Bind(&h.p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	illnessID, _ := strconv.Atoi(c.Param("illnessID"))

	schemes, err := h.sr.ByIllness(illnessID, h.p.Limit, h.p.Offset())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": schemes, "meta": h.p})
}

type createSchemeRequest struct {
	Scheme struct {
		IllnessID uint `json:"illness" validate:"required"`
		Length    uint `json:"length" validate:"required"`
	} `json:"scheme" validate:"required"`
	Days []SchemeDayData `json:"days" validate:"dive,required"`
}

func (r *createSchemeRequest) Bind(c echo.Context, s *database.Scheme) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	s.IllnessID = r.Scheme.IllnessID
	s.Length = r.Scheme.Length
	s.Days = []database.SchemeDay{}

	for _, sd := range r.Days {
		s.Days = append(
			s.Days,
			database.SchemeDay{
				ProcedureID: sd.ProcedureID,
				DrugID:      sd.DrugID,
				Order:       sd.Order,
				Times:       sd.Times,
				Frequency:   sd.Frequency,
			},
		)
	}

	return nil
}

func (h SchemeHandler) Create(c echo.Context) error {
	var (
		r createSchemeRequest
		s database.Scheme
	)

	if err := r.Bind(c, &s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	for _, sd := range r.Days {
		if err := sd.Validate(s); err != nil {
			return echo.NewHTTPError(http.StatusConflict, err)
		}
	}

	_, err := h.ir.ByID(r.Scheme.IllnessID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := h.sr.Add(s.IllnessID, s.Length, s.Days); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
