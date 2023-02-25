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

type createSchemeRequest struct {
	IllnessID uint `json:"illness" binding:"required"`
	Length    uint `json:"length" binding:"required"`
}

func (dh SchemeHandler) Create(c echo.Context) error {
	var req createSchemeRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.create.bind-json", "error": err})
	}

	illness, err := dh.db.Illnesses.GetById(req.IllnessID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.create.not-found-illness", "error": err.Error()})
	}

	scheme := database.Scheme{Illness: *illness, Length: req.Length}
	if err := dh.db.Schemes.Add(&scheme); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "scheme.create.service-request"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
