package server

import (
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type DrugHandler struct {
	env *config.Env
	db  *database.DB
}

func NewDrugsHandler(env *config.Env, db *database.DB) DrugHandler {
	return DrugHandler{env: env, db: db}
}

func (dh DrugHandler) AddRoutes(rg *echo.Group) {
	router := rg.Group("/drugs")
	router.GET("", dh.GetAll)
	router.POST("", dh.Create)
}

func (dh DrugHandler) GetAll(c echo.Context) error {
	p := Pagination{Limit: dh.env.API.Request.Limit, Page: dh.env.API.Request.Page}

	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	drugs, err := dh.db.Drugs.GetAll(p.Limit, p.GetOffset())

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": drugs, "meta": p})
}

func (dh DrugHandler) Create(c echo.Context) error {
	var req createDrugRequest
	drug := &database.Drug{}

	if err := req.bind(c, drug); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			echo.Map{"status": false, "slug": "drug.create.bind-json", "error": err.Error()},
		)
	}

	if err := dh.db.Drugs.Add(drug); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			echo.Map{"status": false, "slug": "drug.create.service-request", "error": err.Error()},
		)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
