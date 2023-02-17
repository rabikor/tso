package route

import (
	"treatment-scheme-organizer/handler"

	"github.com/labstack/echo/v4"
)

type DrugRouter struct {
	dh *handler.DrugHandler
}

func NewDrugRouter(dh *handler.DrugHandler) *DrugRouter {
	return &DrugRouter{dh: dh}
}

func (dr *DrugRouter) AddRoutes(rg *echo.Group) {
	router := rg.Group("/drugs")

	router.GET("", dr.dh.GetAll)
	router.POST("", dr.dh.Create)
}
