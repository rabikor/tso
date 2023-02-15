package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"treatment-scheme-organizer/api/handlers"
	"treatment-scheme-organizer/api/routes"
	"treatment-scheme-organizer/internal/configs"
	"treatment-scheme-organizer/internal/rdb"
	"treatment-scheme-organizer/internal/repositories/mysql"
	"treatment-scheme-organizer/pkg/repositories"
	"treatment-scheme-organizer/pkg/services"
)

var (
	server *gin.Engine

	DrugRepository repositories.DrugRepository
	DrugService    *services.DrugService
	DrugHandler    handlers.DrugHandler
	DrugRouter     routes.DrugRouter
)

func init() {
	DrugRepository = mysql.NewDrugRepository(rdb.DB)
	DrugService = services.NewDrugService(DrugRepository)
	DrugHandler = handlers.NewDrugHandler(DrugService)
	DrugRouter = routes.NewDrugRouter(DrugHandler)

	server = gin.Default()
}

func main() {
	router := server.Group("/api")

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "pong"})
	})

	DrugRouter.AddRoutes(router)

	log.Fatal(server.Run(":" + strconv.FormatUint(uint64(configs.Env.Server.Port), 10)))
}
