package main

import (
	"fmt"
	"net/http"

	"treatment-scheme-organizer/config"
	"treatment-scheme-organizer/database"
	"treatment-scheme-organizer/server"

	"github.com/labstack/echo/v4"
)

func main() {
	srv := echo.New()

	env := &config.Env{}
	if err := env.ParseEnv("./.env"); err != nil {
		srv.Logger.Fatal(err)
	}

	db, err := database.Open(env)
	if err != nil {
		srv.Logger.Fatal(err)
	}
	defer db.Close()

	rtr := srv.Group("/api")

	rtr.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{"status": true, "message": "pong"})
	})

	drugs := server.NewDrugsHandler(env, db)
	drugs.AddRoutes(rtr)

	illnesses := server.NewIllnessesHandler(env, db)
	illnesses.AddRoutes(rtr)

	procedures := server.NewProceduresHandler(env, db)
	procedures.AddRoutes(rtr)

	schemes := server.NewSchemesHandler(env, db)
	schemes.AddRoutes(rtr)

	schemeDays := server.NewSchemeDaysHandler(env, db)
	schemeDays.AddRoutes(rtr)

	srv.Logger.Fatal(srv.Start(fmt.Sprintf(":%d", env.Server.Port)))
}
