package main

import (
	"fmt"
	"log"
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"
	"treatment-scheme-organizer/server"

	"github.com/labstack/echo/v4"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}

	srv := echo.New()
	rtr := srv.Group("/api")

	rtr.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{"status": true, "message": "pong"})
	})

	drugs := server.NewDrugsHandler(db)
	drugs.AddRoutes(rtr)

	srv.Logger.Fatal(srv.Start(fmt.Sprintf(":%d", config.Env.Server.Port)))
}
