package server

import (
	"net/http"
	"strconv"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type SchemeHandler struct {
	env *config.Env
	db  *database.DB
}

func NewSchemesHandler(env *config.Env, db *database.DB) SchemeHandler {
	return SchemeHandler{env: env, db: db}
}

func (dh SchemeHandler) AddRoutes(rtr *echo.Group) {
	illnessGroup := rtr.Group("/illnesses")
	illnessGroup.GET("/:illnessID/schemes", dh.GetByIllness)

	schemeGroup := rtr.Group("/schemes")
	schemeGroup.POST("", dh.Create)
}

func (dh SchemeHandler) GetByIllness(c echo.Context) error {
	p := Pagination{Limit: dh.env.API.Request.Limit, Page: dh.env.API.Request.Page}

	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	illnessID, _ := strconv.Atoi(c.Param("illnessID"))

	schemes, err := dh.db.Schemes.GetByIllness(illnessID, p.Limit, p.GetOffset())

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": schemes, "meta": p})
}

func (dh SchemeHandler) Create(c echo.Context) error {
	var req createSchemeRequest

	s := &database.Scheme{}

	if err := req.bind(c, s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.create.bind-json", "error": err.Error()})
	}

	for _, sd := range s.Days {
		if err := ValidateSchemeDay(&sd); err != nil {
			return echo.NewHTTPError(http.StatusConflict, echo.Map{"status": false, "slug": "scheme.create.day-validation", "error": err.Error()})
		}
	}

	i, err := dh.db.Illnesses.GetById(req.Scheme.IllnessID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.create.not-found-illness", "error": err.Error()})
	}

	s.Illness = *i

	if err := dh.db.Schemes.Add(s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.create.service-request"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
