package server

import (
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type IllnessHandler struct {
	env *config.Env
	db  *database.DB
}

func NewIllnessesHandler(env *config.Env, db *database.DB) IllnessHandler {
	return IllnessHandler{env: env, db: db}
}

func (ih IllnessHandler) AddRoutes(rg *echo.Group) {
	router := rg.Group("/illnesses")
	router.GET("", ih.GetAll)
	router.POST("", ih.Create)
}

func (ih IllnessHandler) GetAll(c echo.Context) error {
	p := Pagination{Limit: ih.env.API.Request.Limit, Page: ih.env.API.Request.Page}

	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	illnesses, err := ih.db.Illnesses.GetAll(p.Limit, p.GetOffset())

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": illnesses, "meta": p})
}

func (ih IllnessHandler) Create(c echo.Context) error {
	var req createIllnessRequest

	i := &database.Illness{}

	if err := req.bind(c, i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "illness.create.bind-json", "error": err.Error()})
	}

	if err := ih.db.Illnesses.Add(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"status": false, "slug": "illness.create.service-request"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
