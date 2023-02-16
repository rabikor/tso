package main

import (
	"fmt"
	"log"
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"
	"treatment-scheme-organizer/server"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}

	srv := gin.Default()
	rtr := srv.Group("/api")

	rtr.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "pong"})
	})

	drugs := server.NewDrugsHandler(db)
	drugs.AddRoutes(rtr)

	srv.Run(fmt.Sprintf(":%d", config.Env.Server.Port))
}
