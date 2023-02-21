package server

import (
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type DrugHandler struct {
	db *database.DB
}

func NewDrugsHandler(db *database.DB) DrugHandler {
	return DrugHandler{db: db}
}

func (dh DrugHandler) AddRoutes(rg *echo.Group) {
	router := rg.Group("/drugs")
	router.GET("", dh.GetAll)
	router.POST("", dh.Create)
}

func (dh DrugHandler) GetAll(c echo.Context) error {
	p := Pagination{Limit: config.Env.API.Request.Limit, Page: config.Env.API.Request.Page}

	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	drugs, err := dh.db.Drugs.GetAll(p.Limit, p.GetOffset())

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": drugs})
}

type createDrugRequest struct {
	Title string `json:"title" binding:"required"`
}

func (dh DrugHandler) Create(c echo.Context) error {
	var req createDrugRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "drug.create.bind-json", "error": err})
	}

	drug := database.Drug{Title: req.Title}
	if err := dh.db.Drugs.Add(&drug); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "drug.create.service-request"})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": drug})
}
