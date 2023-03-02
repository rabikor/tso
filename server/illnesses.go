package server

import (
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"

	"github.com/labstack/echo/v4"
)

type IllnessHandler struct {
	env config.Env
	db  *database.DB
}

func NewIllnessesHandler(env config.Env, db *database.DB) IllnessHandler {
	return IllnessHandler{env: env, db: db}
}

func (h IllnessHandler) AddRoutes(rg *echo.Group) {
	router := rg.Group("/illnesses")
	router.GET("", h.All)
	router.POST("", h.Create)
}

func (h IllnessHandler) All(c echo.Context) error {
	p := NewPagination(h.env)
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	illnesses, err := h.db.Illnesses.All(p.Limit, p.Offset())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": true, "data": illnesses, "meta": p})
}

type createIllnessRequest struct {
	Illness struct {
		Title string `json:"title" validate:"required"`
	} `json:"illness" validate:"required"`
}

func (r createIllnessRequest) Bind(c echo.Context, i *database.Illness) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	i.Title = r.Illness.Title

	return nil
}

func (h IllnessHandler) Create(c echo.Context) error {
	var (
		req createIllnessRequest
		i   database.Illness
	)

	if err := req.Bind(c, &i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if _, err := h.db.Illnesses.Add(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": true})
}
