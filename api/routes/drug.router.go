package routes

import (
	"github.com/gin-gonic/gin"
	"treatment-scheme-organizer/api/handlers"
)

type DrugRouter struct {
	dh handlers.DrugHandler
}

func NewDrugRouter(dh handlers.DrugHandler) DrugRouter {
	return DrugRouter{dh: dh}
}

func (dr *DrugRouter) AddRoutes(rg *gin.RouterGroup) {
	router := rg.Group("drugs")

	router.GET("", dr.dh.GetAll)
	router.POST("", dr.dh.Create)
}
